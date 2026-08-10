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
