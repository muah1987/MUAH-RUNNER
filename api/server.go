package api

import (
"context"
"fmt"
"net/http"
"time"

"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
host string
port int
srv  *http.Server
}

func NewServer(host string, port int) *Server {
return &Server{host: host, port: port}
}

func (s *Server) Start() error {
mux := http.NewServeMux()
mux.HandleFunc("/health", healthHandler)
mux.HandleFunc("/status", statusHandler)
mux.Handle("/metrics", promhttp.Handler())
s.srv = &http.Server{
Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
Handler: mux,
}
return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
if s.srv != nil {
return s.srv.Shutdown(ctx)
}
return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
fmt.Fprintf(w, `{"status":"ok","time":"%s"}`, time.Now().Format(time.RFC3339))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
fmt.Fprintf(w, `{"runner":"muah-runner","status":"running"}`)
}
