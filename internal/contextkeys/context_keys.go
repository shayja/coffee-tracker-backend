// file: internal/contextkeys/context_keys.go
// Package contextkeys defines keys for storing values in request context.
package contextkeys

type ContextKey string

const (
	UserIDKey      ContextKey = "userID"
	CurrentUserKey ContextKey = "currentUser"
)