package queue

import (
	"testing"
)

func TestEnqueueDequeue(t *testing.T) {
	q := New()

	q.Enqueue(&Job{ID: "j1", Name: "low", Priority: 1})
	q.Enqueue(&Job{ID: "j2", Name: "high", Priority: 10})
	q.Enqueue(&Job{ID: "j3", Name: "mid", Priority: 5})

	job := q.Dequeue()
	if job == nil {
		t.Fatal("expected a job, got nil")
	}
	if job.ID != "j2" {
		t.Errorf("expected highest priority job j2, got %s", job.ID)
	}

	job = q.Dequeue()
	if job.ID != "j3" {
		t.Errorf("expected j3, got %s", job.ID)
	}

	job = q.Dequeue()
	if job.ID != "j1" {
		t.Errorf("expected j1, got %s", job.ID)
	}

	if q.Dequeue() != nil {
		t.Error("expected nil from empty queue")
	}
}

func TestCancel(t *testing.T) {
	q := New()
	q.Enqueue(&Job{ID: "j1", Priority: 1})

	if !q.Cancel("j1") {
		t.Error("expected cancel to succeed")
	}
	if q.Cancel("j1") {
		t.Error("expected cancel to fail for already-cancelled job")
	}

	if q.Dequeue() != nil {
		t.Error("expected nil dequeue after cancel")
	}
}

func TestList(t *testing.T) {
	q := New()
	q.Enqueue(&Job{ID: "j1"})
	q.Enqueue(&Job{ID: "j2"})

	list := q.List()
	if len(list) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(list))
	}
}

func TestLen(t *testing.T) {
	q := New()
	if q.Len() != 0 {
		t.Error("expected empty queue")
	}
	q.Enqueue(&Job{ID: "j1"})
	if q.Len() != 1 {
		t.Errorf("expected len 1, got %d", q.Len())
	}
}

func TestAutoID(t *testing.T) {
	q := New()
	job := &Job{Name: "no-id"}
	q.Enqueue(job)
	if job.ID == "" {
		t.Error("expected auto-generated ID")
	}
}
