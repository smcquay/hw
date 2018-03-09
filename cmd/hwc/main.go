package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mcquay.me/hw"
	"mcquay.me/metrics"
)

var success *prometheus.CounterVec

type state struct {
	Hostname string `json:"hostname"`
	V        string `json:"version"`

	sync.RWMutex
	Counts map[string]int `json:"counts"`
}

func (s *state) update(target string) {
	u := fmt.Sprintf("http://%s:8080/", target)
	for {
		time.Sleep(100 * time.Millisecond)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			panic(err)
		}
		req.Close = true
		httpResp, err := http.DefaultClient.Do(req)
		if err != nil {
			success.WithLabelValues("false").Inc()
			continue
		}
		success.WithLabelValues("true").Inc()
		rv := hw.V{}
		if err := json.NewDecoder(httpResp.Body).Decode(&rv); err != nil {
			panic(err)
		}
		if err := httpResp.Body.Close(); err != nil {
			panic(err)
		}
		s.Lock()
		s.Counts[rv.Hostname] += 1
		s.Unlock()
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: hwc <target hostname>")
	}
	target := os.Args[1]

	m, err := metrics.New("hwc")
	if err != nil {
		log.Fatalf("metrics: %v", err)
	}

	success = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hwc_success_count",
			Help: "counts successes and failures of  fetch attempts",
		},
		[]string{"ok"},
	)

	if err := prometheus.Register(success); err != nil {
		log.Fatalf("registering success: %v", err)
	}

	hn, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname: %+v", err)
	}

	fetcher := state{
		Hostname: hn,
		V:        hw.Version,
		Counts:   map[string]int{},
	}
	go fetcher.update(target)

	http.HandleFunc("/", m.WrapFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fetcher.RLock()
		defer fetcher.RUnlock()
		if err := json.NewEncoder(w).Encode(fetcher); err != nil {
			log.Printf("json: %+v", err)
		}
	}))

	http.HandleFunc("/live", m.WrapFunc("/live", hw.OK))
	http.HandleFunc("/ready", m.WrapFunc("/ready", hw.OK))

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("listen and serve: %v", err)
	}
}
