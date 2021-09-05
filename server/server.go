package server

import (
	"fmt"
	"net/http"
)

// Listener handles requests to our API endpoints
type Listener struct {
	Port int
}

func (l Listener) Serve() {
	http.HandleFunc("/test", TestEndpoint)

	addr := fmt.Sprintf("localhost:%d", l.Port)
	http.ListenAndServe(addr, nil)
}
