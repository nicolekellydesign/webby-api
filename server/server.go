package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	l.router.Route("/api", func(r chi.Router) {
		r.Get("/photos", l.GetPhotos)
		r.Get("/gallery", l.GetGalleryItems)

		r.Post("/login", l.PerformLogin)
		r.Post("/logout", l.PerformLogout)
		r.Post("/refresh", l.RefreshSession)

		l.router.Mount("/admin", l.adminRouter())
	})

	addr := fmt.Sprintf("localhost:%d", l.Port)
	http.ListenAndServe(addr, l.router)
}

func (l Listener) adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(l.adminOnly)
	r.Use(middleware.AllowContentType("application/json"))

	r.Route("/gallery", func(r chi.Router) {
		r.Post("/", l.AddGalleryItem)

		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", l.RemoveGalleryItem)

			r.Post("/slides", l.AddSlide)
			r.Delete("/slides/{name}", l.RemoveSlide)
		})
	})

	r.Post("/photos", l.AddPhoto)
	r.Delete("/photos/{fileName}", l.RemovePhoto)

	r.Get("/users", l.GetUsers)
	r.Post("/users", l.AddUser)

	return r
}

func (l Listener) adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get our session cookie if we have one
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get our stored session
		token := cookie.Value
		session, err := l.db.GetSession(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if we have a session
		if session == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Check if the session has expired.
		// If it has expired, remove it from the database.
		if time.Now().After(session.Expires) {
			if err = l.db.RemoveSession(token); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
