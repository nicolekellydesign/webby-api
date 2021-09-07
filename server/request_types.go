package server

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
