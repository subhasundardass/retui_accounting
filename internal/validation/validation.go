package validation

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Validator struct {
	errors   []string
	errorMap map[string]string
}

type Field struct {
	v       *Validator
	name    string
	value   string
	skipped bool
}

func New() *Validator {
	return &Validator{
		errorMap: make(map[string]string),
	}
}

func (v *Validator) Field(name, value string) *Field {
	return &Field{
		v:     v,
		name:  name,
		value: value,
	}
}

func (v *Validator) Errors() []string {
	return v.errors
}

func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

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
	full := fmt.Sprintf("%s: %s", f.name, msg)

	f.v.errors = append(f.v.errors, full)

	// Keep only the first error for each field.
	if _, exists := f.v.errorMap[f.name]; !exists {
		f.v.errorMap[f.name] = msg
	}
}

// --------------------------------------------------
// CONDITIONS
// --------------------------------------------------

// When skips all subsequent rules when condition is false.
func (f *Field) When(condition bool) *Field {
	if !condition {
		f.skipped = true
	}

	return f
}

// RequiredIf is a convenient conditional Required().
func (f *Field) RequiredIf(condition bool) *Field {
	if condition {
		f.Required()
	}

	return f
}

// Nullable allows an empty value to skip subsequent validation.
func (f *Field) Nullable() *Field {
	if strings.TrimSpace(f.value) == "" {
		f.skipped = true
	}

	return f
}

// --------------------------------------------------
// BASIC
// --------------------------------------------------

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

// --------------------------------------------------
// STRING
// --------------------------------------------------

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

	if strings.ContainsAny(f.value, " \t\r\n") {
		f.add("cannot contain spaces")
	}

	return f
}

func (f *Field) Contains(value string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if !strings.Contains(f.value, value) {
		f.add(fmt.Sprintf("must contain %q", value))
	}

	return f
}

func (f *Field) StartsWith(value string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if !strings.HasPrefix(f.value, value) {
		f.add(fmt.Sprintf("must start with %q", value))
	}

	return f
}

func (f *Field) EndsWith(value string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if !strings.HasSuffix(f.value, value) {
		f.add(fmt.Sprintf("must end with %q", value))
	}

	return f
}

func (f *Field) Equal(value string) *Field {
	if f.skipped {
		return f
	}

	if f.value != value {
		f.add("does not match")
	}

	return f
}

func (f *Field) NotEqual(value string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if f.value == value {
		f.add("contains an invalid value")
	}

	return f
}

func (f *Field) In(values ...string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	for _, value := range values {
		if f.value == value {
			return f
		}
	}

	f.add("contains an invalid value")

	return f
}

func (f *Field) NotIn(values ...string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	for _, value := range values {
		if f.value == value {
			f.add("contains an invalid value")
			return f
		}
	}

	return f
}

// --------------------------------------------------
// EMAIL / PASSWORD
// --------------------------------------------------

var emailRegex = regexp.MustCompile(
	`^[^\s@]+@[^\s@]+\.[^\s@]+$`,
)

func (f *Field) Email() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if !emailRegex.MatchString(f.value) {
		f.add("invalid email address")
	}

	return f
}

func (f *Field) Password() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if len([]rune(f.value)) < 8 {
		f.add("must be at least 8 characters")
		return f
	}

	return f
}

func (f *Field) Confirmed(password string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if f.value != password {
		f.add("does not match")
	}

	return f
}

// --------------------------------------------------
// REGEX
// --------------------------------------------------

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

// --------------------------------------------------
// NUMBERS
// --------------------------------------------------

func (f *Field) Number() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if _, err := strconv.ParseFloat(f.value, 64); err != nil {
		f.add("must be a number")
	}

	return f
}

func (f *Field) Integer() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if _, err := strconv.Atoi(f.value); err != nil {
		f.add("must be an integer")
	}

	return f
}

func (f *Field) Decimal() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	value, err := strconv.ParseFloat(f.value, 64)

	if err != nil {
		f.add("must be a decimal number")
		return f
	}

	if !strings.Contains(f.value, ".") {
		_ = value
		f.add("must be a decimal number")
	}

	return f
}

func (f *Field) Min(n float64) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	value, err := strconv.ParseFloat(f.value, 64)

	if err != nil {
		f.add("must be a number")
		return f
	}

	if value < n {
		f.add(fmt.Sprintf("must be at least %g", n))
	}

	return f
}

func (f *Field) Max(n float64) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	value, err := strconv.ParseFloat(f.value, 64)

	if err != nil {
		f.add("must be a number")
		return f
	}

	if value > n {
		f.add(fmt.Sprintf("must be at most %g", n))
	}

	return f
}

