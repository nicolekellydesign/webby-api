package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nicolekellydesign/webby-api/api/v1"
	"github.com/nicolekellydesign/webby-api/database"
)

// Listener handles requests to our API endpoints.
type Listener struct {
	Port int

	db     *database.DB
	router chi.Router
}

// New creates a new HTTP listener on the given port.
func New(port int, db *database.DB) *Listener {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	return &Listener{
		Port:   port,
		db:     db,
		router: r,
	}
}

// Serve sets up our endpoint handlers and begins listening.
func (l Listener) Serve() {
	api := v1.NewAPI(l.db)

	l.router.Mount("/api/v1", api.Routes())

	addr := fmt.Sprintf("localhost:%d", l.Port)
	http.ListenAndServe(addr, l.router)
}
