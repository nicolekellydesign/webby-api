package entities

import "time"

// User represents a user that can manage Webby.
type User struct {
	ID        uint      `json:"id,omitempty" db:"id"`
	Username  string    `json:"username" db:"user_name"`
	Password  string    `json:"password,omitempty" db:"pwdhash"`
	Protected bool      `json:"protected" db:"protected"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
