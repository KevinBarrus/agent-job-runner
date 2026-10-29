package main

import (
	"testing"
)

var tests = []struct {
	name string
	from JobStatus
	to   JobStatus
	want bool
}{
	{
		name: "pending to running",
		from: JobPending,
		to:   JobRunning,
		want: true,
	},
	{
		name: "pending to cancelled",
		from: JobPending,
		to:   JobCancelled,
		want: true,
	},
	{
		name: "succeeded to running",
		from: JobSucceeded,
		to:   JobRunning,
		want: false,
	},
	{
		name: "running to succeeded",
		from: JobRunning,
		to:   JobSucceeded,
		want: true,
	},
	{
		name: "pending to succeeded",
		from: JobPending,
		to:   JobSucceeded,
		want: false,
	},
}

func TestCanTransition(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canTransition(tt.from, tt.to)

			if got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestCreateJobStore(t *testing.T) {
	store := NewJobStore()

	job := &Job{
		ID:     "job-1",
		Goal:   "Analyze Repository",
		Status: JobPending,
	}

	store.Create(job)

	result, ok := store.Get("job-1")

	if !ok {
		t.Fatal("expected job to exist")
	}

	if result.Status != JobPending {
		t.Errorf("want %s, got %s", JobPending, result.Status)
	}
}

func TestUpdateStatus(t *testing.T) {
	tests := []struct {
		name       string
		from       JobStatus
		to         JobStatus
		wantStatus JobStatus
		wantErr    bool
	}{
		{
			name:       "pending to running",
			from:       JobPending,
			to:         JobRunning,
			wantStatus: JobRunning,
			wantErr:    false,
		},
		{
			name:       "pending to succeeded",
			from:       JobPending,
			to:         JobSucceeded,
			wantStatus: JobPending,
			wantErr:    true,
		},
		{
			name:       "running to cancelling",
			from:       JobRunning,
			to:         JobCancelling,
			wantStatus: JobCancelling,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewJobStore()

			job := &Job{
				ID:     "job-1",
				Status: tt.from,
			}

			store.Create(job)

			err := store.UpdateStatus("job-1", tt.to)
			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr %v, got error %v", tt.wantErr, err)
			}

			result, ok := store.Get("job-1")
			if !ok {
				t.Fatal("expected job to exist")
			}

			if result.Status != tt.wantStatus {
				t.Errorf(
					"want status %s, got %s",
					tt.wantStatus,
					result.Status,
				)
			}
		})
	}
}
