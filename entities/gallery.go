package entities

// GalleryItem represents an item in the main project gallery.
type GalleryItem struct {
	Name              string   `json:"id" db:"id"`
	TitleLine1        string   `json:"title_line_1" db:"title_line_1"`
	TitleLine2        string   `json:"title_line_2" db:"title_line_2"`
	ThumbnailLocation string   `json:"thumbnail_location" db:"thumbnail_location"`
	ThumbnailCaption  string   `json:"thumbnail_caption" db:"thumbnail_caption"`
	Slides            []*Slide `json:"slides,omitempty"`
}

// Slide represents a slide in a gallery project page.
type Slide struct {
	GalleryID string `json:"gallery_id" db:"gallery_id"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	Location  string `json:"location"`
}
