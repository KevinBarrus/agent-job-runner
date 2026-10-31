package main

import (
	"log"
	"net/http"
)

func main() {
	// 创建一个 HTTP 请求路由器
	mux := http.NewServeMux()

	// 根据请求的方法和路径，将请求分发给对应的 Handler
	store := NewJobStore()

	handler := &JobHandler{
		store: store,
	}

	mux.HandleFunc("POST /jobs", handler.Create)
	mux.HandleFunc("GET /jobs/{id}", handler.Get)

	log.Println("Server listening on: 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
