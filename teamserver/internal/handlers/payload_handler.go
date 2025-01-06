package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type PayloadHandler struct {
	agentDAL    dal.IAgentDAL
	listenerDAL dal.IListenerDAL
}

func NewPayloadHandler(agentDAL dal.IAgentDAL, listenerDAL dal.IListenerDAL) *PayloadHandler {
	return &PayloadHandler{
		agentDAL:    agentDAL,
		listenerDAL: listenerDAL,
	}
}

func (h *PayloadHandler) CreatePayload(ctx *gin.Context) {
	var agentConfig models.AgentConfig

	if err := ctx.ShouldBindJSON(&agentConfig); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "Invalid Request body. Please enter correct input",
			"status":  s.ERROR,
		})
		logger.Error("Request Body Error while bind the json:\n", err)
		return
	}

	listener, err := h.listenerDAL.GetListenerById(ctx, agentConfig.ListenerID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  s.FAILED,
			"message": "Listener does not exist",
		})
		logger.Error("Error Unable to get the listener", err)
		return
	}

	// Init payload config
	var payloadConfig = models.IntoPayloadConfig(agentConfig)
	payloadConfig.ListenerConfig = listener.Config

	file, err := json.MarshalIndent(payloadConfig, "", "  ")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to create JSON file",
		})
		logger.Error("Error marshalling payload config to JSON", err)
		return
	}

	err = ioutil.WriteFile("payload_config.json", file, 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to write JSON file",
		})
		logger.Error("Error writing JSON file", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  s.SUCCESS,
		"message": "Payload created successfully",
	})
}

func (h *PayloadHandler) DeletePayload(ctx *gin.Context) {

}

func (h *PayloadHandler) GetAllPayloads(ctx *gin.Context) {

}

func (h *PayloadHandler) DownloadPayload(ctx *gin.Context) {

}
