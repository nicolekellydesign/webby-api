package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nicolekellydesign/webby-api/database"
)

var (
	errWrongContentType = errors.New("wrong content type")
	errWrongMethod      = errors.New("wrong method type")
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

	http.HandleFunc("/api/login", l.PerformLogin)     // POST
	http.HandleFunc("/api/logout", l.PerformLogout)   // POST
	http.HandleFunc("/api/refresh", l.RefreshSession) // POST

	http.HandleFunc("/api/users/add", l.AddUser)  // POST
	http.HandleFunc("/api/users/get", l.GetUsers) // GET

	addr := fmt.Sprintf("localhost:%d", l.Port)
	http.ListenAndServe(addr, nil)
}

func checkPreconditions(r *http.Request, method string, checkContentType bool) error {
	if r.Method != method {
		return errWrongMethod
	}

	if checkContentType && r.Header.Get("Content-Type") != "application/json" {
		return errWrongContentType
	}

	return nil
}
