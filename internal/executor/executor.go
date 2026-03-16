package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/muah1987/muah-runner/internal/queue"
)

// Executor defines the interface for running jobs.
type Executor interface {
	Run(ctx context.Context, job *queue.Job, logWriter io.Writer) error
}

// ProcessExecutor runs jobs as local OS processes.
type ProcessExecutor struct {
	WorkDir    string
	MaxRetries int
	Backoff    time.Duration
}

// NewProcessExecutor creates a ProcessExecutor.
func NewProcessExecutor(workDir string, maxRetries, backoffSecs int) *ProcessExecutor {
	return &ProcessExecutor{
		WorkDir:    workDir,
		MaxRetries: maxRetries,
		Backoff:    time.Duration(backoffSecs) * time.Second,
	}
}

// Run executes the job command as an OS process with retry logic.
func (e *ProcessExecutor) Run(ctx context.Context, job *queue.Job, logWriter io.Writer) error {
	if len(job.Command) == 0 {
		return fmt.Errorf("job %s has no command", job.ID)
	}

	jobDir := filepath.Join(e.WorkDir, job.ID)
	if err := os.MkdirAll(jobDir, 0750); err != nil {
		return fmt.Errorf("creating job work dir: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= e.MaxRetries; attempt++ {
		if attempt > 0 {
			fmt.Fprintf(logWriter, "[%s] Retry attempt %d after backoff...\n",
				timestamp(), attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(e.Backoff):
			}
		}

		err := e.runOnce(ctx, job, jobDir, logWriter)
		if err == nil {
			collectArtifacts(job, jobDir, logWriter)
			return nil
		}
		lastErr = err
		fmt.Fprintf(logWriter, "[%s] Attempt %d failed: %v\n", timestamp(), attempt+1, err)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return fmt.Errorf("job %s failed after %d attempts: %w", job.ID, e.MaxRetries+1, lastErr)
}

func (e *ProcessExecutor) runOnce(ctx context.Context, job *queue.Job, workDir string, logWriter io.Writer) error {
	cmd := exec.CommandContext(ctx, job.Command[0], job.Command[1:]...)
	cmd.Dir = workDir

	cmd.Env = os.Environ()
	for k, v := range job.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Stdout = &timestampWriter{w: logWriter}
	cmd.Stderr = &timestampWriter{w: logWriter}

	fmt.Fprintf(logWriter, "[%s] Starting: %s\n", timestamp(), strings.Join(job.Command, " "))
	return cmd.Run()
}

// DockerExecutor runs jobs inside Docker containers.
type DockerExecutor struct {
	Image      string
	WorkDir    string
	MaxRetries int
	Backoff    time.Duration
}

// NewDockerExecutor creates a DockerExecutor.
func NewDockerExecutor(image, workDir string, maxRetries, backoffSecs int) *DockerExecutor {
	return &DockerExecutor{
		Image:      image,
		WorkDir:    workDir,
		MaxRetries: maxRetries,
		Backoff:    time.Duration(backoffSecs) * time.Second,
	}
}

// Run executes the job in a Docker container.
func (e *DockerExecutor) Run(ctx context.Context, job *queue.Job, logWriter io.Writer) error {
	if len(job.Command) == 0 {
		return fmt.Errorf("job %s has no command", job.ID)
	}

	jobDir := filepath.Join(e.WorkDir, job.ID)
	if err := os.MkdirAll(jobDir, 0750); err != nil {
		return fmt.Errorf("creating job work dir: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= e.MaxRetries; attempt++ {
		if attempt > 0 {
			fmt.Fprintf(logWriter, "[%s] Retry attempt %d after backoff...\n",
				timestamp(), attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(e.Backoff):
			}
		}

		err := e.runOnce(ctx, job, jobDir, logWriter)
		if err == nil {
			collectArtifacts(job, jobDir, logWriter)
			return nil
		}
		lastErr = err
		fmt.Fprintf(logWriter, "[%s] Attempt %d failed: %v\n", timestamp(), attempt+1, err)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
	return fmt.Errorf("job %s failed after %d attempts: %w", job.ID, e.MaxRetries+1, lastErr)
}

func (e *DockerExecutor) runOnce(ctx context.Context, job *queue.Job, workDir string, logWriter io.Writer) error {
	args := []string{
		"run", "--rm",
		"-v", workDir + ":/workspace",
		"-w", "/workspace",
	}
	for k, v := range job.Env {
		args = append(args, "-e", k+"="+v)
	}
	args = append(args, e.Image)
	args = append(args, job.Command...)

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = &timestampWriter{w: logWriter}
	cmd.Stderr = &timestampWriter{w: logWriter}

	fmt.Fprintf(logWriter, "[%s] Docker run: %s\n", timestamp(), strings.Join(job.Command, " "))
	return cmd.Run()
}

func collectArtifacts(job *queue.Job, workDir string, logWriter io.Writer) {
	if len(job.Artifacts) == 0 {
		return
	}
	for i, pattern := range job.Artifacts {
		fullPattern := filepath.Join(workDir, pattern)
		matches, err := filepath.Glob(fullPattern)
		if err != nil {
			fmt.Fprintf(logWriter, "[%s] Artifact glob error for %q: %v\n", timestamp(), pattern, err)
			continue
		}
		fmt.Fprintf(logWriter, "[%s] Artifact [%d] %q: %d file(s) found\n",
			timestamp(), i, pattern, len(matches))
	}
}

func timestamp() string {
	return time.Now().Format("2006-01-02T15:04:05.000")
}

// timestampWriter prefixes each line with a timestamp.
type timestampWriter struct {
	w   io.Writer
	buf []byte
}

func (tw *timestampWriter) Write(p []byte) (n int, err error) {
	tw.buf = append(tw.buf, p...)
	for {
		idx := bytes.IndexByte(tw.buf, '\n')
		if idx < 0 {
			break
		}
		line := tw.buf[:idx+1]
		fmt.Fprintf(tw.w, "[%s] %s", timestamp(), string(line))
		tw.buf = tw.buf[idx+1:]
	}
	return len(p), nil
}
