package internal

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func middlewareTopicQueryParam(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topicName := r.URL.Query().Get("topic")
		if topicName == "" {
			http.Error(w, errors.New("Missing 'topic' query param").Error(), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func middlewareRequestId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.NewV7()

		ctx := context.WithValue(r.Context(), "REQUEST_ID", id.String())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
