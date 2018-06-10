package main

import (
	"net/http"
)

func main() {

}

var waiters map[string](chan []byte)

func WaitHandle(w http.ResponseWriter, r *http.Request) {
	_, ok := waiters[r.FormValue("id")]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := make(chan []byte)
	waiters[r.FormValue("id")] = c

	defer func() {
		delete(waiters, r.FormValue("id"))
	}()

	w.Write(<-c)
	return
}
