package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPError holds information about an error that will be sent to the client in
// an HTTP response.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteError writes an error struct to a ResponseWriter as JSON.
func WriteError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := HTTPError{
		Code:    code,
		Message: message,
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&e); err != nil {
		http.Error(w, fmt.Sprintf("error sending previous error: %s; previous error: %s", err.Error(), e.Message), 500)
		return
	}
}
