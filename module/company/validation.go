package company

import (
	"strings"

	validator "github.com/subhasundardass/retui/internal/validation"
)

func ValidateCreate(input FormState) error {
	v := validator.New()

	v.Field("code", strings.TrimSpace(input.Code)).Required().MinLength(2).MaxLength(20).AlphaNumeric().UpperCase()
	v.Field("name", strings.TrimSpace(input.Name)).Required().MinLength(2).MaxLength(100)
	v.Field("legal_name", strings.TrimSpace(input.LegalName)).Required().MinLength(2).MaxLength(100)
	v.IntField("country", input.Country).Required()
	v.IntField("state", input.State).Required()

	return v.Error()
}

func ValidateUpdate(input FormState) error {
	v := validator.New()

	v.Field("code", strings.TrimSpace(input.Code)).Required().MinLength(2).MaxLength(20).AlphaNumeric().UpperCase()
	v.Field("name", strings.TrimSpace(input.Name)).Required().MinLength(2).MaxLength(100)
	v.Field("legal_name", strings.TrimSpace(input.LegalName)).Required().MinLength(2).MaxLength(100)
	v.IntField("country", input.Country).Required()
	v.IntField("state", input.State).Required()

	return v.Error()
}
