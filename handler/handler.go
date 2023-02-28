package handler

import (
	"net/http"
	"sync"
)

type ProxyRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type ProxyResponse struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Lenght  int               `json:"length"`
}

type Handler struct {
	Requests  map[string]ProxyRequest
	Responses map[string]ProxyRequest
	mu        sync.RWMutex
}

func NewHandler() *Handler {
	return &Handler{
		Requests:  make(map[string]ProxyRequest),
		Responses: make(map[string]ProxyRequest),
		mu:        sync.RWMutex{},
	}
}

func (h *Handler) InitRoutes() error {
	http.HandleFunc("/", nil)
	return http.ListenAndServe(":8080", nil)
}
