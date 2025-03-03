package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type ControllerHandler struct {
	controllerDal dal.IControllerDAL
}

func NewControllerHandler(controllerDal dal.IControllerDAL) *ControllerHandler {
	return &ControllerHandler{
		controllerDal: controllerDal,
	}
}

// TODO: After implementing the code for the Controller AuthMiddleware, it has to be moved into the middlewares.go section
// under ./teamserver/internal/server/middlewares.go

func (h *ControllerHandler) ControllerAuthMiddleware() gin.HandlerFunc {
	/*
		return func(ctx *gin.Context) {

			// Not sure where we should be getting the API key from.

			apiKey := ctx.GetHeader("X-API-Key")
			valid, err := h.controllerDal.ValidateAPIKey(ctx, apiKey)

			if err != nil {
				logger.Error("API key validation error:", err)
				models.ResponseError(ctx, http.StatusInternalServerError, "Authentication error", "Failed to validate API key")
				ctx.Abort()
				return
			}

			if !valid {
				logger.Warn("Unauthorized API access attempt with key:", apiKey)
				models.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized", "Invalid API key")
				ctx.Abort()
				return
			}

			ctx.Next()
		}
	*/
	return nil
}

func (h *ControllerHandler) RegisterController(ctx *gin.Context) {
	var registration models.ControllerRegistration

	// Takes what the listener controller sends:
	// ID       string `json:"id"`
	// Endpoint string `json:"endpoint"`
	// and registers a new controller with an ID and endpoint.

	if err := ctx.ShouldBindJSON(&registration); err != nil {
		logger.Error("Invalid controller registration request:", err)
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request", "Request body could not be parsed")
		return
	}

	// Make sure the request isn't null
	if registration.ID == "" || registration.Endpoint == "" {
		logger.Warn("Incomplete controller registration:", registration)
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request", "ID and endpoint are required")
		return
	}

	// Register the controller in the database
	if err := h.controllerDal.RegisterController(ctx, registration); err != nil {
		logger.Error("Failed to register controller:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Registration failed", "Failed to store controller information")
		return
	}

	logger.Info("Controller registered:", registration.ID, "at", registration.Endpoint)
	models.ResponseSuccess(ctx, http.StatusOK, "Controller registered successfully", gin.H{"id": registration.ID})
}

func (h *ControllerHandler) GetKeyPair(ctx *gin.Context) {
	// Not sure where the keys are coming from yet, auth needs to be discussed.
}

func (h *ControllerHandler) ReceiveHeartbeat(ctx *gin.Context) {
	controllerID := ctx.Param(models.ParamControllerID)

	if controllerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request", "Controller ID is required")
		return
	}

	// Binding the heartbeat here:
	// Timestamp int64             `json:"timestamp"`
	// Listeners map[string]string `json:"listeners"`

	var heartbeat models.HeartbeatRequest
	if err := ctx.ShouldBindJSON(&heartbeat); err != nil {
		logger.Error("Invalid heartbeat request:", err)
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request", "Request body could not be parsed")
		return
	}

	exists, err := h.controllerDal.ControllerExists(ctx, controllerID)
	if err != nil {
		logger.Error("Error checking controller existence:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Server error", "Failed to verify controller")
		return
	}

	if !exists {
		models.ResponseError(ctx, http.StatusNotFound, "Not found", "Controller not registered")
		return
	}

	// Update the heartbeat for the listeners in the database
	if err := h.controllerDal.UpdateHeartbeat(ctx, controllerID, heartbeat); err != nil {
		logger.Error("Failed to update heartbeat for controller", controllerID, ":", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Update failed", "Failed to update heartbeat")
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Heartbeat received", nil)
}
