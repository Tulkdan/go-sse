package main

import (
	"log"
	"net/http"
)

func main() {
	pubsub := NewPubSub()

	r := &http.ServeMux{}
	r.HandleFunc("GET /subscribe", pubsub.HandleSubscribe)
	r.HandleFunc("POST /publish", pubsub.HandlePublish)

	server := http.Server{Addr: ":8000", Handler: r}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
