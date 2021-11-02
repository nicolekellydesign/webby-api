package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nicolekellydesign/webby-api/entities"
)

// ChangePortrait handles requests to change the about page portrait.
func (a API) ChangePortrait(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error parsing multipart form: %s\n", err.Error())
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error getting file from form: %s\n", err.Error())
		return
	}
	defer file.Close()

	outPath := filepath.Join(a.imageDir, "about-portrait.jpg")
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

	w.WriteHeader(200)
}

// ChangeResume handles requests to change the about page resume.
func (a API) ChangeResume(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error parsing multipart form: %s\n", err.Error())
		return
	}

	file, _, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error getting file from form: %s\n", err.Error())
		return
	}
	defer file.Close()

	outPath := filepath.Join(a.resourcesDir, "resume.pdf")
	out, err := os.Create(outPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error creating resume file: %s\n", err.Error())
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error copying to file: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}

// GetAbout fetches the about page info from a file and sends it to the client.
func (a API) GetAbout(w http.ResponseWriter, r *http.Request) {
	// Open about page info file
	path := filepath.Join(a.resourcesDir, "about-info.json")

	// Check if the file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// File hasn't been written yet, so return a blank statement
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(200)
			encoder := json.NewEncoder(w)
			encoder.Encode(&entities.About{
				Statement: "",
			})
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("unable to stat about page file: %s\n", err.Error())
		return
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error opening about page file: %s\n", err.Error())
		return
	}
	defer file.Close()

	// Read the file contents
	decoder := json.NewDecoder(file)
	var ret entities.About
	if err := decoder.Decode(&ret); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error reading about page file: %s\n", err.Error())
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	encoder := json.NewEncoder(w)
	encoder.Encode(ret)
}

// UpdateAbout writes the information in the request body to a file.
func (a API) UpdateAbout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Decode from request body
	decoder := json.NewDecoder(r.Body)
	var about entities.About
	if err := decoder.Decode(&about); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error decoding JSON body in about page update request: %s\n", err.Error())
		return
	}

	// Open about page info file
	path := filepath.Join(a.resourcesDir, "about-info.json")
	file, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error opening about page file: %s\n", err.Error())
		return
	}
	defer file.Close()

	// Write out to the file
	encoder := json.NewEncoder(file)

	if err := encoder.Encode(&about); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error writing about info: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}
