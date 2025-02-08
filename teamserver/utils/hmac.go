package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateHMAC generates an HMAC for a given message using the provided key
func generateHMAC(message, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func verifyHMAC(message []byte, receivedHMAC string, key []byte) bool {
	computedHMAC := generateHMAC(message, key)
	return computedHMAC == receivedHMAC
}
