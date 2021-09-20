package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nicolekellydesign/webby-api/entities"
)

// PerformLogin checks if the given credentials match, and if so, generates
// and responds with an auth token.
func (l Listener) PerformLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Decode the user that was sent
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check if the login should be a success
	user, err := l.db.GetLogin(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back the response
	if user != nil && (*user != entities.User{}) {
		session, err := entities.NewSession(user.Username, req.Extended)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Remove any existing session for this user from the database
		if err := l.db.RemoveSessionForName(user.Username); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = l.db.AddSession(session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.Token,
			Expires:  session.Expires,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

// PerformLogout handles when a user wants to log out of their session.
func (l Listener) PerformLogout(w http.ResponseWriter, r *http.Request) {
	// Get our session cookie if we have one
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get our stored session
	token := cookie.Value
	session, err := l.db.GetSession(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
}

// RefreshSession handles requests to refresh a session token.
func (l Listener) RefreshSession(w http.ResponseWriter, r *http.Request) {
	// Get our session cookie if we have one
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get our stored session
	token := cookie.Value
	session, err := l.db.GetSession(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Re-set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session.Token,
		Expires: session.Expires,
	})
	w.WriteHeader(http.StatusOK)
}