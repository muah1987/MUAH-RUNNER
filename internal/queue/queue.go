package queue

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the current state of a job.
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

// Job represents a unit of work.
type Job struct {
	ID          string
	Name        string
	Priority    int
	Status      Status
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Command     []string
	Env         map[string]string
	Artifacts   []string
	Error       string
}

// Queue is a thread-safe priority job queue.
type Queue struct {
	mu   sync.Mutex
	jobs []*Job
}

// New creates a new empty Queue.
func New() *Queue {
	return &Queue{}
}

// Enqueue adds a job to the queue. Jobs with higher Priority values are dequeued first.
func (q *Queue) Enqueue(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if job.ID == "" {
		job.ID = fmt.Sprintf("job-%d", time.Now().UnixNano())
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}
	job.Status = StatusPending
	q.jobs = append(q.jobs, job)
}

// Dequeue removes and returns the highest-priority pending job, or nil if none available.
func (q *Queue) Dequeue() *Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	best := -1
	for i, j := range q.jobs {
		if j.Status != StatusPending {
			continue
		}
		if best == -1 || j.Priority > q.jobs[best].Priority {
			best = i
		}
	}
	if best == -1 {
		return nil
	}
	job := q.jobs[best]
	q.jobs = append(q.jobs[:best], q.jobs[best+1:]...)
	return job
}

// Cancel marks a pending job as cancelled.
func (q *Queue) Cancel(id string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, j := range q.jobs {
		if j.ID == id && j.Status == StatusPending {
			j.Status = StatusCancelled
			return true
		}
	}
	return false
}

// List returns a snapshot of all jobs currently tracked.
func (q *Queue) List() []*Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	out := make([]*Job, len(q.jobs))
	copy(out, q.jobs)
	return out
}

// Len returns the number of jobs in the queue (all statuses).
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.jobs)
}
