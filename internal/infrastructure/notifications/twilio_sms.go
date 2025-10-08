// file: internal/infrastructure/notifications/twilio_sms_service.go
package notifications

import (
	"github.com/google/uuid"
	//"github.com/sfreiberg/gotwilio"
)

type TwilioSMSService struct {
	//twilio *gotwilio.Twilio
	from   string
}

func NewTwilioSMSService(accountSID, authToken, from string) *TwilioSMSService {
	//twilio := gotwilio.NewTwilioClient(accountSID, authToken)
	return &TwilioSMSService{
		//twilio: twilio,
		from:   from,
	}
}

// SendOTP sends an OTP SMS
func (s *TwilioSMSService) SendOTP(userID uuid.UUID, to string, otp string) error {
	// msg := fmt.Sprintf("Your OTP code is: %s", otp)
	// _, exc, err := s.twilio.SendSMS(s.from, to, msg, "", "")
	// if err != nil {
	// 	return err
	// }
	// if exc != nil {
	// 	return fmt.Errorf("twilio exception: %v", exc)
	// }
	return nil
}
