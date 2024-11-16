package utils

import (
	"encoding/base64"
	"errors"
)

func Base64URLSafeEncode(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func Base64URLSafeDecode(encoded string) (string, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("failed to decode base64 string: " + err.Error())
	}
	return string(decodedBytes), nil
}
