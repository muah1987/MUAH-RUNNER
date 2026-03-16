package executor

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/muah1987/muah-runner/internal/queue"
)

func TestProcessExecutorSuccess(t *testing.T) {
	dir := t.TempDir()
	ex := NewProcessExecutor(dir, 0, 0)

	job := &queue.Job{
		ID:      "test-job-1",
		Command: []string{"echo", "hello"},
	}

	var buf bytes.Buffer
	err := ex.Run(context.Background(), job, &buf)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected 'hello' in output, got: %s", buf.String())
	}
}

func TestProcessExecutorNoCommand(t *testing.T) {
	dir := t.TempDir()
	ex := NewProcessExecutor(dir, 0, 0)

	job := &queue.Job{ID: "j1", Command: []string{}}
	err := ex.Run(context.Background(), job, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for empty command")
	}
}

func TestProcessExecutorCancelledContext(t *testing.T) {
	dir := t.TempDir()
	ex := NewProcessExecutor(dir, 0, 0)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	job := &queue.Job{
		ID:      "test-cancelled",
		Command: []string{"sleep", "10"},
	}
	err := ex.Run(ctx, job, &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for cancelled context")
	}
}

func TestProcessExecutorWithEnv(t *testing.T) {
	dir := t.TempDir()
	ex := NewProcessExecutor(dir, 0, 0)

	job := &queue.Job{
		ID:      "test-env",
		Command: []string{"sh", "-c", "echo $MY_VAR"},
		Env:     map[string]string{"MY_VAR": "test_value"},
	}

	var buf bytes.Buffer
	err := ex.Run(context.Background(), job, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "test_value") {
		t.Errorf("expected 'test_value' in output, got: %s", buf.String())
	}
}

func TestTimestampWriter(t *testing.T) {
	var buf bytes.Buffer
	tw := &timestampWriter{w: &buf}

	_, err := tw.Write([]byte("line one\nline two\n"))
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "line one") {
		t.Errorf("expected 'line one' in output, got: %s", out)
	}
}
