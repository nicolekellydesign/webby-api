package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddGalleryItem handles a request to add a new gallery item.
//
// Requires a valid auth token.
func (l Listener) AddGalleryItem(w http.ResponseWriter, r *http.Request) {
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

	var req AddGalleryItemRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if err := l.db.AddGalleryItem(req.Item); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// GetGalleryItems handles a request to get all gallery items from the
// database.
func (l Listener) GetGalleryItems(w http.ResponseWriter, r *http.Request) {
	if err := checkPreconditions(r, http.MethodGet, false); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	ret, err := l.db.GetGalleryItems()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(GalleryResponse{
		ret,
	})
}

// RemoveGalleryItem handles a request to remove a gallery item.
//
// Requires a valid auth token.
func (l Listener) RemoveGalleryItem(w http.ResponseWriter, r *http.Request) {
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

	var req RemoveGalleryItemRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if err := l.db.RemoveGalleryItem(req.ID); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// AddSlide handles a request to add a new gallery slide.
//
// Requires a valid auth token.
func (l Listener) AddSlide(w http.ResponseWriter, r *http.Request) {
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

	var req AddSlideRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if err := l.db.AddSlide(req.Slide); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// RemoveSlide handles a request to remove a gallery slide.
//
// Requires a valid auth token.
func (l Listener) RemoveSlide(w http.ResponseWriter, r *http.Request) {
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

	var req RemoveSlideRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if err := l.db.RemoveSlide(req.GalleryID, req.Name); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}
