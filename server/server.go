package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicolekellydesign/webby-api/database"
	"github.com/nicolekellydesign/webby-api/entities"
)

// Listener handles requests to our API endpoints.
type Listener struct {
	Port int

	db *database.DB
}

// New creates a new HTTP listener on the given port.
func New(port int, db *database.DB) *Listener {
	return &Listener{
		Port: port,
		db:   db,
	}
}

// Serve sets up our endpoint handlers and begins listening.
func (l Listener) Serve() {
	http.HandleFunc("/api/adduser", l.AddUser)  // POST
	http.HandleFunc("/api/login", l.CheckLogin) // GET

	addr := fmt.Sprintf("localhost:%d", l.Port)
	http.ListenAndServe(addr, nil)
}

// AddUser adds a new user into the database.
func (l Listener) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpError := NewError(400, "Bad request: wrong content type")
		httpError.Write(w)
		return
	}

	defer r.Body.Close()

	var user entities.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		httpError := NewError(400, "Bad request: malformed json struct")
		httpError.Write(w)
		return
	}

	if err := l.db.AddUser(&user); err != nil {
		httpError := NewError(500, fmt.Sprintf("Internal error: %s", err.Error()))
		httpError.Write(w)
		return
	}

	w.WriteHeader(200)
}

// CheckLogin looks to see if a login from the API should
// be successful or not.
func (l Listener) CheckLogin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		httpError := NewError(400, "Bad request: wrong content type")
		httpError.Write(w)
		return
	}

	defer r.Body.Close()

	// Decode the user that was sent
	var user entities.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		httpError := NewError(400, "Bad request: malformed json struct")
		httpError.Write(w)
		return
	}

	// Check if the login should be a success
	valid, err := l.db.CheckLogin(&user)
	if err != nil {
		httpError := NewError(500, fmt.Sprintf("Internal error: %s", err.Error()))
		httpError.Write(w)
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(AuthResponse{
		valid,
	})
}
