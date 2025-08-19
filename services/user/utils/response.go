package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type StatusCode int

const (
	Success = "success"
	Error = "error"
)

func Res(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func GenerateResponse(message string, status string) Response {
	return Response{
		Message: message,
		Status:  status,
	}
}
