package auth

//go:generate mockgen -destination=./auth_mock.go -package=auth -source=auth.go

// Authenticator handles user authentication.
type Authenticator interface {
	// Authenticate retuns true if username and password are correct.
	Authenticate(username, password string) (bool, error)
}

type staticAuth struct {
	username string
	password string
}

// NewStataticAuthenticator creates authenticator based on static credentials.
func NewStataticAuthenticator(username, password string) Authenticator {
	return staticAuth{
		username: username,
		password: password,
	}
}

func (sa staticAuth) Authenticate(username, password string) (bool, error) {
	return sa.username == username && sa.password == password, nil
}
