package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nicolekellydesign/webby-api/entities"
)

// CheckSession checks if the request has a valid session.
func (a API) CheckSession(w http.ResponseWriter, r *http.Request) {
	// Get our session cookie if we have one
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(&CheckSessionResponse{Valid: false})
		return
	}

	// Get our stored session
	token := cookie.Value
	session, err := a.db.GetSession(token)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	// Check if we have a session
	if session == nil || (*session == entities.Session{}) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.Encode(&CheckSessionResponse{Valid: false})
		return
	}

	// Check if the session has expired.
	// If it has expired, remove it from the database.
	if session.MaxAge > 0 {
		expires := session.Created.Add(time.Duration(session.MaxAge) * time.Second)
		if time.Now().After(expires) {
			if err = a.db.RemoveSession(token); err != nil {
				http.Error(w, dbError, http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			encoder := json.NewEncoder(w)
			encoder.Encode(&CheckSessionResponse{Valid: false})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(&CheckSessionResponse{Valid: true})
}

// PerformLogin checks if the given credentials match, and if so, generates
// and responds with an auth token.
func (a API) PerformLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Decode the user that was sent
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check if the login should be a success
	user, err := a.db.GetLogin(req.Username, req.Password)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	// Send back the response
	if user != nil && (*user != entities.User{}) {
		session, err := entities.NewSession(user.ID, user.Username, req.Extended)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Remove any existing session for this user from the database
		if err := a.db.RemoveSessionForName(user.Username); err != nil {
			http.Error(w, dbError, http.StatusInternalServerError)
			return
		}

		if err = a.db.AddSession(session); err != nil {
			http.Error(w, dbError, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.Token,
			Path:     "/",
			MaxAge:   session.MaxAge,
			HttpOnly: true,
			Secure:   true,
		})
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

// PerformLogout handles when a user wants to log out of their session.
func (a API) PerformLogout(w http.ResponseWriter, r *http.Request) {
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
	session, err := a.db.GetSession(token)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	// Check if we have a session
	if session == nil || (*session == entities.Session{}) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if the existing session has expired
	if session.MaxAge > 0 {
		expires := session.Created.Add(time.Duration(session.MaxAge) * time.Second)
		if time.Now().After(expires) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	if err = a.db.RemoveSession(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	})
	w.WriteHeader(http.StatusOK)
}
