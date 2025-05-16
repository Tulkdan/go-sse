package internal

import "net/http"

func NewServer(port string) error {
	pubsub := NewPubSub()

	r := &http.ServeMux{}
	r.HandleFunc("GET /subscribe", middlewareRequestId(middlewareTopicQueryParam(pubsub.HandleSubscribe)))
	r.HandleFunc("POST /publish", middlewareRequestId(middlewareTopicQueryParam(pubsub.HandlePublish)))

	server := http.Server{Addr: ":" + port, Handler: r}
	return server.ListenAndServe()
}
