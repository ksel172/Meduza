package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const (
	AES256KeySize = 32
)

type Aes struct {
	key []byte
}

// Returns a new instance of Aes.
func NewAES() (*Aes, error) {
	generateKey, err := GenerateAES256Key()
	if err != nil {
		return nil, fmt.Errorf("invalid key length and failed to generate new AES-256 key: %v", err)
	}
	return &Aes{key: generateKey}, nil
}

// AesEncrypt encrypts data using AES-GCM.
// It's takes data as []byte.
func (a *Aes) AesEncrypt(data []byte) ([]byte, error) {

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf(`error while creating a newcipher: %w`, err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf(`error while creating a newGCM: %w`, err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf(`error while reading a nonce: %w`, err)
	}

	cipherData := aesGCM.Seal(nonce, nonce, data, nil)

	return []byte(base64.StdEncoding.EncodeToString(cipherData)), nil
}

func (a *Aes) AesDecrypt(ciphertext []byte) ([]byte, error) {

	decodeCipherText, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, fmt.Errorf(`error Decoding ciphertext: %w`, err)
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf(`error Creating a cipher block instance: %w`, err)
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf(`error Creating a NewGCM instance: %w`, err)
	}

	nonceSize := aesGcm.NonceSize()
	if len(decodeCipherText) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	nonce, encryptedData := decodeCipherText[:nonceSize], decodeCipherText[nonceSize:]

	data, err := aesGcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf(`error decoding encrypteddata: %w`, err)
	}

	return data, nil
}

// GenerateAES256Key return a 32 bit AESKey in byte format.
func GenerateAES256Key() ([]byte, error) {
	key := make([]byte, AES256KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate AES-256 key: %v", err)
	}
	return key, nil
}
