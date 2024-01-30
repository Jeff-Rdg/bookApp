package httpResponse

import (
	"encoding/json"
	"net/http"
)

type ApiResponse interface {
	RenderJSON(w http.ResponseWriter)
}

type Success struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"response,omitempty"`
}

type Created struct {
	Message string `json:"message,omitempty"`
}

func NewCreatedResponse(message string) *Created {
	return &Created{Message: message}
}

func NewSuccessResponse(data interface{}) *Success {
	return &Success{Data: data}
}

func (s *Created) RenderJSON(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (s *Success) RenderJSON(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.Data)
}
