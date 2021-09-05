package server

import (
	"encoding/json"
	"net/http"
)

// HttpError represents an error response for an API endpoint.
type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewError returns a new HttpError.
func NewError(code int, message string) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

// Write writes the HttpError to the given ResponseWriter.
func (e HttpError) Write(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	encoder := json.NewEncoder(w)
	encoder.Encode(e)
}

// AuthResponse is sent when a login is checked, and sends
// the result of the login attempt.
type AuthResponse struct {
	Valid bool `json:"valid"`
}
