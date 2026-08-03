package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Validator struct {
	errors   []string
	errorMap map[string]string // first error per field
}

type Field struct {
	v       *Validator
	name    string
	value   string
	skipped bool // set by When(false) — skips all subsequent rules
}

func New() *Validator {
	return &Validator{errorMap: make(map[string]string)}
}

func (v *Validator) Field(name, value string) *Field {
	return &Field{v: v, name: name, value: value}
}

func (v *Validator) Errors() []string { return v.errors }
func (v *Validator) HasErrors() bool  { return len(v.errors) > 0 }

func (v *Validator) ErrorMap() map[string]string {
	return v.errorMap
}

func (v *Validator) Error() error {
	if len(v.errors) == 0 {
		return nil
	}
	return errors.New(strings.Join(v.errors, "\n"))
}

func (f *Field) add(msg string) {
	full := fmt.Sprintf("%s: %s ", f.name, msg)
	f.v.errors = append(f.v.errors, full)
	// Record only the first error per field
	if _, exists := f.v.errorMap[f.name]; !exists {
		f.v.errorMap[f.name] = msg
	}
}

// When skips all subsequent rules on this field when condition is false.
//
//	v.Field("phone", val).When(notifyBySMS).Required().Numeric()
func (f *Field) When(condition bool) *Field {
	if !condition {
		f.skipped = true
	}
	return f
}

func (f *Field) Required() *Field {
	if f.skipped {
		return f
	}
	if strings.TrimSpace(f.value) == "" {
		f.add("is required")
	}
	return f
}

func (f *Field) MinLength(n int) *Field {
	if f.skipped || f.value == "" {
		return f
	}
	if len([]rune(f.value)) < n {
		f.add(fmt.Sprintf("minimum length is %d characters", n))
	}
	return f
}

func (f *Field) MaxLength(n int) *Field {
	if f.skipped || f.value == "" {
		return f
	}
	if len([]rune(f.value)) > n {
		f.add(fmt.Sprintf("maximum length is %d characters", n))
	}
	return f
}

func (f *Field) Length(n int) *Field {
	if f.skipped || f.value == "" {
		return f
	}
	if len([]rune(f.value)) != n {
		f.add(fmt.Sprintf("must be exactly %d characters", n))
	}
	return f
}

func (f *Field) Min(n float64) *Field {
	if f.skipped || f.value == "" {
		return f
	}
	v, err := strconv.ParseFloat(f.value, 64)
	if err != nil {
		f.add("must be a number")
		return f
	}
	if v < n {
		f.add(fmt.Sprintf("must be at least %g", n))
	}
	return f
}

func (f *Field) Max(n float64) *Field {
	if f.skipped || f.value == "" {
		return f
	}
	v, err := strconv.ParseFloat(f.value, 64)
	if err != nil {
		f.add("must be a number")
		return f
	}
	if v > n {
		f.add(fmt.Sprintf("must be at most %g", n))
	}
	return f
}

func (f *Field) Between(min, max float64) *Field {
	return f.Min(min).Max(max)
}

func (f *Field) Alpha() *Field {
	if f.skipped || f.value == "" {
		return f
	}
	for _, r := range f.value {
		if !unicode.IsLetter(r) {
			f.add("must contain only letters")
			break
		}
	}
	return f
}

func (f *Field) Numeric() *Field {
	if f.skipped || f.value == "" {
		return f
	}
	for _, r := range f.value {
		if !unicode.IsDigit(r) {
			f.add("must contain only digits")
			break
		}
	}
	return f
}

func (f *Field) AlphaNumeric() *Field {
	if f.skipped || f.value == "" {
		return f
	}
	for _, r := range f.value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			f.add("must contain only letters and digits")
			break
		}
	}
	return f
}

func (f *Field) UpperCase() *Field {
	if f.skipped || f.value == "" {
		return f
	}
	if f.value != strings.ToUpper(f.value) {
		f.add("must be uppercase")
	}
	return f
}

func (f *Field) LowerCase() *Field {
	if f.skipped || f.value == "" {
		return f
	}
	if f.value != strings.ToLower(f.value) {
		f.add("must be lowercase")
	}
	return f
}

func (f *Field) NoSpace() *Field {
	if f.skipped || f.value == "" {
		return f
	}
	if strings.ContainsRune(f.value, ' ') {
		f.add("cannot contain spaces")
	}
	return f
}

func (f *Field) Email() *Field {
	if f.skipped || f.value == "" {
		return f
	}
	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !re.MatchString(f.value) {
		f.add("invalid email address")
	}
	return f
}

func (f *Field) Regex(pattern, message string) *Field {
	if f.skipped || f.value == "" {
		return f
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		f.add("invalid validation pattern: " + err.Error())
		return f
	}
	if !re.MatchString(f.value) {
		f.add(message)
	}
	return f
}

func (f *Field) Custom(fn func(value string) error) *Field {
	if f.skipped || f.value == "" {
		return f
	}
	if err := fn(f.value); err != nil {
		f.add(err.Error())
	}
	return f
}

// =====INT
type IntField struct {
	v       *Validator
	name    string
	value   int
	skipped bool // set by When(false) — skips all subsequent rules
}

// IntField starts a validation chain for an integer value, typically an
// ID/FK into a lookup table (e.g. country, state, category).
func (v *Validator) IntField(name string, value int) *IntField {
	return &IntField{v: v, name: name, value: value}
}

func (f *IntField) add(msg string) {
	f.v.errorMap[f.name] = f.name + " " + msg

}

// When conditionally skips all subsequent rules on this field, matching
// Field.When's semantics.
func (f *IntField) When(cond bool) *IntField {
	if !cond {
		f.skipped = true
	}
	return f
}

// Required fails if value is the zero value (0). Use this for FK/ID
// fields where 0 means "not selected" — not for fields where 0 is a
// legitimate value (use Min/Max or a custom rule instead in that case).
func (f *IntField) Required() *IntField {
	if f.skipped {
		return f
	}
	if f.value == 0 {
		f.add("is required")
	}
	return f
}

func (f *IntField) Min(n int) *IntField {
	if f.skipped {
		return f
	}
	if f.value < n {
		f.add(fmt.Sprintf("must be at least %d", n))
	}
	return f
}

func (f *IntField) Max(n int) *IntField {
	if f.skipped {
		return f
	}
	if f.value > n {
		f.add(fmt.Sprintf("must be at most %d", n))
	}
	return f
}

// OneOf fails if value is not in the given set — useful for validating
// an FK ID against a known list of valid IDs (e.g. loaded from DB).
func (f *IntField) OneOf(allowed ...int) *IntField {
	if f.skipped {
		return f
	}
	for _, a := range allowed {
		if f.value == a {
			return f
		}
	}
	f.add("is not a valid selection")
	return f
}
