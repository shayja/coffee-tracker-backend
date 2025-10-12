// file: internal/infrastructure/utils/string_utils.go
package utils

import "time"

// NowUTC returns the current UTC time.
func NowUTC() time.Time {
    return time.Now().UTC()
}