package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	pubsub := NewPubSub()

	r := &http.ServeMux{}
	r.HandleFunc("GET /subscribe", middlewareRequestId(middlewareTopicQueryParam(pubsub.HandleSubscribe)))
	r.HandleFunc("POST /publish", middlewareRequestId(middlewareTopicQueryParam(pubsub.HandlePublish)))

	server := http.Server{Addr: ":8000", Handler: r}
	fmt.Println("Starting server on 0.0.0.0:8000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Killing server")
	}
}
