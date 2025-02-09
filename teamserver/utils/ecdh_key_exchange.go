package utils

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// Using a P256 curve to generate the public and private keys
func GenerateECDHKeyPair() ([]byte, []byte, error) {
	curve := ecdh.P256()
	privKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	fmt.Printf("Public key: %+v", privKey.PublicKey().Bytes())
	return privKey.Bytes(), privKey.PublicKey().Bytes(), nil
}

func DeriveECDHSharedSecret(privKeyBytes []byte, peerPublicKeyBytes []byte) ([]byte, error) {
	privKey, err := ecdh.P256().NewPrivateKey(privKeyBytes)
	if err != nil {
		return nil, err
	}
	peerPublicKey, err := ecdh.P256().NewPublicKey(peerPublicKeyBytes)
	if err != nil {
		return nil, err
	}
	sharedKey, err := privKey.ECDH(peerPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive shared secret: %v", err)
	}

	// Hash the shared secret with SHA-256 to produce a 32-byte AES key
	// Agent needs to do the same
	//https://crypto.stackexchange.com/questions/57783/aes-encryption-using-a-diffie-hellman-question
	aesKey := sha256.Sum256(sharedKey)
	return aesKey[:], nil
}
