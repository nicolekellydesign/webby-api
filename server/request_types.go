package server

import "github.com/nicolekellydesign/webby-api/entities"

// AddPhotoRequest holds the token and file name of a photo to add
// to the database.
type AddPhotoRequest struct {
	Token    string `json:"token"`
	Filename string `json:"filename"`
}

// RemovePhotoRequest holds the token and file name of a photo to remove
// from the database.
type RemovePhotoRequest struct {
	Token    string `json:"token"`
	Filename string `json:"filename"`
}

// AddGalleryItemRequest holds the token and gallery item to add to the
// database.
type AddGalleryItemRequest struct {
	Token string               `json:"token"`
	Item  entities.GalleryItem `json:"item"`
}

// RemoveGalleryItemRequest holds the token and gallery item id to remove
// from the database.
type RemoveGalleryItemRequest struct {
	Token string `json:"token"`
	ID    string `json:"id"`
}

// AddSlideRequest holds the token and slide information to add to the database.
type AddSlideRequest struct {
	Token string         `json:"token"`
	Slide entities.Slide `json:"slide"`
}

// RemoveSlideRequest holds the token and IDs to use to remove a slide from the
// database.
type RemoveSlideRequest struct {
	Token     string `json:"token"`
	GalleryID string `json:"gallery_id"`
	Name      string `json:"name"`
}

// AddUserRequest is the username and password to create a new user with,
// and may include an auth token.
type AddUserRequest struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is the username and password expected from the login endpoint.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
