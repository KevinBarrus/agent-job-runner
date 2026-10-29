package main

import (
	"errors"
	"sync"
	"time"
)

type JobStatus string

const (
	JobPending    JobStatus = "pending"
	JobRunning    JobStatus = "running"
	JobCancelling JobStatus = "cancelling"
	JobSucceeded  JobStatus = "succeeded"
	JobFailed     JobStatus = "failed"
	JobCancelled  JobStatus = "cancelled"
)

func canTransition(from, to JobStatus) bool {
	switch from {
	case JobPending:
		if to == JobRunning || to == JobCancelled {
			return true
		}

	case JobRunning:
		if to == JobSucceeded || to == JobCancelling || to == JobFailed {
			return true
		}

	case JobCancelling:
		if to == JobCancelled {
			return true
		}

	default:
		return false
	}
	return false
}

type Job struct {
	ID        string
	Goal      string
	Status    JobStatus
	Progress  int
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobStore struct {
	mu   sync.Mutex
	jobs map[string]Job
}

func NewJobStore() *JobStore {
	// 创建并初始化 JobStore
	return &JobStore{
		jobs: make(map[string]Job),
	}
}

func (s *JobStore) Create(job *Job) {
	// 将 Job 保存到 map 中
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = *job
}

func (s *JobStore) Get(id string) (Job, bool) {
	// 根据 ID 查询 Job，返回副本与是否存在
	s.mu.Lock()
	defer s.mu.Unlock()
	job, exist := s.jobs[id]
	if exist == false {
		return Job{}, false
	}
	return job, exist
}

func (s *JobStore) UpdateStatus(id string, next JobStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]

	// 1: Job 不存在时返回 error
	if !ok {
		return errors.New("job not found")
	}

	// 2: 使用 canTransition() 检查转换是否合法
	// 如果不合法，返回 error
	if !canTransition(job.Status, next) {
		return errors.New("illegal status transition")
	}

	// 3: 修改 job.Status
	job.Status = next

	// 4: 更新 job.UpdatedAt
	job.UpdatedAt = time.Now()

	// 5: 将修改后的 job 写回 Map
	s.jobs[id] = job

	return nil
}
