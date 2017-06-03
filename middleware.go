package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

type metricsMiddleware struct {
	requests *prometheus.CounterVec
	latency  *prometheus.HistogramVec
}

func (m *metricsMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(w, req)
	duration := float64(time.Since(start).Nanoseconds()) / 1000000

	rw := w.(negroni.ResponseWriter)

	log.WithFields(logrus.Fields{
		"status":   rw.Status(),
		"method":   req.Method,
		"path":     req.URL.Path,
		"host":     req.Host,
		"duration": duration,
	}).Print("request")
	m.requests.WithLabelValues(fmt.Sprintf("%v", rw.Status()), req.Method, req.URL.Path).Inc()
	m.latency.WithLabelValues(fmt.Sprintf("%v", rw.Status()), req.Method, req.URL.Path).Observe(duration)
}

func newMetricsMiddleware() *metricsMiddleware {
	m := &metricsMiddleware{}

	m.requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.requests)

	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "app_requests_duration_milliseconds",
		Help:    "How long it took to process the request, partitioned by status code, method and HTTP path.",
		Buckets: []float64{10, 100, 300, 1200, 5000},
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)

	return m
}

func newMiddlewares(m http.Handler) *negroni.Negroni {
	n := negroni.New()

	recovery := negroni.NewRecovery()
	recovery.Logger = log
	recovery.PrintStack = config.IsDevelopment
	n.Use(recovery)

	n.Use(newMetricsMiddleware())

	n.Use(negroni.NewStatic(FS(config.IsDevelopment)))

	n.UseHandler(m)

	return n
}
