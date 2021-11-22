package v1

// AddUserRequest is the username and password to create a new user with.
type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangeThumbnailRequest holds the file name of the thumbnail to update a
// project with.
type ChangeThumbnailRequest struct {
	Thumbnail string `json:"thumbnail"`
}

// LoginRequest is the username and password expected from the login endpoint.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Extended bool   `json:"extended,omitempty"`
}
