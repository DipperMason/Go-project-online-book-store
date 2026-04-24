package http

import (
	"encoding/json"
	"net/http"
)

// Response обертка для стандартных ответов API
type Response struct {
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// ErrorResponse стандартный ответ об ошибке
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// SuccessResponse стандартный успешный ответ
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON пишет JSON ответ
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WriteError пишет ошибку в формате JSON
func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, ErrorResponse{
		Error:  message,
		Status: status,
	})
}

// WriteSuccess пишет успех в формате JSON
func WriteSuccess(w http.ResponseWriter, message string, data interface{}) error {
	return WriteJSON(w, http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}
