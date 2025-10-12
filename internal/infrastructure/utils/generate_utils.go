// file: internal/infrastructure/utils/generate_utils.go
package utils

import (
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
  rand.NewSource(NowUTC().UnixNano()))

func GenerateString(length int) string {
  return generateStringWithCharset(length, charset)
}

func generateStringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}
