package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/nacl/sign"
)

func main() {
	http.HandleFunc("/wait", WaitHandle)
	http.HandleFunc("/pair", PairHandle)

	log.Fatal(http.ListenAndServe(":7277", nil))
}

var clients = make(map[[64]byte](chan []byte))

var u = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WaitHandle(w http.ResponseWriter, r *http.Request) {
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer conn.Close()

	t, b, err := conn.ReadMessage()
	if err != nil {
		return
	}

	index, err := getIndex(b)
	if err != nil {
		conn.WriteMessage(t, []byte(err.Error()))
		return
	}

	c := make(chan []byte)
	clients[index] = c
	defer delete(clients, index)

	dur, _ := strconv.Atoi(r.FormValue("timeout"))
	if dur <= 0 {
		dur = 60
	}

	timeout := time.After(time.Second * time.Duration(dur))

	for {
		select {

		// Check if a message is available
		case b = <-c:
			// Send message and close the connection
			conn.WriteMessage(t, b)
			return

		case <-timeout:
			conn.WriteMessage(t, []byte("Request timed out"))
			return

		default:
			// Test connection is still alive
			if err := conn.WriteMessage(t, []byte{}); err != nil {
				return
			}
		}
	}
}

func getIndex(b []byte) (index [64]byte, err error) {
	// 32 Bytes for client sign pubkey
	// 32 Bytes for client box  pubkey
	// 64 Bytes for signature
	if len(b) != 128 {
		err = errors.New("Body must have size 128")
		return
	}

	copy(index[:], b)
	_, ok := clients[index]
	if ok {
		err = errors.New("Request already made")
		// Or the client some how has the same private keys
		// as another client... Probably unlikely
		return
	}

	var pub [32]byte
	copy(pub[:], b[64:])

	if _, verify := sign.Open(nil, b, &pub); !verify {
		err = errors.New("Bad signature")
		// Client doesn't actually have the private keys they claim to have
		return
	}

	return
}

func PairHandle(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 32 Bytes for client sign pubkey
	// 32 Bytes for client box  pubkey
	// 32 Bytes for phone  sign pubkey
	// 32 Bytes for phone  box  pubkey
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

	var pub [32]byte
	copy(pub[:], b[64:])

	if _, verify := sign.Open(nil, b, &pub); !verify {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad signature"))
		// Phone doesn't actually have the private keys they claim to have
		return
	}

	c <- b
}
