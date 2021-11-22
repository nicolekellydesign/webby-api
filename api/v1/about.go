package v1

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nicolekellydesign/webby-api/entities"
)

// GetAbout fetches the about page info from a file and sends it to the client.
func (a API) GetAbout(w http.ResponseWriter, r *http.Request) {
	// Open about page info file
	path := filepath.Join(a.resourcesDir, "about-info.json")
	file, err := openOrCreate(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error opening about page file: %s\n", err.Error())
		return
	}
	defer file.Close()

	// Read the file contents
	decoder := json.NewDecoder(file)
	var ret entities.About
	if err := decoder.Decode(&ret); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error reading about page file: %s\n", err.Error())
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	encoder := json.NewEncoder(w)
	encoder.Encode(ret)
}

// UpdatePortrait updates the about page portrait file name.
func (a API) UpdatePortrait(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Decode from request body
	decoder := json.NewDecoder(r.Body)
	var update entities.About
	if err := decoder.Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error decoding JSON body in about page update request: %s\n", err.Error())
		return
	}

	// Open about page info file
	path := filepath.Join(a.resourcesDir, "about-info.json")
	b, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error reading about page file: %s\n", err.Error())
		return
	}

	// Decode the file contents
	var details entities.About
	if err := json.Unmarshal(b, &details); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error decoding about page contents: %s\n", err.Error())
		return
	}

	// Set the new value
	details.Portrait = update.Portrait

	// Write out to the file
	b2, err := json.Marshal(&details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error encoding new about info: %s\n", err.Error())
		return
	}

	if err := os.WriteFile(path, b2, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error writing about info: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}

// UpdateStatement updates the about page designer statement.
func (a API) UpdateStatement(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Decode from request body
	decoder := json.NewDecoder(r.Body)
	var update entities.About
	if err := decoder.Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error decoding JSON body in about page update request: %s\n", err.Error())
		return
	}

	// Open about page info file
	path := filepath.Join(a.resourcesDir, "about-info.json")
	b, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error reading about page file: %s\n", err.Error())
		return
	}

	// Decode the file contents
	var details entities.About
	if err := json.Unmarshal(b, &details); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error decoding about page contents: %s\n", err.Error())
		return
	}

	// Set the new value
	details.Statement = update.Statement

	// Write out to the file
	b2, err := json.Marshal(&details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error encoding new about info: %s\n", err.Error())
		return
	}

	if err := os.WriteFile(path, b2, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error writing about info: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}

// UpdateResume updates the about page resume file name.
func (a API) UpdateResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// Decode from request body
	decoder := json.NewDecoder(r.Body)
	var update entities.About
	if err := decoder.Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		a.log.Errorf("error decoding JSON body in about page update request: %s\n", err.Error())
		return
	}

	// Open about page info file
	path := filepath.Join(a.resourcesDir, "about-info.json")
	b, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error reading about page file: %s\n", err.Error())
		return
	}

	// Decode the file contents
	var details entities.About
	if err := json.Unmarshal(b, &details); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error decoding about page contents: %s\n", err.Error())
		return
	}

	// Set the new value
	details.Resume = update.Resume

	// Write out to the file
	b2, err := json.Marshal(&details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error encoding new about info: %s\n", err.Error())
		return
	}

	if err := os.WriteFile(path, b2, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		a.log.Errorf("error writing about info: %s\n", err.Error())
		return
	}

	w.WriteHeader(200)
}

// openOrCreate tries to open a file, and creates it if it doesn't exist. In the
// case that the file is created, an empty about page struct is written out.
//
// This has the same use semantics as the file create/open functions in the
// standard library.
func openOrCreate(path string) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return
		}

		// Write out empty About page details
		encoder := json.NewEncoder(file)
		if err = encoder.Encode(&entities.About{
			Portrait:  "",
			Statement: "",
			Resume:    "",
		}); err != nil {
			file.Close()
			file = nil
			return
		}

		return
	}

	return
}
