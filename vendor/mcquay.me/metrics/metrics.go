package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics provides a simple way to track latency and http status
type Metrics struct {
	latency *prometheus.SummaryVec
	status  *prometheus.CounterVec
}

func New(prefix string) (*Metrics, error) {
	m := &Metrics{
		latency: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Namespace: prefix,
			Subsystem: "http",
			Name:      "request_latency_ms",
			Help:      "Latency in ms of http requests grouped by req path",
		}, []string{"path"}),

		status: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: prefix,
			Subsystem: "http",
			Name:      "status_count",
			Help:      "The count of http responses issued classified by status and api endpoint",
		}, []string{"path", "code"}),
	}
	if err := m.registerPromMetrics(); err != nil {
		return nil, errors.Wrap(err, "registration")
	}
	return m, nil
}

// registerPromMetrics registers all the metrics that eventmanger uses.
func (m *Metrics) registerPromMetrics() error {
	if err := prometheus.Register(m.latency); err != nil {
		return errors.Wrap(err, "http request latency")
	}

	if err := prometheus.Register(m.status); err != nil {
		return errors.Wrap(err, "http response counter")
	}

	return nil
}

// Wrap calls a http.Handler and tracks status code and latency.
func (m *Metrics) Wrap(prefix string, h http.Handler) http.HandlerFunc {
	return m.WrapFunc(prefix, h.ServeHTTP)
}

// WrapFunc calls a http.Handler and tracks status code and latency.
func (m *Metrics) WrapFunc(prefix string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			m.HTTPLatency(prefix, start)
		}()

		lw := NewStatusRecorder(w)
		h(lw, req)

		m.HTTPStatus(prefix, lw.Status())
	}
}

// HTTPLatency records a request latency for a given url path.
func (m *Metrics) HTTPLatency(path string, start time.Time) {
	m.latency.WithLabelValues(path).Observe(msSince(start))
}

// HTTPStatus records counts of http statuses, bucketed according to status types.
func (m *Metrics) HTTPStatus(path string, status int) {
	m.status.WithLabelValues(path, fmt.Sprintf("%d", bucketHTTPStatus(status))).Inc()
}

// bucketHTTPStatus rounds down to the nearest hundred to facilitate categorizing http statuses.
func bucketHTTPStatus(i int) int {
	return i - i%100
}

// msSince returns milliseconds since start.
func msSince(start time.Time) float64 {
	return float64(time.Since(start)) / float64(time.Millisecond)
}

// buckets returns the default prometheus buckets scaled to milliseconds.
func buckets() []float64 {
	r := []float64{}

	for _, v := range prometheus.DefBuckets {
		r = append(r, v*float64(time.Second/time.Millisecond))
	}
	return r
}
