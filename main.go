package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed index.html
var content string
var templates = template.Must(template.New("index.html").Parse(content))

type State struct {
	Visits  int
	Healthy bool
}

func index(s *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := templates.Execute(w, s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		s.Visits += 1
	}
}

func health(s *State) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if s.Healthy {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			params := r.URL.Query()
			h, ok := params["h"]
			if !ok {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			if len(h) != 1 || (h[0] != "true" && h[0] != "false") {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			healthy := h[0] == "true"
			if healthy == s.Healthy {
				w.WriteHeader(http.StatusNotModified)
				return
			}

			s.Healthy = healthy
			http.Redirect(w, r, "/", http.StatusSeeOther)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	addr := flag.String("addr", "127.0.0.1", "address to listen on")
	port := flag.Int("port", 8000, "TCP port to listen on")

	flag.Parse()

	s := State{Healthy: true}
	http.HandleFunc("/", index(&s))
	http.HandleFunc("/health", health(&s))

	log.Printf("Listening on %s:%d...", *addr, *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *addr, *port), nil))
}
