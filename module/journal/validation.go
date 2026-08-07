package journal

import (
	"strings"

	validator "github.com/subhasundardass/retui/internal/validation"
)

func ValidateCreate(input FormState) error {
	v := validator.New()

	v.Field("vcNo", strings.TrimSpace(input.VcNo)).Required().MinLength(1).MaxLength(20)
	v.Field("vcDate", strings.TrimSpace(input.VcDate)).Required().MinLength(2).MaxLength(100)

	return v.Error()
}
