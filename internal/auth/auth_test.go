package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStaticAuthenticator_Authenticate(t *testing.T) {
	testCases := []struct {
		desc     string
		username string
		password string
		res      bool
	}{
		{
			desc:     "success",
			username: "admin",
			password: "admin123",
			res:      true,
		},
		{
			desc:     "invalid username",
			username: "test",
			password: "admin123",
			res:      false,
		},
		{
			desc:     "invalid password",
			username: "admin",
			password: "test",
			res:      false,
		},
		{
			desc:     "invalid username and password",
			username: "test",
			password: "test",
			res:      false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			a := NewStataticAuthenticator("admin", "admin123")

			res, err := a.Authenticate(tC.username, tC.password)

			require.NoError(t, err)
			assert.Equal(t, tC.res, res)
		})
	}
}
