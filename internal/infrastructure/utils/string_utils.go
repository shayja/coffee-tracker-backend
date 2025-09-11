package utils

import "strings"

// NullIfEmpty trims the input string and returns nil if it is empty. Otherwise, returns the trimmed string.
func NullIfEmpty(s string) interface{} {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

// SafeToLower trims the input and returns a lowercase version. If the input is empty or only whitespace, returns nil
func ToLower(s string) interface{}  {
	if NullIfEmpty(s) == nil {
		return nil
	}
	return strings.ToLower(strings.TrimSpace(s))
}