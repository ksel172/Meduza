package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const delimiter = "jRlOPs"

type XORUtil struct {
	key []byte
}

func NewXorUtil() (*XORUtil, error) {
	xorUtils := &XORUtil{}

	key, err := xorUtils.GenerateRandomKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	xorUtils.key = key

	return xorUtils, nil
}

func (x *XORUtil) Encrypt(input string) (string, error) {
	key, err := x.GenerateRandomKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate key: %v", err)
	}

	data := []byte(input)
	encrypted := x.xorCrypt(data, key)
	magic := []byte(delimiter)

	combinedKey := make([]byte, len(magic)+len(key)+len(magic))
	copy(combinedKey[0:], magic)
	copy(combinedKey[len(magic):], key)
	copy(combinedKey[len(magic)+len(key):], magic)

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(encrypted))))
	if err != nil {
		return "", fmt.Errorf("failed to generate random position: %v", err)
	}
	startPosition := int(n.Int64()) + 1

	result := make([]byte, len(encrypted)+len(combinedKey))
	copy(result[0:], encrypted[:startPosition])
	copy(result[startPosition:], combinedKey)
	copy(result[startPosition+len(combinedKey):], encrypted[startPosition:])

	return string(result), nil
}

func (x *XORUtil) Decrypt(input string) (string, error) {
	data := []byte(input)
	magic := []byte(delimiter)

	magicPositions := x.search(data, magic)
	if len(magicPositions) < 2 {
		return "", fmt.Errorf("invalid encrypted data: magic delimiters not found")
	}

	keyLength := magicPositions[1] - magicPositions[0] - len(magic)
	key := make([]byte, keyLength)
	copy(key, data[magicPositions[0]+len(magic):magicPositions[0]+len(magic)+keyLength])

	totalDiscardLength := len(magic) + keyLength + len(magic)
	decrypted := make([]byte, len(data)-totalDiscardLength)
	copy(decrypted, data[:magicPositions[0]])
	copy(decrypted[magicPositions[0]:], data[magicPositions[0]+totalDiscardLength:])

	result := x.xorCrypt(decrypted, key)

	return string(result), nil
}

func (x *XORUtil) GenerateRandomKey() ([]byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(9))
	if err != nil {
		return nil, err
	}
	length := int(n.Int64()) + 24

	key := make([]byte, length)
	for i := 0; i < length; i++ {
		for {
			n, err := rand.Int(rand.Reader, big.NewInt(94))
			if err != nil {
				return nil, err
			}
			char := byte(n.Int64()) + 33

			if char != '\\' && char != '/' && char != ' ' {
				key[i] = char
				break
			}
		}
	}
	return key, nil
}

func (x *XORUtil) xorCrypt(data, key []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)

	for i := 0; i < len(result); i++ {
		result[i] ^= key[i%len(key)]
	}
	return result
}

func (x *XORUtil) search(src, pattern []byte) []int {
	var results []int
	maxFirstCharSlot := len(src) - len(pattern) + 1

	for i := 0; i < maxFirstCharSlot; i++ {
		if src[i] != pattern[0] {
			continue
		}

		match := true
		for j := len(pattern) - 1; j >= 1; j-- {
			if src[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			results = append(results, i)
		}
	}
	return results
}
