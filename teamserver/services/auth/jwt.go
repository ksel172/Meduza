package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var secret string

func init() {
	secret = utils.GetEnvString("JWT_SECRET", "jwt_secret")
}

// GenerateToken genrates JSON WEB TOKEN for the loged in user.
// Return an AuthResponse with the token or an error if signing fails.
func GenerateToken(user *models.User) (AuthResponse, error) {
	claims := &UserClaim{
		Id:   user.ID,
		Role: user.Role,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	t, err := token.SignedString(secret)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error while encoding token: %w", err)
	}

	return AuthResponse{Token: t}, nil
}

// ValidateToken parses and verifies the given JWT access token.
// Returns the parsed claims on success or an error on failure.
func ValidateToken(accessToken string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while validating token: %w", err)
	}
	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token or claims")
}

// TokenRefresh generates a new token for an existing valid token.
// Returns the new token as an AuthResponse or an error if validation fails.
func TokenRefresh(accessToken string) (AuthResponse, error) {
	vToken, err := ValidateToken(accessToken)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error validating access token: %w", err)
	}
	if vToken.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return AuthResponse{}, fmt.Errorf("cannot refresh an expired token")
	}
	claims := &UserClaim{
		Id:   vToken.Id,
		Role: vToken.Role,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	t, err := token.SignedString(secret)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error while encoding token: %w", err)
	}

	return AuthResponse{Token: t}, nil
}
