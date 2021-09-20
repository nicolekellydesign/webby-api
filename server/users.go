package server

import (
	"encoding/json"
	"net/http"
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
