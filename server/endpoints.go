package server

import (
	"encoding/json"
	"net/http"
)

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	resp := struct {
		Message string
	}{
		Message: "Hello!",
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}
