package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Response struct {
	Message  string      `json:"message,omitempty"`
	Response interface{} `json:"response,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

type Result struct {
	StatusCode int
	Message    string
	Data       interface{}
}

func RenderJSON(w http.ResponseWriter, result Result) {
	response := switchResponse(result.Data, result.Message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(result.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func switchResponse(data interface{}, message string) Response {
	var result interface{}
	switch data.(type) {
	case error:
		value, ok := data.(error)
		if ok {
			result = errorsToList(value)
		} else {
			result = data
		}
		return Response{
			Message: message,
			Error:   result,
		}
	default:
		result = data
		return Response{
			Message:  message,
			Response: result,
		}
	}
}

func errorsToList(err error) []string {
	var errorList []string

	for _, msg := range strings.Split(err.Error(), "\n") {
		errorList = append(errorList, msg)
	}

	return errorList
}
