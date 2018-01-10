package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type v struct {
	Hostname string `json:"hostname"`
	V        string `json:"version"`
}

func main() {
	panic("a bug")
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname: %+v", err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		r := v{
			Hostname: hn,
			V:        "v0.0.4",
		}
		if err := json.NewEncoder(w).Encode(r); err != nil {
			log.Printf("json: %+v", err)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
