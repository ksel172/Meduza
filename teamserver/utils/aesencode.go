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

// AES stands for Advanced Encryption Standard. It is encryption of electronic data established by the
// U.S. National Institute of Standards and Technology (NIST) in 2001.

// Aes is a structure which provides encryption and decryption functions.
// Key is a 32 bit secret code, which help in encoding and decoding of data.
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

	// NewCipher creates and returns a cipher block using the provided key.
	// A cipher block is a specific-sized portion of data.
	// Each block is processed independently, either encrypted or decrypted, using the same key.
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf(`error while creating a newcipher: %w`, err)
	}

	//NewGCM creates a new AES from the cipher block.
	//AES-GCM provides authenticated encryption, which ensures both confidentiality and integrity of data.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf(`error while creating a newGCM: %w`, err)
	}

	//Nonce ensures that each encryption operation produces unique cipertext.
	//Nonce size is determined by the GCM instance.
	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf(`error while reading a nonce: %w`, err)
	}

	// The Seal method  encrypts and authenticate by append the data to the nonce.
	// It's helps to generates the ciphertext.
	cipherData := aesGCM.Seal(nonce, nonce, data, nil)

	//Base64 enconding ensures the ciphertext can be safely transported over text-based protocols Like HTTP.
	return []byte(base64.StdEncoding.EncodeToString(cipherData)), nil
}

// AesEncrypt decrypts ciphertext using AES-GCM.
// It's takes ciphertext as []byte.
func (a *Aes) AesDecrypt(ciphertext []byte) ([]byte, error) {

	// It decodes the ciphertext from base to raw byte
	decodeCipherText, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, fmt.Errorf(`error Decoding ciphertext: %w`, err)
	}

	// Creates a new AES cipher block using the provided key.
	// The key must be the same used while encryption.
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf(`error Creating a cipher block instance: %w`, err)
	}

	//NewGCM creates a AES-GCM from the cipher block.
	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf(`error Creating a NewGCM instance: %w`, err)
	}

	//AES-GCM provides authenticated encryption, which ensures both confidentiality and integrity of data.
	//It also ensures the decoded ciphertext is long enough to contain both the nonce and encrypted data.
	nonceSize := aesGcm.NonceSize()
	if len(decodeCipherText) < nonceSize {
		return nil, errors.New("ciphertext is too short")
	}

	// It's extract the nonce and encrypted data from the decoded ciphertext.
	// The nonce is the first portion of the ciphertext, while the encrypted data follows it.
	nonce, encryptedData := decodeCipherText[:nonceSize], decodeCipherText[nonceSize:]

	// The `Open` function reverses the encryption process using the nonce and key.
	data, err := aesGcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf(`error decoding encrypteddata: %w`, err)
	}

	// Return the decrypted data.
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
