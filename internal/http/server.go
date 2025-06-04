package http

import (
	"encoding/json"
	"net/http"

	"github.com/dheeraj-sn/distributed-rate-limiter/internal/limiter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	port    string
	limiter limiter.Limiter
}

func NewServer(port string, limiter limiter.Limiter) *Server {
	return &Server{port: port, limiter: limiter}
}

type CheckRequest struct {
	Key      string `json:"key"`
	Rate     int    `json:"rate"`
	Interval int    `json:"interval"` // in seconds
}

type CheckResponse struct {
	Allowed       bool `json:"allowed"`
	RetryAfterSec int  `json:"retry_after_sec"`
}

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rlaas_requests_total",
			Help: "Total number of rate limit check requests",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
}

func (s *Server) Start() error {
	http.HandleFunc("/check", s.checkHandler)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":"+s.port, nil)
}

func (s *Server) checkHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	resp := s.limiter.Allow(r.Context(), limiter.RateLimitRequest{
		Key:      req.Key,
		Rate:     req.Rate,
		Interval: req.Interval,
	})

	if resp.Allowed {
		requestsTotal.WithLabelValues("allowed").Inc()
	} else {
		requestsTotal.WithLabelValues("denied").Inc()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CheckResponse{
		Allowed:       resp.Allowed,
		RetryAfterSec: resp.RetryAfterSec,
	})
}
