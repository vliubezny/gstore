package model

// User represents authenticated person.
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	IsAdmin      bool
}
