package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/conradludgate/pass/libpass"
	"golang.org/x/crypto/nacl/sign"
)

func main() {
	file, err := os.OpenFile(".pass.key", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("Cannot open file. Reason:", err.Error())
		return
	}
	defer file.Close()

	key, err := libpass.NewKeyFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	key.Phone.Generate()
	b := make([]byte, 128)
	copy(b, key.Client.SignPub[:])
	copy(b[32:], key.Client.BoxPub[:])
	copy(b[64:], key.Phone.SignPub[:])
	copy(b[96:], key.Phone.BoxPub[:])

	b = sign.Sign(nil, b, &key.Phone.SignPriv)

	resp, err := http.Post("http://localhost:7277/pair", "", bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode, resp.Status)
	b, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
}
