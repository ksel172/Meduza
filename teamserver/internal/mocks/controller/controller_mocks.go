package controller_mocks

import (
	"context"

	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/mock"
)

type MockCheckInController struct {
	mock.Mock
}

func (m *MockCheckInController) Authenticate(agentPublicKey string, authToken string) (checkin.AuthResponse, error) {
	args := m.Called(agentPublicKey, authToken)
	return args.Get(0).(checkin.AuthResponse), args.Error(1)
}
func (m *MockCheckInController) HandleTaskRequest(ctx context.Context, c2request models.C2Request, sessionToken string) ([]byte, error) {
	args := m.Called(c2request, sessionToken)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockCheckInController) HandleResponseRequest(ctx context.Context, c2request models.C2Request) error {
	args := m.Called(c2request)
	return args.Error(0)
}
func (m *MockCheckInController) HandleRegisterRequest(ctx context.Context, c2request models.C2Request) error {
	args := m.Called(c2request)
	return args.Error(0)
}
