package util

import "strconv"

func StringToInt(s string, defaultValue int) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return num
}

func IntToString(num int) string {
	return strconv.Itoa(num)
}

func StringToFloat(s string, defaultValue float64) float64 {
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return num
}

func FloatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64)
}

// deref safely reads an Optional().Nillable() string field, returning ""
// instead of panicking when the value was never set.
func Deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
