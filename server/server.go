package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicolekellydesign/webby-api/database"
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
	http.HandleFunc("/api/adduser", l.AddUser)    // POST
	http.HandleFunc("/api/login", l.PerformLogin) // POST
	http.HandleFunc("/api/users", l.GetUsers)     // GET

	addr := fmt.Sprintf("localhost:%d", l.Port)
	http.ListenAndServe(addr, nil)
}

// AddUser adds a new user into the database.
func (l Listener) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req AddUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	users, err := l.db.GetUsers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// If no users have been added yet, then it would be impossible to
	// get a token in order to add one. So, only check for a valid token
	// if there are any users in the database.
	if len(users) > 0 {
		if _, _, err := validateToken(req.Token); err != nil {
			if err == errInvalidToken {
				WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
			return
		}
	}

	if err := l.db.AddUser(req.Username, req.Password); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// PerformLogin checks if the given credentials match, and if so, generates
// and responds with an auth token.
func (l Listener) PerformLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	// Decode the user that was sent
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	// Check if the login should be a success
	user, err := l.db.GetLogin(req.Username, req.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// Send back the response
	if user != nil {
		token, err := generateToken(user)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		encoder := json.NewEncoder(w)
		encoder.Encode(AuthResponse{
			token,
		})
	} else {
		WriteError(w, http.StatusUnauthorized, "incorrect login credentials")
	}
}

// GetUsers gets all of the users from the database.
func (l Listener) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	ret, err := l.db.GetUsers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(UsersResponse{
		ret,
	})
}
