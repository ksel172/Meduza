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

func (m *MockAgentDAL) CreateAgentConfig(ctx context.Context, agentConfig models.AgentConfig) error {
	args := m.Called(agentConfig)
	return args.Error(0)
}

func (m *MockAgentDAL) GetAgentConfig(ctx context.Context, agentID string) (models.AgentConfig, error) {
	args := m.Called(agentID)
	return args.Get(0).(models.AgentConfig), args.Error(1)
}

func (m *MockAgentDAL) UpdateAgentConfig(ctx context.Context, agentID string, newConfig models.AgentConfig) error {
	args := m.Called(agentID)
	return args.Error(0)
}

func (m *MockAgentDAL) DeleteAgentConfig(ctx context.Context, agentID string) error {
	args := m.Called(agentID)
	return args.Error(0)
}

func (m *MockAgentDAL) CreateAgentInfo(ctx context.Context, agent models.AgentInfo) error {
	args := m.Called(agent)
	return args.Error(0)
}

func (m *MockAgentDAL) UpdateAgentInfo(ctx context.Context, agent models.AgentInfo) error {
	args := m.Called(agent)
	return args.Error(0)
}

func (m *MockAgentDAL) GetAgentInfo(ctx context.Context, agentID string) (models.AgentInfo, error) {
	args := m.Called(agentID)
	return args.Get(0).(models.AgentInfo), args.Error(1)
}

func (m *MockAgentDAL) DeleteAgentInfo(ctx context.Context, agentID string) error {
	args := m.Called(agentID)
	return args.Error(0)
}

func (m *MockAgentDAL) UpdateAgentLastCallback(ctx context.Context, agentID string, lastCallback string) error {
	args := m.Called(agentID, lastCallback)
	return args.Error(0)
}

func (m *MockAgentDAL) UpdateAgentTask(ctx context.Context, task models.AgentTask) error {
	args := m.Called(task)
	return args.Error(0)
}

type MockCheckInDal struct {
	mock.Mock
}

func (m *MockCheckInDal) CreateAgent(ctx context.Context, agent models.Agent) error {
	args := m.Called(agent)
	return args.Error(0)
}

func (m *MockCheckInDal) GetAgent(ctx context.Context, agentID string) (models.Agent, error) {
	args := m.Called(agentID)
	return args.Get(0).(models.Agent), args.Error(1)
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

type MockListenerDAL struct {
	mock.Mock
}

func (m *MockListenerDAL) CreateListener(ctx context.Context, listener *models.Listener) error {
	args := m.Called(listener)
	return args.Error(0)
}

func (m *MockListenerDAL) GetListenerById(ctx context.Context, lid string) (models.Listener, error) {
	args := m.Called(lid)
	return args.Get(0).(models.Listener), args.Error(1)
}

func (m *MockListenerDAL) GetAllListeners(ctx context.Context) ([]models.Listener, error) {
	args := m.Called()
	return args.Get(0).([]models.Listener), args.Error(1)
}

func (m *MockListenerDAL) DeleteListener(ctx context.Context, lid string) error {
	args := m.Called(lid)
	return args.Error(0)
}

func (m *MockListenerDAL) UpdateListener(ctx context.Context, lid string, updates map[string]any) error {
	args := m.Called(lid, updates)
	return args.Error(0)
}

type MockAdminDal struct {
	mock.Mock
}

func (m *MockAdminDal) CreateDefaultAdmins(ctx context.Context, admin *models.ResAdmin) error {
	args := m.Called(admin)
	return args.Error(0)
}
