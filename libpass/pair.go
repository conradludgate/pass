package pass

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/mdp/qrterminal"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
)

var PASS_SERVER url.URL = url.URL{
	Scheme: "https",
	Host:   "pass.conradludgate.com",
} // Default

type Keys struct {
	EdPriv    [32]byte
	EdPub     [32]byte
	CurvePriv [32]byte
	CurvePub  [32]byte
}

func GenerateKeys(key [64]byte) (k Keys) {
	// Ed25519 Private Key
	copy(k.EdPriv[:], key[:32])

	// Ed25519 Public Key
	copy(k.EdPub[:], key[32:])

	// Curve25519 Private Key
	digest := sha512.Sum512(key[:32])

	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(k.CurvePriv[:], digest[:])

	// Curve25519 Public Key
	curve25519.ScalarBaseMult(&k.CurvePub, &k.CurvePriv)

	return
}

func Pair(key [64]byte) (err error) {
	k := GenerateKeys(key)

	b := make([]byte, 0, 128)
	b = append(b, k.EdPub[:]...)
	b = append(b, k.CurvePub[:]...)

	priv := ed25519.PrivateKey(key[:])
	b = append(b, ed25519.Sign(priv, b)...)

	data := base64.RawStdEncoding.EncodeToString(b)
	fmt.Println(data)
	qrterminal.GenerateHalfBlock(data, qrterminal.L, os.Stdout)

	c := &http.Client{
		Timeout: time.Second * 60,
	}

	resp, err := c.Post(PASS_SERVER, "", b)
	if err != nil {
		return err
	}

	defer resp.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return errors.New(string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status + ": " + string(body))
	}

	if body[:64] != b[:64] {
		fmt.Println("Key exchange failed")
		return nil
	}

	// Verify and save phone's keys

	return nil
}
