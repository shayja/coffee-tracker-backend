// file: internal/entities/user_settings.go
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
type Setting int

const (
	SettingUnknown Setting = iota
	SettingBiometricEnabled
	SettingDarkMode
	SettingNotificationsEnabled
)

func (s Setting) IsValid() bool {
	switch s {
	case SettingBiometricEnabled,
		SettingDarkMode,
		SettingNotificationsEnabled:
		return true
	}
	return false
}

func (s Setting) ColumnName() string {
	switch s {
	case SettingBiometricEnabled:
		return "biometric_enabled"
	case SettingDarkMode:
		return "dark_mode"
	case SettingNotificationsEnabled:
		return "notifications_enabled"
	default:
		return ""
	}
}