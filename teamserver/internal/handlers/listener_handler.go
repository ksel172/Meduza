package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type ListenerHandler struct {
	service *services.ListenersService
}

func NewListenersHandler(service *services.ListenersService) *ListenerHandler {
	return &ListenerHandler{
		service: service,
	}
}

func (h *ListenerHandler) CreateListener(ctx *gin.Context) {
	var listener services.Listener
	if err := ctx.ShouldBindJSON(&listener); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	existingListener, err := h.service.GetListenerByName(ctx.Request.Context(), listener.Name)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to check listener existence", err.Error())
		return
	}
	if existingListener != nil {
		models.ResponseError(ctx, http.StatusConflict, "Failed to create listener", "Listener with the same name already exists")
		return
	}

	// Validate the listener configuration
	switch listener.Type {
	case "http":
		var httpConfig services.HttpListenerConfig
		if err := utils.MapToStruct(listener.Config, &httpConfig); err != nil {
			models.ResponseError(ctx, http.StatusBadRequest, "Failed to validate listener config", "Invalid configuration for HTTP listener")
			return
		}
		if err := httpConfig.Validate(); err != nil {
			models.ResponseError(ctx, http.StatusBadRequest, "Failed to validate listener config", err.Error())
			return
		}
		listener.Config = httpConfig
	}

	err = h.service.AddListener(ctx.Request.Context(), &listener)
	if err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to create listener", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "Listener created successfully", nil)
}

func (h *ListenerHandler) GetAllListeners(ctx *gin.Context) {

	listeners, err := h.service.GetListeners(ctx)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error getting listeners", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listeners retrieved successfully", listeners)
}

func (h *ListenerHandler) GetListener(ctx *gin.Context) {
	listenerID := ctx.Param(services.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}
	listener, err := h.service.GetListener(ctx, listenerID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error getting listener", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener retrieved successfully", listener)
}

func (h *ListenerHandler) TerminateListener(ctx *gin.Context) {
	listenerID := ctx.Param(services.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	err := h.service.TerminateListener(ctx, listenerID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error deleting listener", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener deleted successfully", nil)
}

func (h *ListenerHandler) UpdateListener(ctx *gin.Context) {
	// Get the listener ID from path parameter
	listenerID := ctx.Param(services.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	// Parse the incoming data as a Listener struct
	var listener services.Listener
	if err := ctx.ShouldBindJSON(&listener); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	// Ensure the listener ID in the path matches the one in the request body
	if listener.ID == "" {
		listener.ID = listenerID
	} else if listener.ID != listenerID {
		models.ResponseError(ctx, http.StatusBadRequest, "Listener ID mismatch",
			"Listener ID in URL doesn't match ID in request body")
		return
	}

	// Update the listener
	err := h.service.UpdateListener(ctx, &listener)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to update listener with ID: %s", listener.ID), err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK,
		fmt.Sprintf("Successfully updated listener with ID: %s", listener.ID), nil)
}

func (h *ListenerHandler) StartListener(ctx *gin.Context) {
	listenerID := ctx.Param(services.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	// I am unsure how we should implement this yet
	err := h.service.StartListener(ctx, listenerID, make(chan<- error))
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error starting listener", err.Error())
	}

	models.ResponseSuccess(ctx, http.StatusOK, fmt.Sprintf("Successfully started listener with ID: %s", listenerID), nil)
}

func (h *ListenerHandler) StopListener(ctx *gin.Context) {
	listenerID := ctx.Param(services.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	// I am unsure how we should implement this yet
	err := h.service.StopListener(ctx, listenerID, make(chan<- error))
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error stopping listener", err.Error())
	}

	models.ResponseSuccess(ctx, http.StatusOK, fmt.Sprintf("Successfully stopped listener with ID: %s", listenerID), nil)
}

func (h *ListenerHandler) GetListenerStatuses(ctx *gin.Context) {

	listenerStatuses, err := h.service.GetListenerStatuses(ctx)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed retrieving listener statuses", err.Error())
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Successfully retrieved listener statuses", listenerStatuses)
}

func (h *ListenerHandler) AutoStart(ctx context.Context) error {

	err := h.service.AutoStart(ctx)
	if err != nil {
		return fmt.Errorf("error starting listeners: %v", err)
	}
	return nil
}
