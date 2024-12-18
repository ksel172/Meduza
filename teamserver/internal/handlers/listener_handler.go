package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/pkg/listeners"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type ListenerHandler struct {
	dal *dal.ListenerDAL
}

func NewListenersHandler(dal *dal.ListenerDAL) *ListenerHandler {
	return &ListenerHandler{dal: dal}
}

func (h *ListenerHandler) CreateListener(ctx *gin.Context) {
	var listener listeners.Listener

	if err := ctx.ShouldBindJSON(&listener); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "Invalid Request body.Please type correct input",
			"status":  s.ERROR,
		})
		logger.Error("Request Body Error while bind the json:\n", err)
		return
	}

	reqCtx := ctx.Request.Context()

	// Convert the parsed configuration back to JSON
	configJSON, err := json.Marshal(listener.Config)
	if err != nil {
		logger.Error("Error converting parsed config to JSON:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error processing configuration.",
			"status":  s.ERROR,
		})
		return
	}
	listener.Config = configJSON

	err = h.dal.CreateListener(reqCtx, &listener)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  s.ERROR,
			"message": "Unable to create a listener.",
		})
		logger.Error("Error Occured while Adding Data to listener:\n", err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  s.SUCCESS,
		"message": "listener created successfully",
	})

}

func (h *ListenerHandler) GetAllListeners(ctx *gin.Context) {
	lists, err := h.dal.GetAllListeners(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Unable to process the request",
		})
		logger.Error("Error Unable to get all the listeners", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": s.SUCCESS,
		"data":   lists,
	})
}

func (h *ListenerHandler) GetListenerById(ctx *gin.Context) {
	id := ctx.Param("id")
	listener, err := h.dal.GetListenerById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  s.FAILED,
			"message": "Listener Does Not exists",
		})
		logger.Error("Error Unable to get the listener", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": s.SUCCESS,
		"data":   listener,
	})
}

func (h *ListenerHandler) DeleteListener(ctx *gin.Context) {
	id := ctx.Param("id")
	c := ctx.Request.Context()

	_, err := h.dal.GetListenerById(c, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  s.FAILED,
			"message": "listener does not exists",
		})
		logger.Error("Error Unable to get the listener", err)
		return
	}
	err = h.dal.DeleteListener(c, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  s.FAILED,
			"message": "unable to delete listener",
		})
		logger.Error("Error Delete listener", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  s.SUCCESS,
		"message": "listener deleted successfully",
		"id":      id,
	})
}

func (h *ListenerHandler) UpdateListener(ctx *gin.Context) {
	id := ctx.Param("id")
	c := ctx.Request.Context()

	var listener listeners.Listener

	if err := ctx.ShouldBindJSON(&listener); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "invalid request body. please type correct input",
			"status":  s.ERROR,
		})
		logger.Error("Request Body Error while bind the json:\n", err)
		return
	}

	exists, err := h.dal.GetListenerById(c, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  s.FAILED,
			"message": "Listener Does Not exists",
		})
		logger.Error("Error Unable to get the listener", err)
		return
	}

	// This check make sure that type cannot be changed
	if listener.Type != "" && listener.Type != exists.Type {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  s.ERROR,
			"message": "Updating the 'type' field is not allowed",
		})
		logger.Warn("Unauthorized attempt to update 'type' field.")
		return
	}

	// Updates stores all of the updated queries
	updates := make(map[string]any)

	if listener.Name != "" && listener.Name != exists.Name {
		updates["name"] = listener.Name
	}
	if listener.Status != exists.Status {
		updates["status"] = listener.Status
	}
	if listener.Description != "" && listener.Description != exists.Description {
		updates["description"] = listener.Description
	}

	//TODO: implement a functionality for config updated based on the type.
	if listener.Config != nil && !reflect.DeepEqual(listener.Config, exists.Config) {
		parsedConfig, err := listeners.ValidateAndParseConfig(exists.Type, listener.Config)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  s.ERROR,
				"message": fmt.Sprintf("Invalid config for listener type '%s': %v", exists.Type, err),
			})
			logger.Warn("Invalid config for listener type:", exists.Type, err)
			return
		}
		configJson, err := json.Marshal(parsedConfig)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  s.FAILED,
				"message": "Failed to marshal config field",
			})
			return
		}
		updates["config"] = configJson
	}
	if listener.ResponseRules.Headers != nil && !reflect.DeepEqual(listener.ResponseRules.Headers, exists.ResponseRules.Headers) {
		// Marshal the ResponseRules field only if it has changed
		rrJson, err := json.Marshal(&listener.ResponseRules)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  s.FAILED,
				"message": "Failed to marshal response rules field",
			})
			return
		}
		updates["response_rules"] = rrJson
	}
	if listener.Logging != (listeners.Logging{}) && !reflect.DeepEqual(listener.Logging, exists.Logging) {
		// Marshal the Logging field only if it has changed
		logJson, err := json.Marshal(&listener.Logging)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  s.FAILED,
				"message": "Failed to marshal logging field",
			})
			return
		}
		updates["logging"] = logJson
	}
	if listener.LoggingEnabled != exists.LoggingEnabled {
		updates["logging_enabled"] = listener.LoggingEnabled
	}

	// Only update the fields that are present in the 'updates' map
	if len(updates) > 0 {
		updates["updated_at"] = time.Now().UTC() // Always update the 'updated_at' field
		if err := h.dal.UpdateListener(c, id, updates); err != nil {
			logger.Error("Failed to update listener\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  s.FAILED,
				"message": "Unable to update listener. Please try again later",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  s.SUCCESS,
			"message": "Listener Updated Successfully",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  s.SUCCESS,
			"message": "No fields to update",
		})
	}
}
