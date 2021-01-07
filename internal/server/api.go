package server

// errorResponse represents error response
type errorResponse struct {
	Error string `json:"error"`
}

// category represents category object.
type category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type store struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type item struct {
	ID          int64  `json:"id"`
	StoreID     int64  `json:"storeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}
