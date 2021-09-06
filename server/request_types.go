package server

import "github.com/nicolekellydesign/webby-api/entities"

// AddUserRequest is the data expected when a client requests to add a new user.
type AddUserRequest struct {
	Token string        `json:"token,omitempty"`
	User  entities.User `json:"user"`
}
