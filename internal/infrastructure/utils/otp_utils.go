// file: internal/infrastructure/utils/otp_utils.go
package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"coffee-tracker-backend/internal/infrastructure/config"
)

// GenerateOTP generates a 6-digit OTP based on strength
func GenerateOTP(strength config.OtpStrength) (string, error) {
	switch strength {
	case config.OTP_EASY:
		// easy: 6-digit numeric, deterministic enough for "easy"
		n, err := randIntEasy(100000, 999999)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%06d", n), nil

	case config.OTP_STRONG:
		// strong: 6-digit numeric using crypto/rand
		n, err := randIntCrypto(100000, 999999)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%06d", n), nil

	default:
		return "", errors.New("unknown OTP strength: must be 'easy' or 'strong'")
	}
}

// randIntEasy: simple "easy" OTP generator
func randIntEasy(min, max int) (int, error) {
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		return 0, err
	}
	return min + int(b[0])%(max-min+1), nil
}

// randIntCrypto: secure crypto-based random integer
func randIntCrypto(min, max int) (int, error) {
	diff := int64(max - min + 1)
	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, err
	}
	return min + int(n.Int64()), nil
}
