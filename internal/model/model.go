package model

import "github.com/shopspring/decimal"

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

// Product represents product item.
type Product struct {
	ID          int64
	CategoryID  int64
	Name        string
	Description string
}

// Position represents store prosition.
type Position struct {
	ProductID int64
	StoreID   int64
	Price     decimal.Decimal
}
