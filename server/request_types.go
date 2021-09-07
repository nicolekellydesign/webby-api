package server

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
