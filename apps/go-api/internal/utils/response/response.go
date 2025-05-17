package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Error   interface{} `json:"error",omitempty`
	Message interface{} `json:"message,omitempty"`
}

const (
	Success = "success"
	Error   = "error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := Response{
		Status: Error,
		Error:  err.Error(),
	}
	return json.NewEncoder(w).Encode(response)
}
