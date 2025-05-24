package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
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

	if err := listenerModel.Validate(); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener data", err.Error())
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

	// Temporary X-Session-Token header
	sessionToken := ctx.Request.Header.Get("X-Session-Token")
	if sessionToken == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "missing session token header", "")
		return
	}

	err := ec.checkinController.HandleResponseRequest(ctx, c2request)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "failed to handle response submission", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Response submission processed successfully", nil)
}

// HandleAgentRegistration handles the registration of an agent
func (ec *ExternalController) HandleAgentRegistration(ctx *gin.Context) {
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

	// Temporary X-Auth-Token header
	payloadToken := ctx.Request.Header.Get("X-Auth-Token")
	if payloadToken == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "missing payload token header", "")
		return
	}

	_, err := ec.checkinController.HandleRegisterRequest(ctx, c2request, payloadToken)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "failed to handle agent registration", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Agent registration processed successfully", nil)
}