func (f *Field) Between(min, max float64) *Field {
	return f.Min(min).Max(max)
}

func (f *Field) Positive() *Field {
	return f.Min(0)
}

func (f *Field) Negative() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	value, err := strconv.ParseFloat(f.value, 64)

	if err != nil {
		f.add("must be a number")
		return f
	}

	if value >= 0 {
		f.add("must be negative")
	}

	return f
}

// --------------------------------------------------
// DATE / TIME
// --------------------------------------------------

// Date validates a date using Go's time layout.
//
// Example:
//
//	Date("2006-01-02")
//	Date("02-01-2006")
func (f *Field) Date(layout string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if _, err := time.Parse(layout, f.value); err != nil {
		f.add(fmt.Sprintf("must be a valid date (%s)", layout))
	}

	return f
}

// DateTime validates date and time.
func (f *Field) DateTime(layout string) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if _, err := time.Parse(layout, f.value); err != nil {
		f.add(fmt.Sprintf("must be a valid date and time (%s)", layout))
	}

	return f
}

func (f *Field) Before(layout string, date time.Time) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	value, err := time.Parse(layout, f.value)

	if err != nil {
		f.add(fmt.Sprintf("must be a valid date (%s)", layout))
		return f
	}

	if !value.Before(date) {
		f.add("must be before the specified date")
	}

	return f
}

func (f *Field) After(layout string, date time.Time) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	value, err := time.Parse(layout, f.value)

	if err != nil {
		f.add(fmt.Sprintf("must be a valid date (%s)", layout))
		return f
	}

	if !value.After(date) {
		f.add("must be after the specified date")
	}

	return f
}

// --------------------------------------------------
// NETWORK
// --------------------------------------------------

func (f *Field) IP() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if net.ParseIP(f.value) == nil {
		f.add("must be a valid IP address")
	}

	return f
}

func (f *Field) IPv4() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	ip := net.ParseIP(f.value)

	if ip == nil || ip.To4() == nil {
		f.add("must be a valid IPv4 address")
	}

	return f
}

func (f *Field) IPv6() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	ip := net.ParseIP(f.value)

	if ip == nil || ip.To4() != nil {
		f.add("must be a valid IPv6 address")
	}

	return f
}

// --------------------------------------------------
// URL / UUID
// --------------------------------------------------

func (f *Field) URL() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	u, err := url.ParseRequestURI(f.value)

	if err != nil || u.Scheme == "" || u.Host == "" {
		f.add("invalid URL")
	}

	return f
}

var uuidRegex = regexp.MustCompile(
	`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`,
)

func (f *Field) UUID() *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if !uuidRegex.MatchString(f.value) {
		f.add("invalid UUID")
	}

	return f
}

// --------------------------------------------------
// CUSTOM
// --------------------------------------------------

func (f *Field) Custom(fn func(value string) error) *Field {
	if f.skipped || f.value == "" {
		return f
	}

	if err := fn(f.value); err != nil {
		f.add(err.Error())
	}

	return f
}

// ==================================================
// INT FIELD
// ==================================================

type IntField struct {
	v       *Validator
	name    string
	value   int
	skipped bool
}

func (v *Validator) IntField(name string, value int) *IntField {
	return &IntField{
		v:     v,
		name:  name,
		value: value,
	}
}

func (f *IntField) add(msg string) {
	full := fmt.Sprintf("%s: %s", f.name, msg)

	f.v.errors = append(f.v.errors, full)

	if _, exists := f.v.errorMap[f.name]; !exists {
		f.v.errorMap[f.name] = msg
	}
}

func (f *IntField) When(condition bool) *IntField {
	if !condition {
		f.skipped = true
	}

	return f
}

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

func (f *IntField) Between(min, max int) *IntField {
	return f.Min(min).Max(max)
}

func (f *IntField) Positive() *IntField {
	return f.Min(1)
}

func (f *IntField) Negative() *IntField {
	if f.skipped {
		return f
	}

	if f.value >= 0 {
		f.add("must be negative")
	}

	return f
}

func (f *IntField) OneOf(allowed ...int) *IntField {
	if f.skipped {
		return f
	}

	for _, value := range allowed {
		if f.value == value {
			return f
		}
	}

	f.add("is not a valid selection")

	return f
}

func (f *IntField) NotOneOf(values ...int) *IntField {
	if f.skipped {
		return f
	}

	for _, value := range values {
		if f.value == value {
			f.add("contains an invalid value")
			return f
		}
	}

	return f
}
