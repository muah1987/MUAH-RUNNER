package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/muah1987/muah-runner/internal/config"
	"github.com/muah1987/muah-runner/internal/health"
	"github.com/muah1987/muah-runner/internal/queue"
	"github.com/muah1987/muah-runner/internal/sensor"
)

// Server is the HTTP API server.
type Server struct {
	cfg     *config.Config
	caps    *sensor.Capabilities
	monitor *health.Monitor
	queue   *queue.Queue
	httpSrv *http.Server
}

// NewServer creates a new API server.
func NewServer(cfg *config.Config, caps *sensor.Capabilities, monitor *health.Monitor, q *queue.Queue) *Server {
	s := &Server{
		cfg:     cfg,
		caps:    caps,
		monitor: monitor,
		queue:   q,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.Handle("/metrics", promhttp.HandlerFor(monitor.Metrics().Registry(), promhttp.HandlerOpts{}))
	mux.HandleFunc("/status", s.handleStatus)
	mux.HandleFunc("/jobs", s.handleJobs)
	mux.HandleFunc("/jobs/", s.handleJobByID)
	mux.HandleFunc("/webhooks/github", s.handleGitHubWebhook)
	mux.HandleFunc("/api/v1/info", s.handleInfo)

	s.httpSrv = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	return s.httpSrv.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
