package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/muah1987/muah-runner/internal/queue"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stats := s.monitor.GetStats()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"cpu":    stats.CPUUsage,
		"memory": stats.MemoryUsage,
		"disk":   stats.DiskUsage,
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stats := s.monitor.GetStats()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":         s.cfg.Runner.Name,
		"labels":       s.caps.Labels,
		"jobs_total":   stats.JobsTotal,
		"jobs_running": stats.JobsRunning,
		"jobs_failed":  stats.JobsFailed,
		"queue_len":    s.queue.Len(),
	})
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":    s.cfg.Runner.Name,
		"version": "dev",
		"labels":  s.caps.Labels,
		"os":      s.caps.OS,
		"arch":    s.caps.Arch,
	})
}

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.queue.List())
	case http.MethodPost:
		s.submitJob(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

type jobRequest struct {
	Name      string            `json:"name"`
	Priority  int               `json:"priority"`
	Command   []string          `json:"command"`
	Env       map[string]string `json:"env"`
	Artifacts []string          `json:"artifacts"`
}

func (s *Server) submitJob(w http.ResponseWriter, r *http.Request) {
	var req jobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Command) == 0 {
		http.Error(w, "command is required", http.StatusBadRequest)
		return
	}
	job := &queue.Job{
		Name:      req.Name,
		Priority:  req.Priority,
		Command:   req.Command,
		Env:       req.Env,
		Artifacts: req.Artifacts,
	}
	s.queue.Enqueue(job)
	writeJSON(w, http.StatusCreated, job)
}

func (s *Server) handleJobByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/jobs/")
	if id == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		for _, j := range s.queue.List() {
			if j.ID == id {
				writeJSON(w, http.StatusOK, j)
				return
			}
		}
		http.Error(w, "job not found", http.StatusNotFound)
	case http.MethodDelete:
		if s.queue.Cancel(id) {
			writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
		} else {
			http.Error(w, "job not found or not cancellable", http.StatusNotFound)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	if s.cfg.Webhooks.Secret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if !verifyGitHubSignature(body, s.cfg.Webhooks.Secret, sig) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
	}

	event := r.Header.Get("X-GitHub-Event")
	fmt.Printf("GitHub webhook received: event=%s\n", event)

	writeJSON(w, http.StatusOK, map[string]string{"status": "received", "event": event})
}

func verifyGitHubSignature(body []byte, secret, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
