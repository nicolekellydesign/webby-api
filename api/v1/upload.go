package v1

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Upload handles requests to upload files to the server.
func (a API) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error parsing multipart form: %s\n", err.Error())
		return
	}

	// Get the file from the body
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error getting file from form: %s\n", err.Error())
		return
	}
	defer file.Close()

	// Determine which path the file should be saved to
	var outPath string
	if strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		outPath = filepath.Join(a.imageDir, header.Filename)
	} else {
		outPath = filepath.Join(a.resourcesDir, header.Filename)
	}

	// Create our out file
	out, err := os.Create(outPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error creating new file: %s\n", err.Error())
		return
	}
	defer out.Close()

	// Copy the file from the requset to the out file
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error copying to file: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}
