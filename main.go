package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const version = "v0.0.10"

type v struct {
	Hostname string `json:"hostname"`
	V        string `json:"version"`
}

func main() {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname: %+v", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		r := v{
			Hostname: hn,
			V:        version,
		}
		if err := json.NewEncoder(w).Encode(r); err != nil {
			log.Printf("json: %+v", err)
		}
	})
	http.HandleFunc("/env", func(w http.ResponseWriter, req *http.Request) {
		for _, line := range os.Environ() {
			fmt.Fprintf(w, "%v\n", line)
		}
	})
	http.HandleFunc("/live", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("live")
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("ready")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
