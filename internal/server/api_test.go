package server

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_validate_category(t *testing.T) {
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

func Test_validate_store(t *testing.T) {
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

func Test_validate_position(t *testing.T) {
	testCases := []struct {
		desc string
		req  position
		errs string
	}{
		{
			desc: "valid_price_positive",
			req:  position{Price: decimal.NewFromInt(100)},
			errs: "",
		},
		{
			desc: "invalid_price_0",
			req:  position{Price: decimal.Zero},
			errs: "price must be greater than 0",
		},
		{
			desc: "invalid_price_negative",
			req:  position{Price: decimal.NewFromInt(-10)},
			errs: "price must be greater than 0",
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

func Test_validate_product(t *testing.T) {
	testCases := []struct {
		desc string
		req  product
		errs string
	}{
		{
			desc: "valid_name_3",
			req:  product{CategoryID: 1, Name: "XYZ", Description: "-"},
			errs: "",
		},
		{
			desc: "valid_name_160",
			req:  product{CategoryID: 1, Name: strings.Repeat("x", 160), Description: strings.Repeat("x", 2000)},
			errs: "",
		},
		{
			desc: "invalid_name_required",
			req:  product{CategoryID: 1, Description: "-"},
			errs: "name is a required field",
		},
		{
			desc: "invalid_name_2",
			req:  product{CategoryID: 1, Name: "XY", Description: "-"},
			errs: "name must be at least 3 characters in length",
		},
		{
			desc: "invalid_name_161",
			req:  product{CategoryID: 1, Name: strings.Repeat("x", 161), Description: "-"},
			errs: "name must be at maximum 160 characters in length",
		},
		{
			desc: "invalid_description_required",
			req:  product{CategoryID: 1, Name: "XYZ"},
			errs: "description is a required field",
		},
		{
			desc: "invalid_categoryId_required",
			req:  product{Name: "XYZ", Description: "-"},
			errs: "categoryId is a required field",
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
