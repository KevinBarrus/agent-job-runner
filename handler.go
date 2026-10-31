package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type JobHandler struct {
	store *JobStore
}

type CreateJobRequest struct {
	Goal string `json:"goal"`
}

func (jh *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateJobRequest

	// 将 JSON 数据解析到 req 结构体中
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// 校验 req.Goal
	req.Goal = strings.TrimSpace(req.Goal)
	if req.Goal == "" {
		http.Error(w, "Empty Goal", http.StatusBadRequest)
		return
	}

	// 生成 Job ID
	idBytes := make([]byte, 16)

	_, err = rand.Read(idBytes)
	if err != nil {
		http.Error(w, "Failed to generate job ID", http.StatusInternalServerError)
		return
	}

	id := hex.EncodeToString(idBytes)

	// 创建并保存 Job
	now := time.Now()
	job := &Job{
		ID:        id,
		Goal:      req.Goal,
		Status:    JobPending,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	jh.store.Create(job)

	// 返回 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := map[string]any{
		"id":     job.ID,
		"status": job.Status,
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write job response: %v", err)
	}
}
