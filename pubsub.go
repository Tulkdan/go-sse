package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PubSub struct {
	channel chan string
}

func NewPubSub() *PubSub {
	return &PubSub{channel: make(chan string)}
}

func (p *PubSub) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")

	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)
	for {
		select {
		case <-clientGone:
			fmt.Println("Client disconnected")
			return
		case data := <-p.channel:
			_, err := fmt.Fprintf(w, "data: %s\n", data)
			if err != nil {
				return
			}
			if err = rc.Flush(); err != nil {
				return
			}
			time.Sleep(time.Second)
		}
	}
}

type Input struct {
	Value uint   `json:"value"`
	Name  string `json:"name"`
}

func (p *PubSub) HandlePublish(w http.ResponseWriter, r *http.Request) {
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.channel <- fmt.Sprintf("%+v", input)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message received"))
}
