package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/nicolekellydesign/webby-api/entities"
)

// AddGalleryItem handles a request to add a new gallery item.
//
// Requires a valid auth token.
func (a API) AddGalleryItem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error parsing multipart form: %s\n", err.Error())
		return
	}

	name := r.FormValue("name")

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error getting file from form: %s\n", err.Error())
		return
	}
	defer file.Close()

	fileName := name + "-thumb" + filepath.Ext(header.Filename)
	outPath := filepath.Join(a.uploadDir, fileName)
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

	embedURL := r.FormValue("embed_url")

	galleryItem := entities.GalleryItem{
		Name:        name,
		Title:       r.FormValue("title"),
		Caption:     r.FormValue("caption"),
		ProjectInfo: r.FormValue("project_info"),
		Thumbnail:   fileName,
		EmbedURL: entities.NullString{
			String: embedURL,
			Valid:  embedURL != "",
		},
	}

	if err := a.db.AddGalleryItem(galleryItem); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		a.log.Errorf("error adding gallery item to database: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}

// ChangeThumbnail handles requests to change a thumbnail for a project.
func (a API) ChangeThumbnail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error parsing multipart form: %s\n", err.Error())
		return
	}

	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error getting file from form: %s\n", err.Error())
		return
	}
	defer file.Close()

	fileName := id + "-thumb" + filepath.Ext(header.Filename)
	outPath := filepath.Join(a.uploadDir, fileName)
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

	if err := a.db.ChangeProjectThumbnail(id, fileName); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// GetProject handles a request to get a portfolio project from the database.
func (a API) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "name")

	ret, err := a.db.GetProject(id)
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(ret)
}

// UpdateProject handles requests to update a portfolio project.
func (a API) UpdateProject(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var project entities.GalleryItem
	if err := decoder.Decode(&project); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error decoding JSON body in project update request: %s\n", err.Error())
		return
	}

	if err := a.db.UpdateProject(&project); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// GetGalleryItems handles a request to get all gallery items from the
// database.
func (a API) GetGalleryItems(w http.ResponseWriter, r *http.Request) {
	ret, err := a.db.GetGalleryItems()
	if err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(GalleryResponse{
		ret,
	})
}

// RemoveGalleryItem handles a request to remove a gallery item.
//
// Requires a valid auth token.
func (a API) RemoveGalleryItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := a.db.RemoveGalleryItem(id); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func (a API) AddImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

	outPath := filepath.Join(a.uploadDir, header.Filename)
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

	if err := a.db.AddProjectImage(id, header.Filename); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		a.log.Errorf("error adding project image to database: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}

// RemoveProjectImage deletes an image for a portfolio project and removes
// it from the database.
func (a API) RemoveProjectImage(w http.ResponseWriter, r *http.Request) {
	galleryID := chi.URLParam(r, "id")
	fileName := chi.URLParam(r, "name")

	if err := a.db.RemoveProjectImage(galleryID, fileName); err != nil {
		http.Error(w, dbError, http.StatusInternalServerError)
		a.log.Errorf("error removing project image from database: %s\n", err.Error())
		return
	}

	path := filepath.Join(a.uploadDir, fileName)
	if err := os.Remove(path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error removing project image: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}
