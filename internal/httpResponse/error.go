package httpResponse

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Title      string  `json:"title"`
	Details    string  `json:"details,omitempty"`
	StatusCode int     `json:"-"`
	Error      []Cause `json:"error,omitempty"`
	Instance   string  `json:"instance"`
}

type Cause struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewBadRequestError(message, details string, r *http.Request) *Error {
	return &Error{
		Title:      message,
		Details:    details,
		StatusCode: http.StatusBadRequest,
		Error:      nil,
		Instance:   r.URL.Path,
	}
}

func NewInternalServerError(message, details string, r *http.Request) *Error {
	return &Error{
		Title:      message,
		Details:    details,
		StatusCode: http.StatusInternalServerError,
		Error:      nil,
		Instance:   r.URL.Path,
	}
}

func NewBadRequestValidationError(message string, cause []Cause, r *http.Request) *Error {
	return &Error{
		Title:      message,
		StatusCode: http.StatusBadRequest,
		Error:      cause,
		Instance:   r.URL.Path,
	}
}

func NewNotFoundError(message, details string, r *http.Request) *Error {
	return &Error{
		Title:      message,
		Details:    details,
		StatusCode: http.StatusNotFound,
		Instance:   r.URL.Path,
	}
}

func (e *Error) RenderJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(e.StatusCode)
	json.NewEncoder(w).Encode(e)
}
