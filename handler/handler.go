package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
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
	Length  int               `json:"length"`
}

type Handler struct {
	Requests  map[string]ProxyRequest
	Responses map[string]ProxyResponse
}

func NewHandler() *Handler {
	return &Handler{
		Requests:  make(map[string]ProxyRequest),
		Responses: make(map[string]ProxyResponse),
	}
}

func (h *Handler) InitRoutes() error {
	http.HandleFunc("/", h.proxyHandler)
	return http.ListenAndServe(":8080", nil)
}

func (h *Handler) proxyHandler(w http.ResponseWriter, r *http.Request) {
	var input ProxyRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("status code: %d\nerror body: %s\n", http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	if input.Method == "" ||
		input.Headers == nil ||
		input.URL == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("status code: %d\nerror body: %s\n", http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	id := uuid.NewString()
	h.Requests[id] = input

	req, err := http.NewRequest(input.Method, input.URL, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Printf("status code: %d\nerror body: %s\n", http.StatusBadRequest, err.Error())
		return
	}

	for k, v := range input.Headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("status code: %d\nerror body: %s\n", http.StatusInternalServerError, err.Error())
		return
	}

	defer resp.Body.Close()

	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		respHeaders[k] = v[0]
	}

	repsBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("status code: %d\nerror body: %s\n", http.StatusInternalServerError, err.Error())
		return
	}

	newResponse := ProxyResponse{
		ID:      id,
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Length:  len(repsBody),
	}

	h.Responses[id] = newResponse

	if err := json.NewEncoder(w).Encode(&newResponse); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("status code: %d\nerror body: %s\n", http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("status code: %d\nerror body: %v\n", http.StatusOK, nil)
}
