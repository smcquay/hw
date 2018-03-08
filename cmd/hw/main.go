package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mcquay.me/metrics"
)

const version = "v0.1.0"

type v struct {
	Hostname string `json:"hostname"`
	V        string `json:"version"`
}

func main() {
	m, err := metrics.New("hw")
	if err != nil {
		log.Fatalf("metrics: %v", err)
	}
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname: %+v", err)
	}
	http.HandleFunc("/", m.WrapFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		r := v{
			Hostname: hn,
			V:        version,
		}
		if err := json.NewEncoder(w).Encode(r); err != nil {
			log.Printf("json: %+v", err)
		}
	}))
	http.HandleFunc("/env", m.WrapFunc("/env", func(w http.ResponseWriter, req *http.Request) {
		for _, line := range os.Environ() {
			fmt.Fprintf(w, "%v\n", line)
		}
	}))

	codes := []int{
		http.StatusBadGateway,
		http.StatusBadRequest,
		http.StatusUnauthorized,
	}
	http.HandleFunc("/bad", m.WrapFunc("/bad", func(w http.ResponseWriter, req *http.Request) {
		code := codes[rand.Intn(len(codes))]
		w.WriteHeader(code)
	}))
	http.HandleFunc("/live", m.WrapFunc("/live", ok))
	http.HandleFunc("/ready", m.WrapFunc("/ready", ok))
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}

func ok(w http.ResponseWriter, req *http.Request) {}
