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
	args := m.Called(ctx, agent)
	return args.Get(0).(models.Agent), args.Error(1)
}

func (m *MockAgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	args := m.Called(ctx, agentID)
	return args.Error(0)
}

func (m *MockAgentDAL) CreateAgentTask(ctx context.Context, task models.AgentTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockAgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]models.AgentTask), args.Error(1)
}

func (m *MockAgentDAL) DeleteAgentTask(ctx context.Context, agentID string, taskID string) error {
	args := m.Called(ctx, agentID, taskID)
	return args.Error(0)
}

func (m *MockAgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	args := m.Called(ctx, agentID)
	return args.Error(0)
}