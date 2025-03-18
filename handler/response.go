package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type UnprocessableEntity struct {
	Errors []string `json:"errors"`
}

func (e *UnprocessableEntity) Add(error string) {
	e.Errors = append(e.Errors, error)
}

func (e *UnprocessableEntity) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *UnprocessableEntity) Error() string {
	return strings.Join(e.Errors, "; ")
}

func Response(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&response)
}

func ResponseError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&ErrorResponse{Message: err.Error()})
}
