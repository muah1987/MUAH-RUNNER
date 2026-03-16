package api

import "github.com/prometheus/client_golang/prometheus"

var (
RunsTotal = prometheus.NewCounterVec(
prometheus.CounterOpts{Name: "muah_runs_total", Help: "Total runs"},
[]string{"status"},
)
RunDuration = prometheus.NewHistogramVec(
prometheus.HistogramOpts{Name: "muah_run_duration_seconds", Help: "Run duration"},
[]string{"status"},
)
)

func init() {
prometheus.MustRegister(RunsTotal, RunDuration)
}
