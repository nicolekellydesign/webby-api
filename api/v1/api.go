package v1

import (
	"net/http"
	"time"

	"github.com/DataDrake/waterlog"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nicolekellydesign/webby-api/database"
)

const dbError = "internal database error"

// API is our v1 API that serves and handles endpoints.
type API struct {
	db           *database.DB
	log          *waterlog.WaterLog
	imageDir     string
	resourcesDir string
}

// NewAPI creates a new v1 API.
func NewAPI(db *database.DB, log *waterlog.WaterLog, imagesDir, resourcesDir string) *API {
	return &API{
		db,
		log,
		imagesDir,
		resourcesDir,
	}
}

// Routes sets up our API routes.
func (a API) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/about", a.GetAbout)
	r.Get("/photos", a.GetPhotos)
	r.Get("/gallery", a.GetGalleryItems)
	r.Get("/gallery/{name}", a.GetProject)

	r.Get("/check", a.CheckSession)
	r.Post("/login", a.PerformLogin)
	r.Post("/logout", a.PerformLogout)

	r.Mount("/admin", a.adminRouter())

	return r
}

// adminRouter sets up the admin-only API routes.
func (a API) adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(a.adminOnly)
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))

	r.Route("/about", func(r chi.Router) {
		r.Patch("/", a.UpdateAbout)
	})

	r.Route("/gallery", func(r chi.Router) {
		r.Post("/", a.AddGalleryItem)

		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", a.RemoveGalleryItem)
			r.Put("/", a.UpdateProject)

			r.Patch("/thumbnail", a.ChangeThumbnail)

			r.Post("/images", a.AddImages)
			r.Delete("/images", a.RemoveProjectImages)
		})
	})

	r.Post("/photos", a.AddPhotos)
	r.Delete("/photos", a.RemovePhotos)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", a.GetUsers)
		r.Post("/", a.AddUser)
		r.Delete("/{id}", a.RemoveUser)
	})

	r.Post("/upload", a.Upload)

	return r
}

// adminOnly returns a middleware handler to check for a valid session.
func (a API) adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get our session cookie if we have one
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				WriteError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get our stored session
		token := cookie.Value
		session, err := a.db.GetSession(token)
		if err != nil {
			WriteError(w, dbError, http.StatusInternalServerError)
			return
		}

		// Check if we have a session
		if session == nil {
			WriteError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Check if the session has expired.
		// If it has expired, remove it from the database.
		if session.MaxAge > 0 {
			expires := session.Created.Add(time.Duration(session.MaxAge) * time.Second)
			if time.Now().After(expires) {
				if err = a.db.RemoveSession(token); err != nil {
					WriteError(w, dbError, http.StatusInternalServerError)
					return
				}

				WriteError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
