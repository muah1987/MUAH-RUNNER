package health

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus gauges for the runner.
type Metrics struct {
	CPUUsage    prometheus.Gauge
	MemoryUsage prometheus.Gauge
	DiskUsage   prometheus.Gauge
	JobsTotal   prometheus.Gauge
	JobsRunning prometheus.Gauge
	JobsFailed  prometheus.Gauge
	registry    *prometheus.Registry
}

// NewMetrics creates and registers all Prometheus gauges.
func NewMetrics() *Metrics {
	reg := prometheus.NewRegistry()
	factory := promauto.With(reg)

	return &Metrics{
		CPUUsage: factory.NewGauge(prometheus.GaugeOpts{
			Name: "muah_runner_cpu_usage",
			Help: "Current CPU usage percentage.",
		}),
		MemoryUsage: factory.NewGauge(prometheus.GaugeOpts{
			Name: "muah_runner_memory_usage",
			Help: "Current memory usage percentage.",
		}),
		DiskUsage: factory.NewGauge(prometheus.GaugeOpts{
			Name: "muah_runner_disk_usage",
			Help: "Current disk usage percentage.",
		}),
		JobsTotal: factory.NewGauge(prometheus.GaugeOpts{
			Name: "muah_runner_jobs_total",
			Help: "Total number of jobs processed.",
		}),
		JobsRunning: factory.NewGauge(prometheus.GaugeOpts{
			Name: "muah_runner_jobs_running",
			Help: "Number of currently running jobs.",
		}),
		JobsFailed: factory.NewGauge(prometheus.GaugeOpts{
			Name: "muah_runner_jobs_failed",
			Help: "Total number of failed jobs.",
		}),
		registry: reg,
	}
}

// Registry returns the underlying Prometheus registry.
func (m *Metrics) Registry() *prometheus.Registry {
	return m.registry
}
