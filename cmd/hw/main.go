package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mcquay.me/hw"
	"mcquay.me/metrics"
)

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
		r := hw.V{
			Hostname: hn,
			V:        hw.Version,
			G:        hw.Git,
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
	http.HandleFunc("/live", m.WrapFunc("/live", hw.OK))
	http.HandleFunc("/ready", m.WrapFunc("/ready", hw.OK))
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
