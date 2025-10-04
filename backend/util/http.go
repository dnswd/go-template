package util

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
    Error string `json:"error"`
}

var (
    ErrNotFound      = errors.New("not found")
    ErrBadRequest    = errors.New("bad request")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrForbidden     = errors.New("forbidden")
    ErrConflict      = errors.New("conflict")
)

// HTTPError converts error to appropriate HTTP status
func HTTPError(w http.ResponseWriter, err error) {
    var status int
    switch {
    case errors.Is(err, ErrNotFound):
        status = http.StatusNotFound
    case errors.Is(err, ErrBadRequest):
        status = http.StatusBadRequest
    case errors.Is(err, ErrUnauthorized):
        status = http.StatusUnauthorized
    case errors.Is(err, ErrForbidden):
        status = http.StatusForbidden
    case errors.Is(err, ErrConflict):
        status = http.StatusConflict
    default:
        status = http.StatusInternalServerError
    }
    
    Error(w, status, err)
}

// JSON writes a JSON response
func JSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

// Error writes an error JSON response
func Error(w http.ResponseWriter, status int, err error) {
    JSON(w, status, ErrorResponse{Error: err.Error()})
}

// ErrorMsg writes an error JSON response with a message
func ErrorMsg(w http.ResponseWriter, status int, message string) {
    JSON(w, status, ErrorResponse{Error: message})
}