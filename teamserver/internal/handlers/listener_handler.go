package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"reflect"
// 	"strings"
// 	"sync"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// 	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
// 	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
// 	"github.com/ksel172/Meduza/teamserver/models"
// 	"github.com/ksel172/Meduza/teamserver/pkg/logger"
// 	"github.com/ksel172/Meduza/teamserver/utils"
// )

// type ListenerHandler struct {
// 	dal     dal.IListenerDAL
// 	service *services.ListenersService
// }

// func NewListenersHandler(dal dal.IListenerDAL, service *services.ListenersService) *ListenerHandler {
// 	return &ListenerHandler{
// 		dal:     dal,
// 		service: service,
// 	}
// }

// func (h *ListenerHandler) CreateListener(ctx *gin.Context) {
// 	var listener models.Listener
// 	if err := ctx.ShouldBindJSON(&listener); err != nil {
// 		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
// 		return
// 	}

// 	existingListener, err := h.dal.GetListenerByName(ctx.Request.Context(), listener.Name)
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to check listener existence", err.Error())
// 		return
// 	}
// 	if existingListener.ID != uuid.Nil {
// 		models.ResponseError(ctx, http.StatusConflict, "Failed to create listener", "Listener with the same name already exists")
// 		return
// 	}

// 	// Validate the listener configuration
// 	switch listener.Type {
// 	case models.ListenerTypeHTTP:
// 		var httpConfig models.HTTPListenerConfig
// 		if err := utils.MapToStruct(listener.Config, &httpConfig); err != nil {
// 			models.ResponseError(ctx, http.StatusBadRequest, "Failed to validate listener config", "Invalid configuration for HTTP listener")
// 			return
// 		}
// 		if err := httpConfig.Validate(); err != nil {
// 			models.ResponseError(ctx, http.StatusBadRequest, "Failed to validate listener config", err.Error())
// 			return
// 		}
// 		listener.Config = httpConfig
// 	}

// 	err = h.dal.CreateListener(ctx.Request.Context(), &listener)
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusBadRequest, "Failed to create listener", err.Error())
// 		return
// 	}

// 	models.ResponseSuccess(ctx, http.StatusCreated, "Listener created successfully", nil)
// }

// func (h *ListenerHandler) GetAllListeners(ctx *gin.Context) {
// 	lists, err := h.dal.GetAllListeners(ctx.Request.Context())
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get listeners", err.Error())
// 		logger.Error("Error getting all listeners:", err)
// 		return
// 	}
// 	models.ResponseSuccess(ctx, http.StatusOK, "Listeners retrieved successfully", lists)
// }

// func (h *ListenerHandler) GetListenerById(ctx *gin.Context) {
// 	id := ctx.Param(models.ParamListenerID)
// 	listener, err := h.dal.GetListenerById(ctx.Request.Context(), id)
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusNotFound, "Failed to get listener", "Listener does not exist")
// 		logger.Error("Error getting listener:", err)
// 		return
// 	}

// 	models.ResponseSuccess(ctx, http.StatusOK, "Listener retrieved successfully", listener)
// }

// func (h *ListenerHandler) DeleteListener(ctx *gin.Context) {
// 	id := ctx.Param(models.ParamListenerID)
// 	c := ctx.Request.Context()

// 	_, err := h.dal.GetListenerById(c, id)
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusNotFound, "Failed to get listener", "Listener does not exist")
// 		logger.Error("Error getting listener:", err)
// 		return
// 	}

// 	err = h.dal.DeleteListener(c, id)
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusBadRequest, "Failed to delete listener", err.Error())
// 		logger.Error("Error deleting listener:", err)
// 		return
// 	}

// 	models.ResponseSuccess(ctx, http.StatusOK, "Listener deleted successfully", gin.H{"id": id})
// }

// func (h *ListenerHandler) UpdateListener(ctx *gin.Context) {
// 	id := ctx.Param(models.ParamListenerID)
// 	c := ctx.Request.Context()

// 	var listener models.Listener
// 	if err := ctx.ShouldBindJSON(&listener); err != nil {
// 		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
// 		logger.Error("Request body error:", err)
// 		return
// 	}

// 	exists, err := h.dal.GetListenerById(c, id)
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusNotFound, "Failed to get listener", "Listener does not exist")
// 		logger.Error("Error getting listener:", err)
// 		return
// 	}

// 	if listener.Type != "" && listener.Type != exists.Type {
// 		models.ResponseError(ctx, http.StatusForbidden, "Failed to update listener", "Updating the 'type' field is not allowed")
// 		logger.Warn("Unauthorized attempt to update 'type' field")
// 		return
// 	}

// 	updates := make(map[string]any)

// 	if listener.Name != "" && listener.Name != exists.Name {
// 		updates["name"] = listener.Name
// 	}
// 	if listener.Status != exists.Status {
// 		updates["status"] = listener.Status
// 	}
// 	if listener.Description != "" && listener.Description != exists.Description {
// 		updates["description"] = listener.Description
// 	}

// 	if listener.Config != nil && !reflect.DeepEqual(listener.Config, exists.Config) {
// 		parsedConfig, err := services.ValidateAndParseConfig(exists.Type, listener.Config)
// 		if err != nil {
// 			models.ResponseError(ctx, http.StatusBadRequest, "Failed to validate listener config",
// 				fmt.Sprintf("Invalid config for listener type '%s': %v", exists.Type, err))
// 			logger.Warn("Invalid config for listener type:", exists.Type, err)
// 			return
// 		}
// 		configJson, err := json.Marshal(parsedConfig)
// 		if err != nil {
// 			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to process listener config", "Failed to marshal config field")
// 			return
// 		}
// 		updates["config"] = configJson
// 	}

