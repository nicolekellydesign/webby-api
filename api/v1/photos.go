package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// AddPhoto handles a request to add a file name to the
// photos database.
//
// It requires a valid auth token.
func (a API) AddPhoto(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error parsing multipart form: %s\n", err.Error())
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error getting file from form: %s\n", err.Error())
		return
	}
	defer file.Close()

	outPath := filepath.Join(a.imageDir, header.Filename)
	out, err := os.Create(outPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error creating new image file: %s\n", err.Error())
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error copying to file: %s\n", err.Error())
		return
	}

	if err := a.db.AddPhoto(header.Filename); err != nil {
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(PhotosResponse{
		ret,
	})
}

// RemovePhotos handles a request to remove a list of files from the
// photos database.
//
// It requires a valid auth token.
func (a API) RemovePhotos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var files []string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&files); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error decoding image files to remove: %s\n", err.Error())
		return
	}

	if err := a.db.RemovePhotos(files); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		path := filepath.Join(a.imageDir, file)
		if err := os.Remove(path); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			a.log.Errorf("error removing image: %s\n", err.Error())
			return
		}
	}

	w.WriteHeader(200)
}
