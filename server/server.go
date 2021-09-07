package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicolekellydesign/webby-api/database"
)

// Listener handles requests to our API endpoints.
type Listener struct {
	Port int

	db *database.DB
}

// New creates a new HTTP listener on the given port.
func New(port int, db *database.DB) *Listener {
	return &Listener{
		Port: port,
		db:   db,
	}
}

// Serve sets up our endpoint handlers and begins listening.
func (l Listener) Serve() {
	http.HandleFunc("/api/photos/add", l.AddPhoto)       // POST
	http.HandleFunc("/api/photos/get", l.GetPhotos)      // GET
	http.HandleFunc("/api/photos/remove", l.RemovePhoto) // POST

	http.HandleFunc("/api/gallery/add", l.AddGalleryItem)       // POST
	http.HandleFunc("/api/gallery/get", l.GetGalleryItems)      // GET
	http.HandleFunc("/api/gallery/remove", l.RemoveGalleryItem) // POST

	http.HandleFunc("/api/gallery/slides/add", l.AddSlide)       // POST
	http.HandleFunc("/api/gallery/slides/remove", l.RemoveSlide) // POST

	http.HandleFunc("/api/login", l.PerformLogin) // POST
	http.HandleFunc("/api/users/add", l.AddUser)  // POST
	http.HandleFunc("/api/users/get", l.GetUsers) // GET

	addr := fmt.Sprintf("localhost:%d", l.Port)
	http.ListenAndServe(addr, nil)
}

// AddPhoto handles a request to add a file name to the
// photos database.
//
// It requires a valid auth token.
func (l Listener) AddPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req AddPhotoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if _, _, err := validateToken(req.Token); err != nil {
		if err == errInvalidToken {
			WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
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
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
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
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req RemovePhotoRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if _, _, err := validateToken(req.Token); err != nil {
		if err == errInvalidToken {
			WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	if err := l.db.RemovePhoto(req.Filename); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// AddGalleryItem handles a request to add a new gallery item.
//
// Requires a valid auth token.
func (l Listener) AddGalleryItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req AddGalleryItemRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if _, _, err := validateToken(req.Token); err != nil {
		if err == errInvalidToken {
			WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
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
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
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
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req RemoveGalleryItemRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if _, _, err := validateToken(req.Token); err != nil {
		if err == errInvalidToken {
			WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
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
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req AddSlideRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if _, _, err := validateToken(req.Token); err != nil {
		if err == errInvalidToken {
			WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
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
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req RemoveSlideRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	if _, _, err := validateToken(req.Token); err != nil {
		if err == errInvalidToken {
			WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	if err := l.db.RemoveSlide(req.GalleryID, req.Name); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// AddUser adds a new user into the database.
func (l Listener) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	var req AddUserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	users, err := l.db.GetUsers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// If no users have been added yet, then it would be impossible to
	// get a token in order to add one. So, only check for a valid token
	// if there are any users in the database.
	if len(users) > 0 {
		if _, _, err := validateToken(req.Token); err != nil {
			if err == errInvalidToken {
				WriteError(w, http.StatusUnauthorized, err.Error())
				return
			}

			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
			return
		}

		// Make sure usernames are unique
		for _, user := range users {
			if user.Username == req.Username {
				WriteError(w, http.StatusBadRequest, "username already exists")
				return
			}
		}
	}

	if err := l.db.AddUser(req.Username, req.Password); err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	w.WriteHeader(200)
}

// PerformLogin checks if the given credentials match, and if so, generates
// and responds with an auth token.
func (l Listener) PerformLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		WriteError(w, http.StatusBadRequest, "wrong content type")
		return
	}

	defer r.Body.Close()

	// Decode the user that was sent
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	// Check if the login should be a success
	user, err := l.db.GetLogin(req.Username, req.Password)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// Send back the response
	if user != nil {
		token, err := generateToken(user)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		encoder := json.NewEncoder(w)
		encoder.Encode(AuthResponse{
			token,
		})
	} else {
		WriteError(w, http.StatusUnauthorized, "incorrect login credentials")
	}
}

// GetUsers gets all of the users from the database.
func (l Listener) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusBadRequest, "wrong HTTP method type")
		return
	}

	ret, err := l.db.GetUsers()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Internal error: %s", err.Error()))
		return
	}

	// Send back the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	encoder := json.NewEncoder(w)
	encoder.Encode(UsersResponse{
		ret,
	})
}
