package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/ksel172/Meduza/teamserver/models"
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
	var listenerConfig services.ListenerConfig
	if err := ctx.ShouldBindJSON(&listenerConfig); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())

	}
	err := h.service.AddListener(ctx, listenerConfig)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error adding listener", err.Error())
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
	var listenerConfig services.ListenerConfig
	if err := ctx.ShouldBindJSON(&listenerConfig); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())

	}

	err := h.service.UpdateListener(ctx, listenerConfig)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to update listener with ID: ", listenerConfig.ID), err.Error())
	}

	models.ResponseSuccess(ctx, http.StatusOK, fmt.Sprintf("Successfully updated listener with ID: ", listenerConfig.ID), nil)
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
