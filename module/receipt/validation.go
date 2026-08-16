package receipt

import (
	"strings"

	validator "github.com/subhasundardass/retui/internal/validation"
)

func ValidateInput(input FormState) error {
	v := validator.New()

	v.Field("date", strings.TrimSpace(input.Date)).Required()
	// v.IntField("data", input.Amount).Required()
	v.IntField("data", input.RcptAccount).Required()

	return v.Error()
}
