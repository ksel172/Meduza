package mocks

import (
	"time"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateTokens(userID, role string) (*models.AuthResponse, error) {
	args := m.Called(userID, role)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenStr string) (*models.UserClaim, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(*models.UserClaim), args.Error(1)
}

func (m *MockJWTService) RefreshTokens(refreshToken string) (*models.AuthResponse, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockJWTService) RevokeToken(token string, expiresAt time.Time) {
	m.Called(token, expiresAt)
}

func (m *MockJWTService) IsTokenRevoked(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}
