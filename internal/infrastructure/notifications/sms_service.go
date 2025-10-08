// file: internal/infrastructure/notifications/sms_service.go
package notifications

import "github.com/google/uuid"

// SMSService defines the contract for sending SMS messages
type SMSService interface {
    SendOTP(userID uuid.UUID, mobile string, otp string) error
}
