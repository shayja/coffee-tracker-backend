// // file: internal/infrastructure/notifications/twilio_sms.go
// file: internal/infrastructure/notifications/noop_sms_service.go
package notifications

import (
	"log"

	"github.com/google/uuid"
)

type NoOpSMSService struct{}

func NewNoOpSMSService() *NoOpSMSService {
	return &NoOpSMSService{}
}

// SendOTP just logs the OTP instead of sending it
func (s *NoOpSMSService) SendOTP(userID uuid.UUID, mobile string, otp string) error {
	log.Printf("[NoOpSMS] OTP for user %s, mobile %s: %s", userID, mobile, otp)
	return nil
}
