package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddUser adds a new user into the database.
func (l Listener) AddUser(w http.ResponseWriter, r *http.Request) {
	if err := checkPreconditions(r, http.MethodPost, true); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
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
		ok, code, err := checkSession(l.db, r)
		if !ok {
			if err != nil {
				WriteError(w, code, err.Error())
				return
			}

			// Session was not okay, but no error
			// That means the session is not valid
			w.WriteHeader(code)
			return
		}

		// Make sure usernames are unique
		for _, user := range users {
			if user.Username == req.Username {
				WriteError(w, http.StatusBadRequest, "username already exists")
				return
			}
		}
	}

	if err := l.db.AddUser(req.Username, req.Password); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// GetUsers gets all of the users from the database.
func (l Listener) GetUsers(w http.ResponseWriter, r *http.Request) {
	if err := checkPreconditions(r, http.MethodGet, false); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
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
