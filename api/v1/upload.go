package v1

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (a API) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(8 * 1024 * 1024); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error parsing multipart form: %s\n", err.Error())
		return
	}

	file, header, err := r.FormFile("file")
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
		a.log.Errorf("error creating new file: %s\n", err.Error())
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
