package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	var payloadRequest models.PayloadRequest

	// Parse and validate the request body
	if err := ctx.ShouldBindJSON(&payloadRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body. Please enter correct input.",
			"status":  s.ERROR,
		})
		logger.Error("Request body error while binding the JSON:", err)
		return
	}

	// Ensure listenerDAL is initialized
	if h.listenerDAL == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Internal server error: listenerDAL is not initialized.",
		})
		logger.Error("listenerDAL is nil")
		return
	}

	// Retrieve the listener configuration
	listener, err := h.listenerDAL.GetListenerById(ctx.Request.Context(), payloadRequest.ListenerID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  s.FAILED,
			"message": "Listener does not exist.",
		})
		logger.Error("Error retrieving the listener:", err)
		return
	}

	// Create payload configuration
	payloadConfig := models.IntoPayloadConfig(payloadRequest)
	payloadConfig.ConfigID = uuid.New().String()
	payloadConfig.ListenerConfig = listener.Config

	file, err := json.MarshalIndent(payloadConfig, "", "  ")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to serialize payload config to JSON.",
		})
		logger.Error("Error marshalling payload config to JSON:", err)
		return
	}

	// Write configuration to a JSON file
	err = ioutil.WriteFile("./agent/Agent/baseconf.json", file, 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to write JSON configuration file.",
		})
		logger.Error("Error writing JSON file:", err)
		return
	}

	// Define paths for Docker container
	// sourcePath := "C:/Users/MagicMan/Documents/Golang/Meduza/agent"
	// outputPath := "C:/Users/MagicMan/Downloads/meduza-publish:/app/output"

	// Run Docker container to compile the agent
	dockerCmd := exec.Command(
		"dotnet", "publish", "--configuration", "Release", "--self-contained", "true", "-o", "/app/output", "-p:PublishSingleFile=true", "-r", "win-x64", "agent/Agent/Agent.csproj",
	)

	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	if err := dockerCmd.Run(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to compile the C# agent.",
		})
		logger.Error("Error running Docker container to compile agent:", err)
		return
	}

	// Save the agent configuration in the database
	agentConfig := models.IntoAgentConfig(payloadConfig)
	if err := h.agentDAL.CreateAgentConfig(ctx.Request.Context(), agentConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to save agent configuration.",
		})
		logger.Error("Error saving agent configuration:", err)
		return
	}

	// TODO: Clean the written temp baseconf.json in the agent dir and avoid creation of unnecessary folders such as /agent/output or /teamserver/agent
	// Add compile types and make payloads save under a specific directory (probably by adding something like payload names)
	// Clean up compilation and improve error handling and logging

	ctx.JSON(http.StatusOK, gin.H{
		"status":  s.SUCCESS,
		"message": "Payload created and compiled successfully.",
	})
}

func (h *PayloadHandler) DeletePayload(ctx *gin.Context) {

}

func (h *PayloadHandler) DownloadPayload(ctx *gin.Context) {

}

func (h *PayloadHandler) GetAllPayloads(ctx *gin.Context) {

}
