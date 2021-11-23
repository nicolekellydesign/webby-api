package entities

// About holds information for the about page.
type About struct {
	Portrait  string `json:"portrait,omitempty"`
	Statement string `json:"statement,omitempty"`
	Resume    string `json:"resume,omitempty"`
}

// MergeAbout combines two About structs, returning a new struct with the merged
// values.
func MergeAbout(left, right About) *About {
	ret := left

	if right.Portrait != "" {
		ret.Portrait = right.Portrait
	}

	if right.Statement != "" {
		ret.Statement = right.Statement
	}

	if right.Resume != "" {
		ret.Resume = right.Resume
	}

	return &ret
}
