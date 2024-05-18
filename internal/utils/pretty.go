package utils

import (
	"unicode/utf8"
)

// TruncateString truncates a string to the specified number of characters
func TruncateString(str string, limit int) string {
	if limit <= 0 {
		return "" // return empty string if limit is not positive
	}

	if utf8.RuneCountInString(str) <= limit {
		return str // return the original string if it's shorter than the limit
	}

	// Truncate the string to the specified number of characters
	runes := []rune(str)
	return string(runes[:limit])
}
