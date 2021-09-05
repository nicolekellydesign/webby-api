package entities

// User represents a user that can manage Webby.
type User struct {
	Username string `json:"username" db:"user_name"`
	Password string `json:"password" db:"pwdhash"`
}
