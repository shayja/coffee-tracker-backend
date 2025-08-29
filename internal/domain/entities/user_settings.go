// file: internal/domain/entities/user_settings.go
package entities

import "time"

// UserSettings entity mapped to DB
type UserSettings struct {
    UserID               string    `db:"user_id"`
    BiometricEnabled     bool      `db:"biometric_enabled"`
    DarkMode             bool      `db:"dark_mode"`
    NotificationsEnabled bool      `db:"notifications_enabled"`
    CreatedAt            time.Time `db:"created_at"`
    UpdatedAt            time.Time `db:"updated_at"`
}

// Enum-like type for allowed settings
type Setting string

const (
    SettingBiometricEnabled     Setting = "biometric_enabled"
    SettingDarkMode             Setting = "dark_mode"
    SettingNotificationsEnabled Setting = "notifications_enabled"
)

// AllowedSettings set for validation
var AllowedSettings = map[Setting]struct{}{
    SettingBiometricEnabled:     {},
    SettingDarkMode:             {},
    SettingNotificationsEnabled: {},
}

// Helper: check if a setting is valid
func (s Setting) IsValid() bool {
    _, ok := AllowedSettings[s]
    return ok
}
