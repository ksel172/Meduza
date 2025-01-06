package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

func New(dal dal.IAgentDAL) *AgentController {
	return &AgentController{
		dal: dal,
	}
}

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
	var payloadRequest models.PayloadRequest

	// Agent config is taken
	if err := ctx.ShouldBindJSON(&payloadRequest); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "Invalid Request body. Please enter correct input",
			"status":  s.ERROR,
		})
		logger.Error("Request Body Error while bind the json:\n", err)
		return
	}

	// Check if listenerDAL is nil
	if h.listenerDAL == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Internal server error",
		})
		logger.Error("listenerDAL is nil")
		return
	}

	// listener is extracted from agent config listener ID
	listener, err := h.listenerDAL.GetListenerById(ctx.Request.Context(), payloadRequest.ListenerID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  s.FAILED,
			"message": "Listener does not exist",
		})
		logger.Error("Error Unable to get the listener", err)
		return
	}

	// Payload config is initialized to be embedded into the payload executable
	var payloadConfig = models.IntoPayloadConfig(payloadRequest)
	payloadConfig.ID = uuid.New().String()
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
	// Payload config is written to a file
	err = ioutil.WriteFile("baseconf.json", file, 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to write JSON file",
		})
		logger.Error("Error writing JSON file", err)
		return
	}

	// TODO: Make payload generation and output and modify the payload to contain only vital
	// information for the embedded config. Also need to make the agentID also assign by making
	// a payload creation type that contains only the necessary data for the API

	// TODO: Need to make the payload also be saved as an agent config in the DB for future
	// Need to resolve issue #58 Add endpoints for agent config management first before taking on
	// the rest so that I could save the config in the DB before writing it to a file

	ctx.JSON(http.StatusOK, gin.H{
		"status":  s.SUCCESS,
		"message": "Payload created successfully",
	})
}

func (h *PayloadHandler) DeletePayload(ctx *gin.Context) {

}

func (h *PayloadHandler) DownloadPayload(ctx *gin.Context) {

}

func (h *PayloadHandler) GetAllPayloads(ctx *gin.Context) {

}
