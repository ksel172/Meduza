package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type PayloadController struct {
	agentDAL    dal.IAgentDAL
	listenerDAL dal.IListenerDAL
	payloadDAL  dal.IPayloadDAL
}

func NewPayloadController(agentDAL dal.IAgentDAL, listenerDAL dal.IListenerDAL, payloadDAL dal.IPayloadDAL) *PayloadController {
	return &PayloadController{
		agentDAL:    agentDAL,
		listenerDAL: listenerDAL,
		payloadDAL:  payloadDAL,
	}
}

// TODO: added configID to payload request, verify config exists in handler
func (h *PayloadController) CreatePayload(ctx *gin.Context) {
	var payloadRequest models.PayloadRequest

	if err := ctx.ShouldBindJSON(&payloadRequest); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		logger.Error("Request body error while binding the JSON:", err)
		return
	}

	if h.listenerDAL == nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Internal server error", "ListenerDAL is not initialized")
		logger.Error("listenerDAL is nil")
		return
	}

	// listener, err := h.listenerDAL.GetListenerById(ctx.Request.Context(), payloadRequest.ListenerID)
	// if err != nil {
	// 	models.ResponseError(ctx, http.StatusNotFound, "Listener not found", err.Error())
	// 	logger.Error("Error retrieving the listener:", err)
	// 	return
	// }

	payloadConfig := models.IntoPayloadConfig(payloadRequest)

	// TODO: might have to first marshal here, maybe update the listener config into json.RawMessage?
	// payloadConfig.ListenerConfig = listener.RawConfig

	privateKey, publicKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to generate server ECDH keys", err.Error())
		logger.Error("Error generating ECDH keys:", err)
		return
	}

	payloadConfig.PublicKey = publicKey
	payloadConfig.PrivateKey = privateKey
	payloadConfig.Token = uuid.New().String()

	file, err := json.MarshalIndent(payloadConfig, "", "  ")
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to serialize payload config", err.Error())
		logger.Error("Error marshalling payload config to JSON:", err)
		return
	}

	baseconfPath := conf.GetBaseConfPath()
	if err = os.WriteFile(baseconfPath, file, 0644); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to write configuration file", err.Error())
		logger.Error("Error writing JSON file:", err)
		return
	}

	// args := []string{
	// 	"publish",
	// 	"--configuration", "Release",
	// 	"--self-contained", strings.ToLower(fmt.Sprintf("%t", payloadRequest.SelfContained)),
	// 	"-o", "/app/build/payload-" + payloadConfig.PayloadID,
	// 	"-p:PublishSingleFile=true",
	// 	// "-p:DefineConstants=TYPE_" + listener.Type,
	// 	"-r", payloadConfig.Arch,
	// 	conf.GetAgentProjectFilepath(),
	// }

	// cmd := exec.Command("dotnet", args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err := cmd.Run(); err != nil {
	// 	models.ResponseError(ctx, http.StatusInternalServerError, "Failed to compile payload", err.Error())
	// 	logger.Error("Error running Docker container to compile agent:", err)
	// 	return
	// }

	if err := h.payloadDAL.CreatePayload(ctx.Request.Context(), payloadConfig); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to save payload configuration", err.Error())
		logger.Error("Error saving payload configuration:", err)
		return
	}

	// TODO: remove, to create a payload, an agent_config ID must be passed in
	// No longer creating configs alongside the payload
	// agentConfig := models.IntoAgentConfig(payloadConfig)
	// if err := h.agentDAL.CreateAgentConfig(ctx.Request.Context(), agentConfig); err != nil {
	// 	models.ResponseError(ctx, http.StatusInternalServerError, "Failed to save agent configuration", err.Error())
	// 	logger.Error("Error saving agent configuration:", err)
	// 	return
	// }

	defer func() {
		if err := os.Truncate(baseconfPath, 0); err != nil {
			logger.Error("Error cleaning baseconf.json:", err)
		}
	}()

	models.ResponseSuccess(ctx, http.StatusCreated, "Payload created successfully", payloadConfig)
}

func (h *PayloadController) DeletePayload(ctx *gin.Context) {
	payloadId := ctx.Param(models.ParamPayloadID)
	if payloadId == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamPayloadID))
		return
	}

	filePath := "./teamserver/build/payload-" + payloadId
	logger.Info(filePath)

	if err := h.payloadDAL.DeletePayload(ctx.Request.Context(), payloadId); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete payload", err.Error())
		logger.Error("Error deleting payload:", err)
		return
	}

	if err := os.RemoveAll(filePath); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete payload files", err.Error())
		logger.Error("Error deleting payload folder:", err)
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Payload deleted successfully", nil)
}

func (h *PayloadController) DeleteAllPayloads(ctx *gin.Context) {
	dirPath := "./teamserver/build"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to read payload directory", err.Error())
		logger.Error("Error reading payload directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "payload-") {
			filePath := filepath.Join(dirPath, file.Name())
			if err := os.RemoveAll(filePath); err != nil {
				logger.Error("Error deleting payload directory:", err)
				continue
			}
		}
	}

	if err := h.payloadDAL.DeleteAllPayloads(ctx.Request.Context()); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete payloads from database", err.Error())
		logger.Error("Error deleting payloads:", err)
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "All payloads deleted successfully", nil)
}

func (h *PayloadController) DownloadPayload(ctx *gin.Context) {
	payloadId := ctx.Param(models.ParamPayloadID)
	if payloadId == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamPayloadID))
		return
	}

	extensions := []string{".exe", ".bin", ".dll", ""}
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
		models.ResponseError(ctx, http.StatusNotFound, "Executable not found", "No matching payload file found")
		logger.Error("Executable not found for payload:", payloadId)
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(executablePath)))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(executablePath)
}

func (h *PayloadController) GetAllPayloads(ctx *gin.Context) {
	payloads, err := h.payloadDAL.GetAllPayloads(ctx.Request.Context())
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get payloads", err.Error())
		logger.Error("Error getting payloads:", err)
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Payloads retrieved successfully", payloads)
}

func (h *PayloadController) GetPayloadByToken(ctx *gin.Context) {
	authToken := ctx.Param(models.ParamPayloadToken)

	payload, err := h.payloadDAL.GetPayloadByToken(ctx, authToken)
	if err != nil {
		logger.Error("Error getting payload payload by token:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get payload by token", err.Error())
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Payload retrieved successfully", payload)
}

// Unexported for users, internal use only
func (h *PayloadController) GetToken(ctx *gin.Context) {
	payloadID := ctx.Param(models.ParamPayloadID)

	token, err := h.payloadDAL.GetToken(ctx, payloadID)
	if err != nil {
		logger.Error("Error getting payload token:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get payload token", err.Error())
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Payload token retrieved successfully", token)
}

// Unexported for users, internal use only
func (h *PayloadController) GetKeys(ctx *gin.Context) {
	authToken := ctx.Param(models.ParamPayloadToken)

	privKey, pubKey, err := h.payloadDAL.GetKeys(ctx, authToken)
	if err != nil {
		logger.Error("Error getting payload keys:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get payload keys", err.Error())
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Payload keys retrieved successfully", storage.KeyPair{
		PublicKey:  pubKey,
		PrivateKey: privKey,
	})
}
