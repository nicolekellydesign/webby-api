package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddPhoto handles a request to add a file name to the
// photos database.
//
// It requires a valid auth token.
func (l Listener) AddPhoto(w http.ResponseWriter, r *http.Request) {
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

	if err := checkPreconditions(r, http.MethodPost, true); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	var req AddPhotoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if err := l.db.AddPhoto(req.Filename); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// GetPhotos handles requests to get all photos from the database.
func (l Listener) GetPhotos(w http.ResponseWriter, r *http.Request) {
	if err := checkPreconditions(r, http.MethodGet, false); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ret, err := l.db.GetPhotos()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(PhotosResponse{
		ret,
	})
}

// RemovePhoto handles a request to remove a file name from the
// photos database.
//
// It requires a valid auth token.
func (l Listener) RemovePhoto(w http.ResponseWriter, r *http.Request) {
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

	if err := checkPreconditions(r, http.MethodPost, true); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()

	var req RemovePhotoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if err := l.db.RemovePhoto(req.Filename); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}
