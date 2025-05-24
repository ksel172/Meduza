package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

const (
	AES256KeySize = 32
)

// GenerateAES256Key returns a 32-byte AES key.
func GenerateAES256Key() ([]byte, error) {
	key := make([]byte, AES256KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate AES-256 key: %v", err)
	}
	return key, nil
}

// AesEncrypt encrypts data using AES-GCM with the provided key.
func AesEncrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error while creating a new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error while creating a new GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("error while reading a nonce: %w", err)
	}

	cipherData := aesGCM.Seal(nil, nonce, data, nil)

	fullData := append(nonce, cipherData...) // nonce + encrypted data

	return fullData, nil
}

// AesDecrypt decrypts data using AES-GCM with the provided key.
func AesDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating a cipher block instance: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating a NewGCM instance: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	data, err := aesGCM.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("error decoding encrypted data: %w", err)
	}

	return data, nil
}
