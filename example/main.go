package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/dim13/sse"
)

const index = `<!DOCTYPE html>
<html>
<head>
<style>
	body { font-family: 'Go', sans-serif; }
</style>
<script>
	const events = new EventSource("/events");
	events.addEventListener("now", function(e) {
		document.getElementById("now").innerHTML = e.data;
	});
</script>
<body>
	<div id="now"></div>
</body>
</html>`

func main() {
	e := sse.New("now", 60)
	http.Handle("/events", e)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.Must(template.New("index").Parse(index)).Execute(w, nil)
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
