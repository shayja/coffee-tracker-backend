package utils

import "strings"

// NullIfEmpty trims the input string and returns nil if it is empty or nil.
// Supports both string and *string.
func NullIfEmpty(v any) any {
	switch s := v.(type) {
	case string:
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return nil
		}
		return trimmed
	case *string:
		if s == nil {
			return nil
		}
		trimmed := strings.TrimSpace(*s)
		if trimmed == "" {
			return nil
		}
		return trimmed
	default:
		return nil
	}
}

// SafeToLower trims and lowercases the string. Returns nil if input is nil/empty.
// Supports both string and *string.
func SafeToLower(v any) any {
	switch s := v.(type) {
	case string:
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return nil
		}
		return strings.ToLower(trimmed)
	case *string:
		if s == nil {
			return nil
		}
		trimmed := strings.TrimSpace(*s)
		if trimmed == "" {
			return nil
		}
		return strings.ToLower(trimmed)
	default:
		return nil
	}
}
