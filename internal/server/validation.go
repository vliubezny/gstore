package server

import (
	"bytes"
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/shopspring/decimal"
)

var (
	validation *validator.Validate
	trans      ut.Translator
)

func init() {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ = uni.GetTranslator("en")

	validation = validator.New()
	en_translations.RegisterDefaultTranslations(validation, trans)

	validation.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validation.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if d, ok := field.Interface().(decimal.Decimal); ok {
			f, _ := d.Float64()
			return f
		}
		return nil
	}, decimal.Decimal{})
}

// validate validates struct base on field tags and returns tranlated error.
func validate(target interface{}) error {
	err := validation.Struct(target)
	if errs, ok := err.(validator.ValidationErrors); ok {

		b := bytes.NewBufferString("")
		for _, e := range errs {
			b.WriteString(e.Translate(trans))
			b.WriteString("\n")
		}

		return errors.New(strings.TrimSpace(b.String()))
	}
	return err
}
