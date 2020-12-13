package server

// Error represents error response
type Error struct {
	Error string `json:"error"`
}

// Category represents category object.
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetCategoriesResponse represents category list response.
type GetCategoriesResponse struct {
	Categories []*Category `json:"categories"`
}
