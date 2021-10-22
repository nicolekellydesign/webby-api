package entities

// GalleryItem represents an item in the main project gallery.
type GalleryItem struct {
	Name        string   `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	Caption     string   `json:"caption" db:"caption"`
	ProjectInfo string   `json:"projectInfo" db:"project_info"`
	Thumbnail   string   `json:"thumbnail" db:"thumbnail"`
	Images      []string `json:"images"`
}
