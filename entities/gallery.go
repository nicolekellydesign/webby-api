package entities

import "github.com/nicolekellydesign/webby-api/internal/db"

// GalleryItem represents an item in the main project gallery.
type GalleryItem struct {
	Name        string        `json:"name" db:"id"`
	Title       string        `json:"title" db:"title"`
	Caption     string        `json:"caption" db:"caption"`
	ProjectInfo string        `json:"projectInfo" db:"project_info"`
	Thumbnail   string        `json:"thumbnail" db:"thumbnail"`
	VideoKey    db.NullString `json:"videoKey,omitempty" db:"video_key"`
	Images      []string      `json:"images"`
}
