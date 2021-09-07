package server

import (
	"encoding/json"
	"net/http"

	"github.com/nicolekellydesign/webby-api/entities"
)

// httpError represents an error response for an API endpoint.
type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteError sends an error response with the code and error message
// in the body encoded as JSON.
func WriteError(w http.ResponseWriter, code int, message string) {
	err := httpError{
		code,
		message,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	encoder := json.NewEncoder(w)
	encoder.Encode(err)
}

// AuthResponse is sent when a login is checked, and sends
// the result of the login attempt.
type AuthResponse struct {
	Token string `json:"token"`
}

// UsersResponse is sent when a request to list all useres
// is received.
type UsersResponse struct {
	Users []*entities.User `json:"users"`
}
