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

func RenderJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	var result interface{}
	var response Response

	switch data.(type) {
	case error:
		value, ok := data.(error)
		if ok {
			result = errorsToList(value)
		} else {
			result = data
		}
		response = Response{
			Message: message,
			Error:   result,
		}
	default:
		result = data
		response = Response{
			Message:  message,
			Response: result,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(response)
}

func errorsToList(err error) []string {
	var errorList []string

	for _, msg := range strings.Split(err.Error(), "\n") {
		errorList = append(errorList, msg)
	}

	return errorList
}
