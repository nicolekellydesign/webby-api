package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// AddGalleryItem handles a request to add a new gallery item.
//
// Requires a valid auth token.
func (l Listener) AddGalleryItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req AddGalleryItemRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := l.db.AddGalleryItem(req.Item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// GetGalleryItems handles a request to get all gallery items from the
// database.
func (l Listener) GetGalleryItems(w http.ResponseWriter, r *http.Request) {
	ret, err := l.db.GetGalleryItems()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	id := chi.URLParam(r, "id")
	if err := l.db.RemoveGalleryItem(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// AddSlide handles a request to add a new gallery slide.
//
// Requires a valid auth token.
func (l Listener) AddSlide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req AddSlideRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	req.Slide.GalleryID = id
	if err := l.db.AddSlide(req.Slide); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// RemoveSlide handles a request to remove a gallery slide.
//
// Requires a valid auth token.
func (l Listener) RemoveSlide(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	name := chi.URLParam(r, "name")
	if err := l.db.RemoveSlide(id, name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
