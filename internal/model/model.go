package model

// Category represents product cetegory.
type Category struct {
	ID   int64
	Name string
}

// Store represents product store.
type Store struct {
	ID   int64
	Name string
}

// Item represents product item.
type Item struct {
	ID          int64
	StoreID     int64
	Name        string
	Description string
	Price       int64
}
