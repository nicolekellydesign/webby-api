package server

import (
	"encoding/json"
	"net/http"

	"github.com/nicolekellydesign/webby-api/entities"
)

// HTTPError represents an error response for an API endpoint.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewError returns a new HttpError.
func NewError(code int, message string) HTTPError {
	return HTTPError{
		Code:    code,
		Message: message,
	}
}

// SendErrMalformedBody creates and sends an error for a malformed JSON body.
func SendErrMalformedBody(w http.ResponseWriter) {
	httpError := NewError(400, "Bad request: malformed body")
	httpError.Write(w)
}

// SendErrWrongMethod creates and sends an error for wrong HTTP method.
func SendErrWrongMethod(w http.ResponseWriter) {
	httpError := NewError(400, "Bad request: wrong HTTP method")
	httpError.Write(w)
}

// SendErrWrongType creates and sends an error for wrong content type.
func SendErrWrongType(w http.ResponseWriter) {
	httpError := NewError(400, "Bad request: wrong content type")
	httpError.Write(w)
}

// Write writes the HttpError to the given ResponseWriter.
func (e HTTPError) Write(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	encoder := json.NewEncoder(w)
	encoder.Encode(e)
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
