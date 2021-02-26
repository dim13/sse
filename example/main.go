package main

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/dim13/sse"
)

//go:embed index.tmpl
var static embed.FS

func main() {
	e := sse.New("now", 60)
	http.Handle("/events", e)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.ParseFS(static, "index.tmpl")).Execute(w, nil)
	})
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for now := range t.C {
			fmt.Fprint(e, now)
		}
	}()
	http.ListenAndServe("localhost:6060", nil)
}
