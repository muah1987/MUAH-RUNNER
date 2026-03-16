package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muah1987/muah-runner/internal/config"
	"github.com/muah1987/muah-runner/internal/health"
	"github.com/muah1987/muah-runner/internal/queue"
	"github.com/muah1987/muah-runner/internal/sensor"
)

func newTestServer() *Server {
	cfg := config.Defaults()
	caps := &sensor.Capabilities{
		OS:     "linux",
		Arch:   "amd64",
		Labels: []string{"os:linux"},
	}
	monitor := health.New(5, "", "test-runner")
	q := queue.New()
	return NewServer(cfg, caps, monitor, q)
}

func TestHandleHealth(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	s.handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status ok, got %v", resp["status"])
	}
}

func TestHandleStatus(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rr := httptest.NewRecorder()
	s.handleStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestHandleInfo(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/info", nil)
	rr := httptest.NewRecorder()
	s.handleInfo(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestSubmitJob(t *testing.T) {
	s := newTestServer()
	body := `{"name":"test","priority":5,"command":["echo","hello"]}`
	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()
	s.handleJobs(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
	if s.queue.Len() != 1 {
		t.Errorf("expected 1 job in queue, got %d", s.queue.Len())
	}
}

func TestListJobs(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	rr := httptest.NewRecorder()
	s.handleJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCancelJob(t *testing.T) {
	s := newTestServer()
	s.queue.Enqueue(&queue.Job{ID: "test-id"})

	req := httptest.NewRequest(http.MethodDelete, "/jobs/test-id", nil)
	req.URL.Path = "/jobs/test-id"
	rr := httptest.NewRecorder()
	s.handleJobByID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGitHubWebhookNoSecret(t *testing.T) {
	s := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/webhooks/github", bytes.NewBufferString(`{}`))
	req.Header.Set("X-GitHub-Event", "push")
	rr := httptest.NewRecorder()
	s.handleGitHubWebhook(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
