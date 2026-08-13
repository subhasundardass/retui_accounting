package journal

import (
	"fmt"
	"strings"

	validator "github.com/subhasundardass/retui/internal/validation"
)

const epsilon = 0.005 // tolerance for float rounding on money comparisons

// Balanced reports whether debit and credit are equal within floating
// point tolerance. Never compare money totals with == directly.
func Balanced(debit, credit float64) bool {
	diff := debit - credit
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

func ValidateFormShape(input FormState) error {
	v := validator.New()

	v.Field("vcNo", strings.TrimSpace(input.VcNo)).Required().MinLength(1).MaxLength(20)
	v.Field("vcDate", strings.TrimSpace(input.VcDate)).Required().MinLength(2).MaxLength(100)

	return v.Error()
}

// Totals returns the sum of debits and credits across all lines.
func Totals(in VoucherInput) (debit, credit float64) {
	for _, l := range in.Lines {
		debit += l.Debit
		credit += l.Credit
	}
	return debit, credit
}

// Validate enforces the rules every journal voucher must satisfy,
// regardless of which module (receipt, payment, sale, ...) created it.
func Validate(in VoucherInput) error {
	if in.Type == "" {
		return fmt.Errorf("voucher type is required")
	}
	if in.VoucherNo == "" {
		return fmt.Errorf("voucher number is required")
	}
	if in.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	if len(in.Lines) < 2 {
		return fmt.Errorf("a voucher needs at least two lines (one debit, one credit)")
	}

	for i, l := range in.Lines {
		if l.LedgerID == 0 {
			return fmt.Errorf("line %d: ledger is required", i+1)
		}
		if l.Debit < 0 || l.Credit < 0 {
			return fmt.Errorf("line %d: amounts cannot be negative", i+1)
		}
		if l.Debit > 0 && l.Credit > 0 {
			return fmt.Errorf("line %d: a line cannot be both debit and credit", i+1)
		}
		if l.Debit == 0 && l.Credit == 0 {
			return fmt.Errorf("line %d: amount is required", i+1)
		}
	}

	debit, credit := Totals(in)
	if debit == 0 {
		return fmt.Errorf("total debit cannot be zero")
	}
	if diff := debit - credit; diff > epsilon || diff < -epsilon {
		return fmt.Errorf("voucher does not balance: debit %.2f, credit %.2f", debit, credit)
	}

	return nil
}
