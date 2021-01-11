package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_category_Validate(t *testing.T) {
	testCases := []struct {
		desc  string
		req   category
		valid bool
	}{
		{
			desc:  "valid_name_2",
			req:   category{Name: "IT"},
			valid: true,
		},
		{
			desc:  "valid_name_80",
			req:   category{Name: strings.Repeat("x", 80)},
			valid: true,
		},
		{
			desc:  "invalid_name_empty",
			req:   category{Name: ""},
			valid: false,
		},
		{
			desc:  "invalid_name_1",
			req:   category{Name: "x"},
			valid: false,
		},
		{
			desc:  "invalid_name_81",
			req:   category{Name: strings.Repeat("x", 81)},
			valid: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.valid {
				assert.NoError(t, tC.req.Validate())
			} else {
				assert.Error(t, tC.req.Validate())
			}
		})
	}
}
