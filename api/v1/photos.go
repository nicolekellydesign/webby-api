package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// AddPhoto handles a request to add a file name to the
// photos database.
//
// It requires a valid auth token.
func (a API) AddPhoto(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req AddPhotoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := a.db.AddPhoto(req.Filename); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// GetPhotos handles requests to get all photos from the database.
func (a API) GetPhotos(w http.ResponseWriter, r *http.Request) {
	ret, err := a.db.GetPhotos()
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
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
func (a API) RemovePhoto(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	if err := a.db.RemovePhoto(fileName); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}
