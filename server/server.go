package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/DataDrake/waterlog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/nicolekellydesign/webby-api/api/v1"
	"github.com/nicolekellydesign/webby-api/database"
)

// Listener handles requests to our API endpoints.
type Listener struct {
	Port int

	db           *database.DB
	log          *waterlog.WaterLog
	router       chi.Router
	rootDir      string
	imagesDir    string
	resourcesDir string

	errs chan error
}

// New creates a new HTTP listener on the given port.
func New(port int, db *database.DB, log *waterlog.WaterLog, rootDir string, errs chan error) *Listener {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	return &Listener{
		Port:         port,
		db:           db,
		log:          log,
		router:       r,
		rootDir:      rootDir,
		imagesDir:    filepath.Join(rootDir, "images"),
		resourcesDir: filepath.Join(rootDir, "resources"),
		errs:         errs,
	}
}

// Serve sets up our endpoint handlers and begins listening.
func (l Listener) Serve() {
	if err := os.MkdirAll(l.rootDir, 0755); err != nil {
		l.errs <- fmt.Errorf("root dir does not exist and could not create it: %s", err.Error())
	}

	if err := os.MkdirAll(l.imagesDir, 0755); err != nil {
		l.errs <- fmt.Errorf("images dir does not exist and could not create it: %s", err.Error())
	}

	if err := os.MkdirAll(l.resourcesDir, 0755); err != nil {
		l.errs <- fmt.Errorf("resources dir does not exist and could not create it: %s", err.Error())
	}

	api := v1.NewAPI(l.db, l.log, l.imagesDir, l.resourcesDir)
	l.router.Mount("/api/v1", api.Routes())

	addr := fmt.Sprintf("localhost:%d", l.Port)
	l.errs <- http.ListenAndServe(addr, l.router)
}
