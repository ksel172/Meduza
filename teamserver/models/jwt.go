package models

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService handles token operations
type JWTService struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	revokedTokens   map[string]time.Time
	mu              sync.Mutex
}

func NewJWTService(secret string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTService {
	return &JWTService{
		secret:          secret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		revokedTokens:   make(map[string]time.Time),
	}
}

// UserClaim represents custom claims for JWT
type UserClaim struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokens generates an access and refresh token pair
func (j *JWTService) GenerateTokens(userID, role string) (*AuthResponse, error) {
	now := time.Now()
	accessClaims := &UserClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err := j.generateToken(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshClaims := &UserClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken, err := j.generateToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken validates and parses a token
func (j *JWTService) ValidateToken(tokenStr string) (*UserClaim, error) {

	if j.IsTokenRevoked(tokenStr) {
		return nil, errors.New("token has been revoked")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		log.Printf("Token validation error: %v", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaim)
	if !ok || !token.Valid {
		log.Println("Invalid claims or token")
		return nil, fmt.Errorf("invalid claims or token")
	}

	return claims, nil
}

// RefreshTokens generates a new access and refresh token if the refresh token is valid
func (j *JWTService) RefreshTokens(refreshToken string) (*AuthResponse, error) {
	// Step 1: Validate the refresh token
	if j.IsTokenRevoked(refreshToken) {
		return nil, fmt.Errorf("revoked token")
	}

	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Step 2: Ensure the refresh token is not expired
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token is expired")
	}

	// Step 3: Generate a new access token and refresh token
	newTokens, err := j.GenerateTokens(claims.UserID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Step 4: Return both tokens
	return &AuthResponse{
		Token:        newTokens.Token,        // New access token
		RefreshToken: newTokens.RefreshToken, // New refresh token
	}, nil
}

func (j *JWTService) generateToken(claims *UserClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.secret))
}

// RevokeToken adds a token to the revoked list an expiration time
func (j *JWTService) RevokeToken(token string, expiresAt time.Time) {
	j.mu.Lock()         // Lock the mutex
	defer j.mu.Unlock() //Ensure it's unlocked when the function exits

	j.revokedTokens[token] = expiresAt
}

func (j *JWTService) IsTokenRevoked(token string) bool {
	j.mu.Lock()
	defer j.mu.Unlock()

	expiry, exists := j.revokedTokens[token]
	if !exists {
		return false
	}

	if time.Now().After(expiry) {
		delete(j.revokedTokens, token)
		return false
	}

	return true
}
