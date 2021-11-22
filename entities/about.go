package entities

// About holds information for the about page.
type About struct {
	Portrait  string `json:"portrait,omitempty"`
	Statement string `json:"statement,omitempty"`
	Resume    string `json:"resume,omitempty"`
}
