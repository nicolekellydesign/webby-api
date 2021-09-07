package entities

// Photo represents a photography photo.
type Photo struct {
	Filename string `json:"filename" db:"file_name"`
}
