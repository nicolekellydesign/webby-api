package server

import (
	"github.com/nicolekellydesign/webby-api/entities"
)

// GalleryResponse is sent when a client requests all of the
// gallery items from the database.
type GalleryResponse struct {
	Items []*entities.GalleryItem `json:"items"`
}

// PhotosResponse is sent when a client requests all of the
// photos from the database.
type PhotosResponse struct {
	Photos []*entities.Photo `json:"photos"`
}

// UsersResponse is sent when a request to list all useres
// is received.
type UsersResponse struct {
	Users []*entities.User `json:"users"`
}
