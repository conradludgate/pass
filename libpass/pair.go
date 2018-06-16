package libpass

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mdp/qrterminal"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
)

var PASS_SERVER *url.URL = &url.URL{
	Scheme: "https",
	Host:   "pass.conradludgate.com",
} // Default

type Keys struct {
	EdPriv    [32]byte
	EdPub     [32]byte
	CurvePriv [32]byte
	CurvePub  [32]byte
}

func GenerateKeys(key KeyData) (k Keys) {
	k.EdPriv = key.PrivateKey
	k.EdPub = key.PublicKey

	// Curve25519 Private Key
	digest := sha512.Sum512(k.EdPriv[:])

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(k.CurvePriv[:], digest[:])

	// Curve25519 Public Key
	curve25519.ScalarBaseMult(&k.CurvePub, &k.CurvePriv)

	return
}

func Pair(key KeyData, w io.Writer) (err error) {
	k := GenerateKeys(key)

	// Generate public key encoding
	b := make([]byte, 0, 128)
	b = append(b, k.EdPub[:]...)
	b = append(b, k.CurvePub[:]...)

	priv := ed25519.PrivateKey(append(key.PrivateKey[:], key.PublicKey[:]...))
	b = append(b, ed25519.Sign(priv, b)...)

	data := base64.RawStdEncoding.EncodeToString(b)

	// Display QR Code
	fmt.Println()
	fmt.Println()
	fmt.Println("Scan this QR code on the Pass app to start the pairing process")
	qrterminal.GenerateHalfBlock(data, qrterminal.L, os.Stdout)

	// Setup Wait request
	c := &http.Client{
		Timeout: time.Second * 60,
	}

	u := *PASS_SERVER
	u.Path = "wait"

	// Make wait request
	resp, err := c.Post(u.String(), "", bytes.NewReader(b))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Read request body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Bad request
	if resp.StatusCode == http.StatusBadRequest {
		return errors.New(string(body))
	}

	// Some other request error
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status + ": " + string(body))
	}

	// Bad key exchange. Compromised server?
	if len(body) != 196 {
		return errors.New("Something went wrong. Try using a different server")
	}

	// body[:64] != b[:64]
	for i, v := range body[:64] {
		if b[i] != v {
			return errors.New("Something went wrong. Try using a different server")
		}
	}

	// Verify and save phone's keys
	ed_phone := body[64:96]
	if !ed25519.Verify(ed_phone, body[:128], body[128:]) {
		return errors.New("Cannot verify keys from phone. Try again")
	}

	copy(key.PhoneEdKey[:], ed_phone)
	copy(key.PhoneCurveKey[:], body[96:128])

	// Write public keys to file
	err = key.WriteKeyData(w)
	if err != nil {
		return err
	}

	// Now that keys have

	fmt.Println("Succesfully exchanged keys")
	fmt.Println("Please run 'pass verify' to verify the keys are correct")

	return
}
