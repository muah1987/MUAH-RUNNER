package health

import (
	"testing"
)

func TestNewMonitor(t *testing.T) {
	m := New(5, "", "test-runner")
	if m == nil {
		t.Fatal("expected non-nil monitor")
	}
	stats := m.GetStats()
	if stats.JobsTotal != 0 {
		t.Error("expected zero jobs initially")
	}
}

func TestUpdateJobs(t *testing.T) {
	m := New(5, "", "test-runner")
	m.UpdateJobs(10, 2, 1)
	stats := m.GetStats()
	if stats.JobsTotal != 10 {
		t.Errorf("expected 10, got %d", stats.JobsTotal)
	}
	if stats.JobsRunning != 2 {
		t.Errorf("expected 2 running, got %d", stats.JobsRunning)
	}
	if stats.JobsFailed != 1 {
		t.Errorf("expected 1 failed, got %d", stats.JobsFailed)
	}
}

func TestMetrics(t *testing.T) {
	m := New(5, "", "test-runner")
	metrics := m.Metrics()
	if metrics == nil {
		t.Fatal("expected non-nil metrics")
	}
	if metrics.Registry() == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestHeartbeatNoURL(t *testing.T) {
	m := New(5, "", "test-runner")
	if err := m.Heartbeat(); err != nil {
		t.Errorf("expected no error for empty coordinator URL, got %v", err)
	}
}
