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
	"golang.org/x/crypto/nacl/sign"
)

var DefaultServer *url.URL = &url.URL{
	Scheme:   "https",
	Host:     "pass.conradludgate.com",
	RawQuery: "timeout=60",
} // Default

func Pair(kf KeyFile, pass_server url.URL) (err error) {
	k := kf.Client

	// Encode public keys
	b := make([]byte, 64)
	copy(b, (*k.SignPub)[:])
	copy(b[32:], (*k.BoxPub)[:])

	b = sign.Sign(nil, b, k.SignPriv)

	// Websocket Scheme
	pass_server.Scheme = strings.Replace(pass_server.Scheme, "http", "ws", 1)
	pass_server.Path += "wait"

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
	fmt.Println("Scan this QR code on the Pass app to start the pairing process")
	qrterminal.GenerateHalfBlock(data, qrterminal.L, os.Stdout)

	// Timeout
	dur, _ := strconv.Atoi(pass_server.Query().Get("timeout"))
	if dur <= 0 {
		dur = 60
	}

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

	// Bad key exchange. Compromised server?
	if len(body) != 196 {
		return errors.New(string(body))
	}

	// body[:64] != b[:64]
	for i, v := range body[:64] {
		if b[i] != v {
			return errors.New("Something went wrong. Try using a different server")
		}
	}

	// Verify and save phone's keys
	phone := kf.Phone
	copy((*phone.SignPub)[:], body[64:])
	copy((*phone.BoxPub)[:], body[96:])

	if _, verify := sign.Open(nil, body, phone.SignPub); !verify {
		return errors.New("Cannot verify keys from phone. Try again")
	}

	// Write public keys to file
	err = kf.Save()
	if err != nil {
		return err
	}

	fmt.Println("Succesfully exchanged keys")
	fmt.Println("Please run 'pass verify' to verify the keys are correct")

	return
}
