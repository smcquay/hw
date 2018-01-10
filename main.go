package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type v struct {
	V string `json:"version"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		r := v{"v0.0.1"}
		if err := json.NewEncoder(w).Encode(r); err != nil {
			log.Printf("json: %+v", err)
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
