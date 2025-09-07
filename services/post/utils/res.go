package utils

import (
	"encoding/json"
	"net/http"
)

const (
	Success = "success"
	Error   = "error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Res(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func ResSuccess(message string) *Response {
	return &Response{
		Status:  Success,
		Message: message,
	}
}

func ResError(message string) *Response {
	return &Response{
		Status:  Error,
		Message: message,
	}
}
