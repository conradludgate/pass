package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/ed25519"
)

func main() {
	http.HandleFunc("/wait", WaitHandle)
	http.HandleFunc("/pair", PairHandle)

	log.Fatal(http.ListenAndServe(":7277", nil))
}

var clients map[[64]byte](chan []byte)

func WaitHandle(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 32 Bytes for client ed25519    pubkey
	// 32 Bytes for client curve25519 pubkey
	// 64 Bytes for signature
	if len(b) != 128 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Body must have size 128"))
		return
	}

	var index [64]byte
	copy(index[:], b)
	_, ok := clients[index]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request already made"))
		// Or the client some how has the private keys... Probably unlikely
		return
	}

	if !ed25519.Verify(ed25519.PublicKey(b[:32]), b[:64], b[64:]) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad signature"))
		// Client doesn't actually have the private keys they claim to have
		return
	}

	c := make(chan []byte)
	clients[index] = c

	defer func() {
		delete(clients, index)
	}()

	w.Write(<-c)
	return
}

func PairHandle(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 32 Bytes for client ed25519    pubkey
	// 32 Bytes for client curve25519 pubkey
	// 32 Bytes for phone  ed25519    pubkey
	// 32 Bytes for phone  curve25519 pubkey
	// 64 Bytes for signature
	if len(b) != 192 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Body must have size 192"))
		return
	}

	var client [64]byte
	copy(client[:], b[:64])

	c, ok := clients[client]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !ed25519.Verify(ed25519.PublicKey(b[64:96]), b[:128], b[128:]) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad signature"))
		// Phone doesn't actually have the private keys they claim to have
		return
	}

	c <- b
}
