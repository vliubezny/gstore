package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_category_Validate(t *testing.T) {
	testCases := []struct {
		desc string
		req  category
		errs string
	}{
		{
			desc: "valid_name_2",
			req:  category{Name: "IT"},
			errs: "",
		},
		{
			desc: "valid_name_80",
			req:  category{Name: strings.Repeat("x", 80)},
			errs: "",
		},
		{
			desc: "invalid_name_empty",
			req:  category{Name: ""},
			errs: "name is a required field",
		},
		{
			desc: "invalid_name_1",
			req:  category{Name: "x"},
			errs: "name must be at least 2 characters in length",
		},
		{
			desc: "invalid_name_81",
			req:  category{Name: strings.Repeat("x", 81)},
			errs: "name must be at maximum 80 characters in length",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := validate(&tC.req)

			if tC.errs == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tC.errs)
			}
		})
	}
}

func Test_store_Validate(t *testing.T) {
	testCases := []struct {
		desc string
		req  store
		errs string
	}{
		{
			desc: "valid_name_2",
			req:  store{Name: "IT"},
			errs: "",
		},
		{
			desc: "valid_name_80",
			req:  store{Name: strings.Repeat("x", 80)},
			errs: "",
		},
		{
			desc: "invalid_name_empty",
			req:  store{Name: ""},
			errs: "name is a required field",
		},
		{
			desc: "invalid_name_1",
			req:  store{Name: "x"},
			errs: "name must be at least 2 characters in length",
		},
		{
			desc: "invalid_name_81",
			req:  store{Name: strings.Repeat("x", 81)},
			errs: "name must be at maximum 80 characters in length",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := validate(&tC.req)

			if tC.errs == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tC.errs)
			}
		})
	}
}
