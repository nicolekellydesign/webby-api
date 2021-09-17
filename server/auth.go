package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicolekellydesign/webby-api/database"
	"github.com/nicolekellydesign/webby-api/entities"
)

// checkSession tries to get our session cookie and checks if it is
// valid, meaning it exists in the database and the expiration time
// has not yet passed.
func checkSession(db *database.DB, r *http.Request) (ok bool, code int, err error) {
	ok = false

	// Get our session cookie if we have one
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			code = http.StatusUnauthorized
			return
		}

		code = http.StatusBadRequest
		return
	}

	// Get our stored session
	token := cookie.Value
	session, err := db.GetSession(token)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	// Check if we have a session
	if session == nil {
		code = http.StatusUnauthorized
		return
	}

	// Check if the session has expired.
	// If it has expired, remove it from the database.
	if time.Now().After(session.Expires) {
		if err = db.RemoveSession(token); err != nil {
			code = http.StatusInternalServerError
			return
		}

		code = http.StatusUnauthorized
		return
	}

	// Everything passed, so we have a valid session
	ok = true
	code = http.StatusOK
	return
}

// PerformLogin checks if the given credentials match, and if so, generates
// and responds with an auth token.
func (l Listener) PerformLogin(w http.ResponseWriter, r *http.Request) {
	if err := checkPreconditions(r, http.MethodPost, true); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
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
	if user != nil && (*user != entities.User{}) {
		session, err := entities.NewSession(user.Username, req.Extended)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
			return
		}

		// Remove any existing session for this user from the database
		if err := l.db.RemoveSessionForName(user.Username); err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
			return
		}

		if err = l.db.AddSession(session); err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.Token,
			Expires:  session.Expires,
			HttpOnly: true,
		})
	} else {
		WriteError(w, http.StatusUnauthorized, "incorrect login credentials")
	}
}

// PerformLogout handles when a user wants to log out of their session.
func (l Listener) PerformLogout(w http.ResponseWriter, r *http.Request) {
	if err := checkPreconditions(r, http.MethodPost, false); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get our session cookie if we have one
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get our stored session
	token := cookie.Value
	session, err := l.db.GetSession(token)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if we have a session
	if session == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the existing session has expired
	if session.Expires.Before(time.Now()) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = l.db.RemoveSession(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// RefreshSession handles requests to refresh a session token.
func (l Listener) RefreshSession(w http.ResponseWriter, r *http.Request) {
	if err := checkPreconditions(r, http.MethodPost, false); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get our session cookie if we have one
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get our stored session
	token := cookie.Value
	session, err := l.db.GetSession(token)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if we have a session
	if session == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the existing session has expired
	if session.Expires.Before(time.Now()) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Update the session to be valid for another 5 mins
	session.Expires = time.Now().Add(300 * time.Second).UTC()
	if err = l.db.UpdateSession(session); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// Re-set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session.Token,
		Expires: session.Expires,
	})
}
