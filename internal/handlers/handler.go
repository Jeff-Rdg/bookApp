package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ApiResponse[T SuccessResponse | ErrorResponse] interface {
	RenderJSON(w http.ResponseWriter)
}

type SuccessResponse struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"response,omitempty"`
}
type ErrorResponse struct {
	Title      string   `json:"title"`
	Detail     string   `json:"detail,omitempty"`
	StatusCode int      `json:"status"`
	Error      []string `json:"error,omitempty"`
	Instance   string   `json:"instance"`
}

func MakeSuccessResponse(statusCode int, message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func MakeErrorResponse(title, detail string, status int, err error, instance string) *ErrorResponse {
	if detail == "" && err != nil {
		detail = "please, refer to the errors property for additional details"
	}
	errs := errorsToList(err)
	return &ErrorResponse{
		Title:      title,
		Detail:     detail,
		StatusCode: status,
		Error:      errs,
		Instance:   instance,
	}
}

func (s *SuccessResponse) RenderJSON(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s.StatusCode)
	json.NewEncoder(w).Encode(s)
}

func (f *ErrorResponse) RenderJSON(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(f.StatusCode)
	json.NewEncoder(w).Encode(f)
}

func errorsToList(err error) []string {
	var errorList []string

	if err != nil {
		for _, msg := range strings.Split(err.Error(), "\n") {
			errorList = append(errorList, msg)
		}
	}
	return errorList
}
