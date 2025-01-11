package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type PayloadHandler struct {
	agentDAL    dal.IAgentDAL
	listenerDAL dal.IListenerDAL
	payloadDAL  dal.IPayloadDAL
}

func NewPayloadHandler(agentDAL dal.IAgentDAL, listenerDAL dal.IListenerDAL, payloadDAL dal.IPayloadDAL) *PayloadHandler {
	return &PayloadHandler{
		agentDAL:    agentDAL,
		listenerDAL: listenerDAL,
		payloadDAL:  payloadDAL,
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
	payloadConfig.PayloadID = uuid.New().String()
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
	baseconfPath := "./agent/Agent/baseconf.json"
	// Write configuration to a JSON file
	err = ioutil.WriteFile(baseconfPath, file, 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to write JSON configuration file.",
		})
		logger.Error("Error writing JSON file:", err)
		return
	}

	// Run Docker container to compile the agent
	// arch := payloadConfig.Arch
	args := []string{
		"publish",
		"--configuration", "Release",
		"--self-contained", "true",
		"-o", "/app/build/agent-" + payloadConfig.PayloadID,
		"-p:PublishSingleFile=true",
		"-r", payloadConfig.Arch, // interesting fix here
		"agent/Agent/Agent.csproj",
	}

	// logger.Info("Compiling the agent with the following arguments:", args)

	// Prepend the dotnet executable
	cmd := exec.Command("dotnet", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to compile the payload.",
		})
		logger.Error("Error running Docker container to compile agent:", err)
		return
	}

	// Save the payload configuration in the database
	if err := h.payloadDAL.CreatePayload(ctx.Request.Context(), payloadConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to save payload configuration.",
		})
		logger.Error("Error saving payload configuration:", err)
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

	truncErr := os.Truncate(baseconfPath, 0)
	if truncErr != nil {
		logger.Error("Error cleaning baseconf.json:", truncErr)
	}

	// TODO: Code payload DAL to save the payload
	// Clean up compilation and improve error handling and logging

	ctx.JSON(http.StatusOK, gin.H{
		"status":  s.SUCCESS,
		"message": "Payload created and compiled successfully.",
	})
}

func (h *PayloadHandler) DeletePayload(ctx *gin.Context) {
	payloadId := ctx.Param("id")
	filePath := "./teamserver/build/agent-" + payloadId

	logger.Info(filePath)
	err := h.payloadDAL.DeletePayload(ctx.Request.Context(), payloadId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to delete payload.",
		})
		logger.Error("Error deleting payload:", err)
		return
	}

	// Delete the payload folder and all its contents
	err = os.RemoveAll(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to delete payload folder.",
		})
		logger.Error("Error deleting payload folder:", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *PayloadHandler) DeleteAllPayloads(ctx *gin.Context) {
	dirPath := "./teamserver/build"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to read payload directory.",
		})
		logger.Error("Error reading payload directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "agent-") {
			payloadId := file.Name()
			filePath := filepath.Join(dirPath, payloadId)

			// Delete the payload directory and all its contents
			err = os.RemoveAll(filePath)
			if err != nil {
				logger.Error("Error deleting payload directory:", err)
				continue
			}
		}
	}

	// Delete the payloads from the database
	delErr := h.payloadDAL.DeleteAllPayloads(ctx.Request.Context())
	if err != nil {
		logger.Error("Error deleting payloads:", delErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  s.SUCCESS,
		"message": "All payloads deleted successfully.",
	})
}

func (h *PayloadHandler) DownloadPayload(ctx *gin.Context) {
	payloadId := ctx.Param("id")
	extensions := []string{".exe", ".bin", ".dll", ""} // keeping extensions here for now,
	// maybe later move them out to somewhere more modifiable like the .env file or at least the payload models
	var executablePath string
	found := false

	for _, ext := range extensions {
		executablePath = fmt.Sprintf("teamserver/build/agent-%s/agent%s", payloadId, ext)
		if _, err := os.Stat(executablePath); err == nil {
			found = true
			break
		}
	}

	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  s.FAILED,
			"message": "Executable not found.",
		})
		logger.Error("Executable not found")
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(executablePath)))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(executablePath)
}

func (h *PayloadHandler) GetAllPayloads(ctx *gin.Context) {

	payloads, err := h.payloadDAL.GetAllPayloads(ctx.Request.Context())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  s.FAILED,
			"message": "Failed to get payloads.",
		})
		logger.Error("Error getting payloads:", err)
		return
	}
	ctx.JSON(http.StatusOK, payloads)
}
