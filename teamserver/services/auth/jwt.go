package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService handles token operations
type JWTService struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(secret string, accessTokenTTL, refreshTokenTTL time.Duration) *JWTService {
	return &JWTService{
		secret:          secret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
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
	accessClaims := &UserClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
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
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid claims or token")
	}

	return claims, nil
}

// RefreshTokens generates a new access token if the refresh token is valid
func (j *JWTService) RefreshTokens(refreshToken string) (*AuthRefreshToken, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token is expired")
	}

	newTokens, err := j.GenerateTokens(claims.UserID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token: %w", err)
	}

	return &AuthRefreshToken{
		RefreshToken: newTokens.RefreshToken,
	}, nil

}

func (j *JWTService) generateToken(claims *UserClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.secret))
}
