package entities

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// NullString wraps sql.NullString and implements some interfaces to make life easier.
type NullString sql.NullString

// MarshalJSON implements the JSON marshal interface for NullString.
func (s *NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		s.String = ""
		s.Valid = true
	}

	return json.Marshal(s.String)
}

// UnmarshalJSON implements the JSON unmarshal interface for NullString.
func (s *NullString) UnmarshalJSON(data []byte) error {
	var str *string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != nil {
		s.String = *str
		s.Valid = true
	} else {
		s.Valid = false
	}

	return nil
}

// Scan implements the Scanner interface for NullString.
func (s *NullString) Scan(value interface{}) error {
	var str sql.NullString
	if err := str.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*s = NullString{str.String, false}
	} else {
		*s = NullString{s.String, true}
	}

	return nil
}

// GalleryItem represents an item in the main project gallery.
type GalleryItem struct {
	Name        string     `json:"name" db:"id"`
	Title       string     `json:"title" db:"title"`
	Caption     string     `json:"caption" db:"caption"`
	ProjectInfo string     `json:"projectInfo" db:"project_info"`
	Thumbnail   string     `json:"thumbnail" db:"thumbnail"`
	EmbedURL    NullString `json:"embedURL,omitempty" db:"embed_url"`
	Images      []string   `json:"images"`
}
