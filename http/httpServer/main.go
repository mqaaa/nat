package main

import (
	"log"
	"net/http"
	"time"
)

var Addr = ":9090"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bye", sayBye)

	server := &http.Server{
		Addr:         Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 3,
	}
	log.Println("Server Start at " + Addr)
	log.Fatal(server.ListenAndServe())
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	_, _ = w.Write([]byte("Bay!!!!"))
}
