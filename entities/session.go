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
	Created  time.Time
	MaxAge   int `db:"max_age"`
}

// NewSession creates a new login session with a unique
// token generated as a UUID, and an expiry time of
// five minutes.
func NewSession(id uint, username string, extended bool) (*Session, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	created := time.Now().UTC()
	maxAge := 0
	if extended {
		maxAge = created.Add(30 * 24 * time.Hour).Second()
	}

	return &Session{
		Token:    token.String(),
		Username: username,
		ID:       id,
		Created:  created,
		MaxAge:   maxAge,
	}, nil
}
