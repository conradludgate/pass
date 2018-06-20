package libpass

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mdp/qrterminal"
)

// var DefaultServer *url.URL = &url.URL{
// 	Scheme: "https",
// 	Host:   "pass.conradludgate.com",
// 	//RawQuery: "timeout=60",
// } // Default

var DefaultServer string = "https://pass.conradludgate.com/"

func Pair(kf *KeyFile, server string, dur int, force bool) (err error) {
	if kf.Phone.SignPub != Zeros {
		if force {
			fmt.Println("Warning, this client is already paired to a device")
			fmt.Println("This will overwrite your currently paired device")
		} else {
			return errors.New(`Client is already paired.
If you want to ignore this error, supply the '--force' flag`)
		}
	}

	// Encode public keys
	b := make([]byte, 64)
	copy(b, kf.Client.SignPub[:])
	copy(b[32:], kf.Client.BoxPub[:])

	b = kf.Sign(b)

	// Parse URL
	pass_server, err := url.Parse(server)
	if err != nil {
		return err
	}

	pass_server.Scheme = strings.Replace(pass_server.Scheme, "http", "ws", 1)
	pass_server.Path += "wait"

	if dur <= 0 {
		dur = 60
	}
	if dur > 60 {
		dur = 60
	}

	v := url.Values{}
	v.Set("timeout", strconv.Itoa(dur))

	pass_server.RawQuery = v.Encode()

	// Start websocket connection
	conn, _, err := websocket.DefaultDialer.Dial(pass_server.String(), nil)
	if err != nil {
		return err
	}

	defer conn.Close()

	// Write public keys to server
	if err := conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return err
	}

	data := base64.RawStdEncoding.EncodeToString(b)

	// Display QR Code
	fmt.Println("On the devices page of the Pass app, tap the plus icon and scan the following QR code")
	qrterminal.GenerateHalfBlock(data, qrterminal.L, os.Stdout)

	timeout := time.After(time.Second * time.Duration(dur))

	// Wait for response
	var body []byte
	for len(body) == 0 {
		select {
		case <-timeout:
			return errors.New("Request timed out")

		default:
			_, body, err = conn.ReadMessage()
			if err != nil {
				return err
			}
		}
	}

	/*
		body := struct {
			Checksum           [64]byte // body[0:64]
			kf.Client.SignPub [32]byte // body[64:96]
			kf.Client.BoxPub  [32]byte // body[96:128]
			kf.Phone.SignPub  [32]byte // body[128:160]
			kf.Phone.BoxPub   [32]byte // body[160:192]
		}
	*/

	if len(body) != 192 {
		return errors.New(string(body))
	}

	// body[64:128] != b[64:128]
	for i, v := range body[64:128] {
		if b[64+i] != v {
			return errors.New("Something went wrong. Try using a different server")
		}
	}

	// Verify and save phone's keys
	copy(kf.Phone.SignPub[:], body[128:])

	if _, ok := kf.Open(body); !ok {
		return errors.New("Cannot verify keys from phone. Try again")
	}

	copy(kf.Phone.BoxPub[:], body[160:])

	// Write public keys to file
	err = kf.Save()
	if err != nil {
		return err
	}

	fmt.Println("Succesfully exchanged keys")
	fmt.Println("Please run 'pass verify' to verify the keys are correct")

	return
}
