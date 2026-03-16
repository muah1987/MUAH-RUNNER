package health

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Stats holds current health statistics.
type Stats struct {
	CPUUsage    float64
	MemoryUsage float64
	DiskUsage   float64
	JobsTotal   int64
	JobsRunning int64
	JobsFailed  int64
}

// Monitor periodically collects health metrics.
type Monitor struct {
	mu             sync.RWMutex
	stats          Stats
	checkInterval  time.Duration
	coordinatorURL string
	runnerName     string
	stopCh         chan struct{}
	metrics        *Metrics
}

// New creates a new health Monitor.
func New(checkIntervalSecs int, coordinatorURL, runnerName string) *Monitor {
	m := &Monitor{
		checkInterval:  time.Duration(checkIntervalSecs) * time.Second,
		coordinatorURL: coordinatorURL,
		runnerName:     runnerName,
		stopCh:         make(chan struct{}),
		metrics:        NewMetrics(),
	}
	return m
}

// Start begins the periodic health check loop.
func (m *Monitor) Start() {
	go m.loop()
}

// Stop signals the monitor to stop.
func (m *Monitor) Stop() {
	close(m.stopCh)
}

// GetStats returns a snapshot of the latest stats.
func (m *Monitor) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stats
}

// UpdateJobs updates job counters used in metrics.
func (m *Monitor) UpdateJobs(total, running, failed int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stats.JobsTotal = total
	m.stats.JobsRunning = running
	m.stats.JobsFailed = failed
	m.metrics.JobsTotal.Set(float64(total))
	m.metrics.JobsRunning.Set(float64(running))
	m.metrics.JobsFailed.Set(float64(failed))
}

// Heartbeat sends a heartbeat to the coordinator.
func (m *Monitor) Heartbeat() error {
	if m.coordinatorURL == "" {
		return nil
	}
	stats := m.GetStats()
	payload := map[string]interface{}{
		"runner": m.runnerName,
		"time":   time.Now().UTC().Format(time.RFC3339),
		"cpu":    stats.CPUUsage,
		"memory": stats.MemoryUsage,
		"disk":   stats.DiskUsage,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(m.coordinatorURL+"/heartbeat", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("heartbeat returned status %d", resp.StatusCode)
	}
	return nil
}

// Metrics returns the Prometheus metrics registry.
func (m *Monitor) Metrics() *Metrics {
	return m.metrics
}

func (m *Monitor) loop() {
	ticker := time.NewTicker(m.checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.collect()
		case <-m.stopCh:
			return
		}
	}
}

func (m *Monitor) collect() {
	cpu := collectCPU()
	mem := collectMemory()
	disk := collectDisk()

	m.mu.Lock()
	m.stats.CPUUsage = cpu
	m.stats.MemoryUsage = mem
	m.stats.DiskUsage = disk
	m.mu.Unlock()

	m.metrics.CPUUsage.Set(cpu)
	m.metrics.MemoryUsage.Set(mem)
	m.metrics.DiskUsage.Set(disk)
}

func collectCPU() float64 {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0
	}
	var user, nice, system, idle, iowait, irq, softirq uint64
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "cpu ") {
			fmt.Sscanf(line, "cpu %d %d %d %d %d %d %d",
				&user, &nice, &system, &idle, &iowait, &irq, &softirq)
			total := user + nice + system + idle + iowait + irq + softirq
			if total == 0 {
				return 0
			}
			used := total - idle - iowait
			return float64(used) / float64(total) * 100
		}
	}
	return 0
}

func collectMemory() float64 {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	var total, available uint64
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fmt.Sscanf(strings.TrimPrefix(line, "MemTotal:"), "%d", &total)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fmt.Sscanf(strings.TrimPrefix(line, "MemAvailable:"), "%d", &available)
		}
	}
	if total == 0 {
		return 0
	}
	used := total - available
	return float64(used) / float64(total) * 100
}

func collectDisk() float64 {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return 0
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	if total == 0 {
		return 0
	}
	used := total - free
	return float64(used) / float64(total) * 100
}
