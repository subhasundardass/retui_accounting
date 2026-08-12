package ledger

import (
	"strings"

	validator "github.com/subhasundardass/retui/internal/validation"
)

func ValidateForm(input LedgerState) error {
	v := validator.New()

	v.Field("code", strings.TrimSpace(input.Code)).Required().MinLength(2).MaxLength(20).UpperCase()
	v.Field("name", strings.TrimSpace(input.Name)).Required().MinLength(2).MaxLength(100)
	v.IntField("groupId", input.GroupID).Required()
	v.IntField("country", input.CountryID).Required()
	v.IntField("state", input.StateID).Required()

	return v.Error()
}
