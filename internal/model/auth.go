package model

// User represents authenticated person.
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	IsAdmin      bool
}

// TokenPair groups access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}
