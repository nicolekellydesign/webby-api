package v1

// CheckSessionResponse is sent when a client is trying to check
// if they have a valid session.
type CheckSessionResponse struct {
	Valid bool `json:"valid"`
}
