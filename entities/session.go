package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

// Session holds a user's login session information.
type Session struct {
	Token    string
	Username string `db:"user_name"`
	Expires  time.Time
}

// NewSession creates a new login session with a unique
// token generated as a UUID, and an expiry time of
// two minutes.
func NewSession(username string) (*Session, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &Session{
		Token:    token.String(),
		Username: username,
		Expires:  time.Now().Add(120 * time.Second),
	}, nil
}
