package company

import "github.com/subhasundardass/retui/internal/validation"

func ValidateCreate(input FormState) error {
	v := validation.New()
	v.Field("code", input.Code).Required().MinLength(2).MaxLength(20).AlphaNumeric().UpperCase()
	v.Field("name", input.Name).Required().MinLength(2).MaxLength(100)

	return v.Error()
}
