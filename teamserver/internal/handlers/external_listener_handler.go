package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// Entrypoint for external listener requests

// All external listeners call back to this server
// The server is responsible for receiving listener C2 requests
type ExternalController struct {
	agentDAL          dal.IAgentDAL
	listenerDAL       dal.IListenerDAL
	checkinController checkin.CheckInController
}

func NewExternalController(agentDal dal.IAgentDAL, listenerDal dal.IListenerDAL, checkinController checkin.CheckInController) *ExternalController {
	return &ExternalController{
		agentDAL:          agentDal,
		listenerDAL:       listenerDal,
		checkinController: checkinController,
	}
}

// TODO: We need to figure out how exactly to handle the registration process of listeners.
// First we need to save the parameters and only then can we create a config.
// The most viable approach for now seems to be to create an empty config and have the
// listener not be able to start/stop until it is updated with a config.

// We can also save the parameters under a different table but that can and most likely
// will create routing problems when creating multiple listeners of the same external kind.

func (ec *ExternalController) RegisterListener(ctx *gin.Context) {
	var listenerModel models.Listener

	if err := ctx.ShouldBindJSON(&listenerModel); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Invalid request body", err.Error())
		return
	}

	// TODO: Move all data validation as methods to the model in models
	if listenerModel.Host == "" || listenerModel.Port < 1 || listenerModel.Port > 65535 {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing or invalid required fields", "Host and Port are required. Port must be between 1 and 65535")
		return
	}

	// Make sure the listener is marked as external, there is no way a listener registered this way
	// isn't external
	listenerModel.External = true

	err := ec.listenerDAL.CreateListener(ctx, &listenerModel)
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
func (ec *ExternalController) HandleTaskRequest(ctx *gin.Context) {

	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Temporary X-Session-Token header
	sessionToken := ctx.Request.Header.Get("X-Session-Token")
	if sessionToken == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "missing session token header", "")
		return
	}

	c2response, err := ec.checkinController.HandleTaskRequest(ctx.Request.Context(), c2request, sessionToken)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "failed to handle task request", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Task request processed successfully", c2response)
}

// HandleResponseSubmission processes the response from a task executed by an agent
func (ec *ExternalController) HandleResponseSubmission(ctx *gin.Context) {
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
func (es *ExternalController) HandleAgentRegistration(ctx *gin.Context) {

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
