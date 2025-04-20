package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// Entrypoint for external listener requests

// All external listeners call back to this server
// The server is responsible for receiving listener C2 requests
type ExternalServer struct {
	host        string
	port        int
	server      *gin.Engine
	registry    *listener.ListenerRegistry
	agentDAL    dal.IAgentDAL
	listenerDAL dal.IListenerDAL
}

func NewExternalServer(agentDal dal.IAgentDAL, listenerDal dal.IListenerDAL) *ExternalServer {
	return &ExternalServer{
		agentDAL:    agentDal,
		listenerDAL: listenerDal,
	}
}

// TODO: Implement registration of listener paramaters to dynamically display them on the client application

func (es *ExternalServer) RegisterListener(ctx *gin.Context) {
	var listenerModel models.Listener

	if err := ctx.ShouldBindJSON(&listenerModel); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Invalid request body", err.Error())
		return
	}

	if listenerModel.Host == "" || listenerModel.Port < 1 || listenerModel.Port > 65535 {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing or invalid required fields", "Host and Port are required. Port must be between 1 and 65535")
		return
	}
	// Make sure the listener is marked as external, there is no way a listener registered this way
	// isn't external
	listenerModel.External = true

	err := es.listenerDAL.CreateListener(ctx, &listenerModel)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to create external listener of kind: %s", listenerModel.Kind), err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener controller registered successfully", nil)
}

// The external listener handler processes the incoming requests from the agent
// which were already processed by the external listener. It shouldn't be responsible for
// encryption because the external listener already handles that part. This makes
// the algorithms that are used for encryption and decryption interchangeable.

// HandleTaskRequest handles a task request from an agent
func (es *ExternalServer) HandleTaskRequest(ctx *gin.Context) {

	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	tasks, err := es.agentDAL.GetAgentTasks(ctx, c2request.AgentID)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to get tasks for agent %s: %v", c2request.AgentID, err))
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get tasks", err.Error())
		return
	}

	// Create a new slice for pending tasks
	pendingTasks := make([]models.AgentTask, 0)

	// Only process non-completed tasks
	for _, task := range tasks {
		if task.Status == models.TaskComplete {
			continue
		}

		// The modules will be revived as an agent specific command later.
		// It shouldn't be handled from here because it forces the agent to support modules which
		// can't be done with every codebase.

		// Handle module commands
		// 	if task.Type == models.ModuleCommand {
		// 		moduleDirPath := filepath.Join(conf.GetModuleUploadPath(), task.Module)
		// 		moduleName := task.Command.Name

		// 		modulePath := filepath.Join(moduleDirPath, moduleName)
		// 		mainModuleBytes, err := utils.LoadAssembly(filepath.Join(modulePath, moduleName+".dll"))
		// 		if err != nil {
		// 			logger.Info(fmt.Sprintf("Failed to load main module: %v", err))
		// 			return nil, fmt.Errorf("failed to load main module: %w", ErrInternalServer)
		// 		}

		// 		loadingModulePath := moduleDirPath + "/" + moduleName + "/"
		// 		dependencyBytes := make(map[string][]byte)
		// 		files, err := os.ReadDir(loadingModulePath)
		// 		if err != nil {
		// 			logger.Info(fmt.Sprintf("Failed to read module directory: %v", err))
		// 			return nil, fmt.Errorf("failed to read module directory: %w", ErrInternalServer)
		// 		}

		// 		for _, file := range files {
		// 			if file.Name() != moduleName+".dll" && strings.HasSuffix(file.Name(), ".dll") {
		// 				depBytes, err := utils.LoadAssembly(filepath.Join(loadingModulePath, file.Name()))
		// 				if err != nil {
		// 					logger.Info(fmt.Sprintf("Failed to load main module: %v", err))
		// 					return nil, fmt.Errorf("failed to load dependency :%w", ErrInternalServer)
		// 				}
		// 				dependencyBytes[file.Name()] = depBytes
		// 			}
		// 		}

		// 		moduleBytes := models.ModuleBytes{
		// 			ModuleBytes:     mainModuleBytes,
		// 			DependencyBytes: dependencyBytes,
		// 		}

		// 		moduleBytesJSON, err := json.Marshal(moduleBytes)
		// 		if err != nil {
		// 			logger.Info(fmt.Sprintf("Failed to marshal module bytes: %v", err))
		// 			return nil, fmt.Errorf("failed to marshal module: %w", ErrInternalServer)
		// 		}

		// 		task.Module = base64.StdEncoding.EncodeToString(moduleBytesJSON)
		// 	}

		// Add non-completed task to pending tasks
		pendingTasks = append(pendingTasks, task)
	}

	// Update the agent's last callback time
	lastCallback := time.Now().Format(time.RFC3339)
	if err := es.agentDAL.UpdateAgentLastCallback(ctx, c2request.AgentID, lastCallback); err != nil {
		logger.Info(fmt.Sprintf("Failed to update agent last callback: %v", err))
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update agent last callback", err.Error())
		return
	}

	// Use pendingTasks instead of tasks for the response
	tasksJSON, err := json.Marshal(pendingTasks)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to marshal tasks: %v", err))
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to marshal tasks", err.Error())
		return
	}

	var c2response models.C2Request
	c2response.AgentID = c2request.AgentID
	c2response.Reason = models.Task
	c2response.Message = string(tasksJSON)

	models.ResponseSuccess(ctx, http.StatusOK, "Task request processed successfully", c2response)
}

// HandleResponseSubmission processes the response from a task executed by an agent
func (es *ExternalServer) HandleResponseSubmission(ctx *gin.Context) {
	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	var agentTask models.AgentTask
	if err := json.Unmarshal([]byte(c2request.Message), &agentTask); err != nil {
		logger.Info(fmt.Sprintf("Failed to unmarshal agent message: %v", err))
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to unmarshal agent message", err.Error())
		return
	}

	err := es.agentDAL.UpdateAgentTask(ctx, agentTask)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to update agent task: %v", err))
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update agent task", err.Error())
	}

	logger.Info(fmt.Sprintf("Successfully updated agent task: %s", agentTask.TaskID))
	models.ResponseSuccess(ctx, http.StatusOK, "Response submission processed successfully", nil)
}

// HandleAgentRegistration handles the registration of an agent
func (es *ExternalServer) HandleAgentRegistration(ctx *gin.Context) {

	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	logger.Info(fmt.Sprintf("Received register request from agent: %s", c2request.AgentID))

	var agentInfo models.AgentInfo
	if err := json.Unmarshal([]byte(c2request.Message), &agentInfo); err != nil {
		logger.Info(fmt.Sprintf("Failed to parse agent info from decrypted message: %v", err))
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to parse agent info", err.Error())
		return
	}

	if _, err := es.agentDAL.GetAgent(ctx, agentInfo.AgentID); err == nil {
		logger.Info("Agent already exists:", c2request.AgentID)
		models.ResponseError(ctx, http.StatusConflict, "Agent already exists", nil)
	}

	newAgent := c2request.IntoNewAgent()
	newAgent.Name = utils.RandomString(6)

	if err := es.agentDAL.RegisterAgent(ctx, newAgent); err != nil {
		logger.Info(fmt.Sprintf("Failed to create agent: %v", err))
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create agent", err.Error())
	}

	if err := es.agentDAL.CreateAgentInfo(ctx, agentInfo); err != nil {
		logger.Info(fmt.Sprintf("Failed to create agent info: %v", err))
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create agent info", err.Error())
	}
}
