package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

// Session holds a user's login session information.
type Session struct {
	Token    string
	Username string `db:"user_name"`
	ID       uint   `db:"user_id"`
	Expires  time.Time
}

// NewSession creates a new login session with a unique
// token generated as a UUID, and an expiry time of
// five minutes.
func NewSession(id uint, username string, extended bool) (*Session, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	var expires time.Time
	if extended {
		expires = time.Now().Add(30 * 24 * time.Hour).UTC()
	} else {
		expires = time.Now().Add(5 * time.Minute).UTC()
	}

	return &Session{
		Token:    token.String(),
		Username: username,
		ID:       id,
		Expires:  expires,
	}, nil
}
