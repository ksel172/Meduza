package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
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
			"status":  utils.Status.ERROR,
		})
		logger.Error("Request body error while binding the JSON:", err)
		return
	}

	// Ensure listenerDAL is initialized
	if h.listenerDAL == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Internal server error: listenerDAL is not initialized.",
		})
		logger.Error("listenerDAL is nil")
		return
	}

	// Retrieve the listener configuration
	listener, err := h.listenerDAL.GetListenerById(ctx.Request.Context(), payloadRequest.ListenerID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  utils.Status.FAILED,
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

	// Generate the public and private keys for the server for this payload
	privateKey, publicKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "failed to generate server ECDH keys",
		})
	}
	payloadConfig.PublicKey = publicKey
	payloadConfig.PrivateKey = privateKey
	payloadConfig.Token = uuid.New().String()

	// TODO: ADD SHARED KEY SHARING WITH AGENT FOR HMAC

	file, err := json.MarshalIndent(payloadConfig, "", "  ")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to serialize payload config to JSON.",
		})
		logger.Error("Error marshalling payload config to JSON:", err)
		return
	}

	baseconfPath := conf.GetBaseConfPath()
	// Write configuration to a JSON file
	err = os.WriteFile(baseconfPath, file, 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to write JSON configuration file.",
		})
		logger.Error("Error writing JSON file:", err)
		return
	}

	args := []string{
		"publish",
		"--configuration", "Release",
		"--self-contained", strings.ToLower(fmt.Sprintf("%t", payloadRequest.SelfContained)),
		"-o", "/app/build/payload-" + payloadConfig.PayloadID,
		"-p:PublishSingleFile=true",
		"-p:DefineConstants=TYPE_" + listener.Type, // Specify comm type to cut out pieces of the code
		"-r", payloadConfig.Arch, // interesting fix here
		"agent/Agent/Agent.csproj",
	}

	cmd := exec.Command("dotnet", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to compile the payload.",
		})
		logger.Error("Error running Docker container to compile agent:", err)
		return
	}

	// Save the payload configuration in the database
	if err := h.payloadDAL.CreatePayload(ctx.Request.Context(), payloadConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to save payload configuration.",
		})
		logger.Error("Error saving payload configuration:", err)
		return
	}

	// Save the agent configuration in the database
	agentConfig := models.IntoAgentConfig(payloadConfig)
	if err := h.agentDAL.CreateAgentConfig(ctx.Request.Context(), agentConfig); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to save agent configuration.",
		})
		logger.Error("Error saving agent configuration:", err)
		return
	}

	defer func() {
		truncErr := os.Truncate(baseconfPath, 0)
		if truncErr != nil {
			logger.Error("Error cleaning baseconf.json:", truncErr)
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  utils.Status.SUCCESS,
		"message": "Payload created and compiled successfully.",
	})
}

func (h *PayloadHandler) DeletePayload(ctx *gin.Context) {
	payloadId := ctx.Param(models.ParamPayloadID)
	filePath := "./teamserver/build/payload-" + payloadId

	logger.Info(filePath)
	err := h.payloadDAL.DeletePayload(ctx.Request.Context(), payloadId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to delete payload.",
		})
		logger.Error("Error deleting payload:", err)
		return
	}

	// Delete the payload folder and all its contents
	err = os.RemoveAll(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to delete payload folder.",
		})
		logger.Error("Error deleting payload folder:", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  utils.Status.SUCCESS,
		"message": "Payload deleted successfully.",
	})
}

func (h *PayloadHandler) DeleteAllPayloads(ctx *gin.Context) {
	dirPath := "./teamserver/build"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.Status.FAILED,
			"message": "Failed to read payload directory.",
		})
		logger.Error("Error reading payload directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "payload-") {
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
		"status":  utils.Status.SUCCESS,
		"message": "All payloads deleted successfully.",
	})
}

func (h *PayloadHandler) DownloadPayload(ctx *gin.Context) {
	payloadId := ctx.Param(models.ParamPayloadID)
	extensions := []string{".exe", ".bin", ".dll", ""} // keeping extensions here for now,
	// maybe later move them out to somewhere more modifiable like the .env file or at least the payload models
	var executablePath string
	found := false

	for _, ext := range extensions {
		executablePath = fmt.Sprintf("teamserver/build/payload-%s/agent%s", payloadId, ext)
		if _, err := os.Stat(executablePath); err == nil {
			found = true
			break
		}
	}

	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  utils.Status.FAILED,
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
			"status":  utils.Status.FAILED,
			"message": "Failed to get payloads.",
		})
		logger.Error("Error getting payloads:", err)
		return
	}
	ctx.JSON(http.StatusOK, payloads)
}
