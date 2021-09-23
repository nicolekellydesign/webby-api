package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AddUser adds a new user into the database.
func (l Listener) AddUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req AddUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	users, err := l.db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make sure usernames are unique
	for _, user := range users {
		if user.Username == req.Username {
			http.Error(w, "username already exists", http.StatusBadRequest)
			return
		}
	}

	if err := l.db.AddUser(req.Username, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// GetUsers gets all of the users from the database.
func (l Listener) GetUsers(w http.ResponseWriter, r *http.Request) {
	ret, err := l.db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// RemoveUser deletes a user from the database, and invalidates any
// active sessions that the user had.
func (l Listener) RemoveUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// We don't want users to be able to delete themselves, so
	// we have to get the session to check the user ID.
	session, err := l.db.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	converted, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if session.ID == uint(converted) {
		http.Error(w, "can't delete yourself", http.StatusBadRequest)
		return
	}

	if err := l.db.RemoveUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
