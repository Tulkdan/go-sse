package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"
)

type SubClients struct {
	id      string
	channel chan string
}

type PubSub struct {
	topics map[string][]*SubClients
	mu     sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{topics: make(map[string][]*SubClients)}
}

func (p *PubSub) getOrCreateTopic(topicName, requestId string) *SubClients {
	p.mu.Lock()
	defer p.mu.Unlock()

	clients, exists := p.topics[topicName]
	if exists == false {
		clients = []*SubClients{}
		p.topics[topicName] = clients
	}

	var currClient *SubClients
	for _, client := range clients {
		if client.id == requestId {
			currClient = client
			break
		}
	}

	if currClient == nil {
		currClient = &SubClients{id: requestId, channel: make(chan string)}
		p.topics[topicName] = append(p.topics[topicName], currClient)
	}

	return currClient
}

func (p *PubSub) clientDisconnect(topicName, requestId string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var index int
	for i, client := range p.topics[topicName] {
		if client.id == requestId {
			index = i
			break
		}
	}

	p.topics[topicName] = slices.Delete(p.topics[topicName], index, index+1)

	if len(p.topics[topicName]) == 0 {
		delete(p.topics, topicName)
	}
}

func (p *PubSub) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Query().Get("topic")
	requestId := r.Context().Value("REQUEST_ID").(string)
	clientInTopic := p.getOrCreateTopic(topicName, requestId)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")

	rc := http.NewResponseController(w)
	for {
		select {
		case <-r.Context().Done():
			fmt.Println("Client disconnected")

			p.clientDisconnect(topicName, requestId)

			return
		case data := <-clientInTopic.channel:
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

func (p *PubSub) HandlePublish(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Query().Get("topic")
	topic, exists := p.topics[topicName]
	if !exists {
		http.Error(w, errors.New("'topic' is not registered").Error(), http.StatusBadRequest)
		return
	}

	var input interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, _ := json.Marshal(input)
	for _, client := range topic {
		client.channel <- string(b)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"msg": "Message received"}`))
}