// 	if listener.Logging != (models.Logging{}) && !reflect.DeepEqual(listener.Logging, exists.Logging) {
// 		logJson, err := json.Marshal(&listener.Logging)
// 		if err != nil {
// 			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to process listener logging", "Failed to marshal logging field")
// 			return
// 		}
// 		updates["logging"] = logJson
// 	}

// 	if listener.LoggingEnabled != exists.LoggingEnabled {
// 		updates["logging_enabled"] = listener.LoggingEnabled
// 	}

// 	if len(updates) > 0 {
// 		updates["updated_at"] = time.Now().UTC()
// 		if err := h.dal.UpdateListener(c, id, updates); err != nil {
// 			logger.Error("Failed to update listener:", err)
// 			models.ResponseError(ctx, http.StatusBadRequest, "Failed to update listener", err.Error())
// 			return
// 		}
// 		models.ResponseSuccess(ctx, http.StatusOK, "Listener updated successfully", nil)
// 	} else {
// 		models.ResponseSuccess(ctx, http.StatusOK, "No listener fields to update", nil)
// 	}
// }

// func (h *ListenerHandler) StartListener(ctx *gin.Context) {
// 	c := ctx.Request.Context()
// 	id := ctx.Param(models.ParamListenerID)

// 	listener, err := h.dal.GetListenerById(c, id)
// 	if err != nil {
// 		models.ResponseError(ctx, http.StatusNotFound, "Failed to get listener", "Listener does not exist")
// 		logger.Error("Failed to retrieve listener from database:", err)
// 		return
// 	}

// 	logger.Info("Attempting to create listener of type:", listener.Type)
// 	if err := h.service.CreateListenerController(listener.Type, listener.Config); err != nil {
// 		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create listener controller", err.Error())
// 		logger.Error("Failed to create listener controller:", err)
// 		return
// 	}

// 	logger.Info("Starting listener with ID:", id)
// 	if err := h.service.Start(listener); err != nil {
// 		logger.Error("Failed to start the listener:", err)
// 		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to start listener", err.Error())
// 		return
// 	}

// 	err = h.dal.UpdateListener(c, id, map[string]any{"status": 1})
// 	if err != nil {
// 		logger.Error("Failed to update listener status:", err)
// 		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update listener status", err.Error())
// 		return
// 	}

// 	models.ResponseSuccess(ctx, http.StatusOK, "Listener started successfully", gin.H{"id": id})
// }

// func (h *ListenerHandler) StopListener(ctx *gin.Context) {
// 	c := ctx.Request.Context()
// 	id := ctx.Param(models.ParamListenerID)

// 	if err := h.service.Stop(id, 10*time.Second); err != nil {
// 		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to stop listener", err.Error())
// 		return
// 	}

// 	err := h.dal.UpdateListener(c, id, map[string]any{"status": 2})
// 	if err != nil {
// 		logger.Error("Failed to update listener status:", err)
// 		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update listener status", err.Error())
// 		return
// 	}

// 	models.ResponseSuccess(ctx, http.StatusOK, "Listener stopped successfully", nil)
// }

// func (h *ListenerHandler) CheckRunningListener(ctx *gin.Context) {
// 	id := ctx.Param(models.ParamListenerID)

// 	_, running := h.service.GetListener(id)
// 	if !running {
// 		models.ResponseError(ctx, http.StatusNotFound, "Failed to check listener status", "Listener is not running")
// 		return
// 	}

// 	models.ResponseSuccess(ctx, http.StatusOK, "Listener is running", nil)
// }

// func (h *ListenerHandler) AutoStart(ctx context.Context) error {
// 	activeListeners, err := h.dal.GetActiveListeners(ctx)
// 	if err != nil {
// 		logger.Error("Error fetching active listeners:", err)
// 		return fmt.Errorf("failed to fetch active listeners: %w", err)
// 	}

// 	totalActiveListeners := len(activeListeners)
// 	if totalActiveListeners == 0 {
// 		logger.Info("No active listeners found to start.")
// 		return nil
// 	}

// 	logger.Info("Found", totalActiveListeners, "listeners to start.")

// 	var wg sync.WaitGroup
// 	errChan := make(chan error, totalActiveListeners)

// 	if h.service == nil {
// 		return fmt.Errorf("listener service is nil")
// 	}

// 	for _, listener := range activeListeners {
// 		wg.Add(1)
// 		go func(listener models.Listener) {
// 			defer wg.Done()

// 			id := listener.ID.String()
// 			logger.Info("Starting listener:", id)

// 			if err := h.service.CreateListenerController(listener.Type, listener.Config); err != nil {
// 				logger.Error("Error creating listener controller", id, err)
// 				errChan <- fmt.Errorf("failed to create listener controller %s: %w", id, err)
// 				return
// 			}

// 			if err := h.service.Start(listener); err != nil {
// 				logger.Error("Error starting listener", id, err)
// 				errChan <- fmt.Errorf("failed to start listener %s: %w", id, err)
// 			} else {
// 				logger.Info("Listener started successfully:", id)
// 			}
// 		}(listener)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(errChan)
// 	}()

// 	var errors []string
// 	for err := range errChan {
// 		if err != nil {
// 			errors = append(errors, err.Error())
// 		}
// 	}

// 	if len(errors) > 0 {
// 		return fmt.Errorf("autostart encountered errors:\n%s", strings.Join(errors, "\n"))
// 	}

// 	logger.Info("All listeners started successfully.")
// 	return nil
// }
