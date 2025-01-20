package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
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
	return &CheckInController{
		checkInDAL: checkInDAL,
		agentDAL:   agentDAL,
	}
}

func (cc *CheckInController) Checkin(ctx *gin.Context) {

	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
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
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling encryption key request for agent %s", c2request.AgentID))
		cc.handleTaskRequest(ctx, c2request)
	} else if c2request.Reason == models.Response {
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling encryption key request for agent %s", c2request.AgentID))
		cc.handleResponseRequest(ctx, c2request)
	} else if c2request.Reason == models.Register {
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling encryption key request for agent %s", c2request.AgentID))
		cc.handleRegisterRequest(ctx, c2request)
	}
}

func (cc *CheckInController) handleTaskRequest(ctx *gin.Context, c2request models.C2Request) {
	tasks, err := cc.agentDAL.GetAgentTasks(ctx, c2request.AgentID)
	if err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to get tasks for agent %s: %v", c2request.AgentID, err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
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
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to marshal agent task to JSON: %v", err))
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

	logger.Info("Successfully updated agent task:", agentTask.TaskID)
	ctx.JSON(http.StatusOK, "successfully updated")
}

func (cc *CheckInController) handleRegisterRequest(ctx *gin.Context, c2request models.C2Request) {
	logger.Info("Received check-in request from agent:", c2request.AgentID)
	// Parse the message as AgentInfo
	var agentInfo models.AgentInfo
	if err := json.Unmarshal([]byte(c2request.Message), &agentInfo); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to parse agent info from message: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent info"})
		return
	}
	// Check if the agent already exists
	if _, err := cc.agentDAL.GetAgent(agentInfo.AgentID); err == nil {
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Agent already exists: %v", err))
		ctx.JSON(http.StatusConflict, gin.H{"error": "agent already exists"})
		return
	}

	// Create agent in the database
	newAgent := c2request.IntoNewAgent()
	newAgent.Name = utils.RandomString(6)

	if err := cc.checkInDAL.CreateAgent(ctx.Request.Context(), newAgent); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to create agent: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create agent info in the database
	if err := cc.agentDAL.CreateAgentInfo(ctx.Request.Context(), agentInfo); err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to create agent info: %v", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"agent": newAgent})
}
func (cc *CheckInController) handleEncryptionKeyRequest(ctx *gin.Context, c2request models.C2Request) {
	// Generate an AES key for this session to comunicate with the agent
	key, err := utils.GenerateAES256Key()
	if err != nil {
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to generate AES256 key: %v", err))
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Store the key in the database and associate with the agent
	// Implement storage method with an expiry period, redis being the most sensible option
	KeyRegistry.writeKey(c2request.AgentID, key)

	ctx.JSON(http.StatusOK, gin.H{"key": key})
}
