package v1

import "github.com/nicolekellydesign/webby-api/entities"

// AddGalleryItemRequest holds the token and gallery item to add to the
// database.
type AddGalleryItemRequest struct {
	Item entities.GalleryItem `json:"item"`
}

// AddSlideRequest holds the token and slide information to add to the database.
type AddSlideRequest struct {
	Slide entities.Slide `json:"slide"`
}

// AddUserRequest is the username and password to create a new user with,
// and may include an auth token.
type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is the username and password expected from the login endpoint.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Extended bool   `json:"extended,omitempty"`
}
