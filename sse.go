package sse

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Broker for SSE messages
type Broker struct {
	eventName string
	queueSize int
	clients   *sync.Map
}

// Write implements io.Writer
func (b *Broker) Write(p []byte) (n int, err error) {
	b.clients.Range(func(key, value any) bool {
		ch := key.(chan string)
		select {
		case ch <- string(p):
		default:
		}
		return true
	})
	return len(p), nil
}

// New SSE Broker
func New(eventName string, queueSize int) *Broker {
	if queueSize == 0 {
		queueSize = 10
	}
	return &Broker{
		eventName: eventName,
		queueSize: queueSize,
		clients:   new(sync.Map),
	}
}

// ServeHTTP implements http.Handler
func (b Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "not a flusher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan string, b.queueSize)
	defer close(ch)

	b.clients.Store(ch, nil)
	defer b.clients.Delete(ch)

	for {
		select {
		case <-r.Context().Done():
			return
		case data := <-ch:
			if b.eventName != "" {
				fmt.Fprintf(w, "event: %s\n", b.eventName)
			}
			for _, s := range strings.Split(data, "\n") {
				fmt.Fprintf(w, "data: %s\n", s)
			}
			fmt.Fprintf(w, "\n")
			flusher.Flush()
		}
	}
}
