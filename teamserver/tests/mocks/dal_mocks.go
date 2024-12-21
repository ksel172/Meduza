package mocks

import (
	"context"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/mock"
)

type MockAgentDAL struct {
	mock.Mock
}

func (m *MockAgentDAL) GetAgent(agentID string) (models.Agent, error) {
	args := m.Called(agentID)
	return args.Get(0).(models.Agent), args.Error(1)
}

func (m *MockAgentDAL) UpdateAgent(ctx context.Context, agent models.UpdateAgentRequest) (models.Agent, error) {
	args := m.Called(agent)
	return args.Get(0).(models.Agent), args.Error(1)
}

func (m *MockAgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	args := m.Called(agentID)
	return args.Error(0)
}

func (m *MockAgentDAL) CreateAgentTask(ctx context.Context, task models.AgentTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockAgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	args := m.Called(agentID)
	return args.Get(0).([]models.AgentTask), args.Error(1)
}

func (m *MockAgentDAL) DeleteAgentTask(ctx context.Context, agentID string, taskID string) error {
	args := m.Called(agentID, taskID)
	return args.Error(0)
}

func (m *MockAgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	args := m.Called(agentID)
	return args.Error(0)
}

type MockCheckInDal struct {
	mock.Mock
}

func (m *MockCheckInDal) CreateAgent(ctx context.Context, agent models.Agent) error {
	args := m.Called(agent)
	return args.Error(0)
}

type MockUserDAL struct {
	mock.Mock
}

func (m *MockUserDAL) AddUsers(ctx context.Context, user *models.ResUser) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDAL) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserDAL) GetUserByUsername(ctx context.Context, username string) (*models.ResUser, error) {
	args := m.Called(username)
	return args.Get(0).(*models.ResUser), args.Error(1)
}

func (m *MockUserDAL) GetUserById(ctx context.Context, id string) (*models.ResUser, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ResUser), args.Error(1)
}
