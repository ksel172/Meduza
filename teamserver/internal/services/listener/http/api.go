package http_listener

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// Handle agent check in endpoint
func (l *HTTPListener) HandleCheckIn(ctx *gin.Context) {
	// Read request body to unmarshal into c2request
	// body might be encrypted with AES or not
	// depending on whether the agent is already authenticated)
	body, _ := io.ReadAll(ctx.Request.Body)
	if len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing request body"})
		return
	}
	var c2request models.C2Request

	// If a session token is sent, that means the agent is authenticated
	// Otherwise, treat is as an authentication request
	sessionTokenBase64 := ctx.GetHeader("Session-Token")
	if sessionTokenBase64 == "" {
		// Get base64 Auth-Token header and decode it
		authTokenBase64 := ctx.GetHeader("Auth-Token")
		if authTokenBase64 == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing Auth-Token header"})
			return
		}
		authToken, err := base64.StdEncoding.DecodeString(authTokenBase64)
		if err != nil {
			logger.Info(fmt.Sprintf("failed to base64 decode Auth-Token header: %v", authToken))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Auth-Token header"})
			return
		}

		// The C2Request should not be encrypted at this point, only base64 encoded, unmarshal and use it to authenticate
		// The c2 request should contain only a single message field with the agent public key:
		// BASE 64 ENCODED REQUEST BODY:
		// {"message": <AGENT_PUBLIC_KEY>}
		decodedC2Request, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			logger.Info(fmt.Sprintf("failed to decode c2request: %v, data: %v", err, decodedC2Request))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid encrypted data"})
			return
		}
		var c2request models.C2Request
		if err := json.Unmarshal(decodedC2Request, &c2request); err != nil {
			logger.Info(fmt.Sprintf("Invalid request body: %v", err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Get the agent public key from the request message
		agentPublicKeyBase64 := c2request.Message
		if agentPublicKeyBase64 == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing agent public key"})
			return
		}

		// Decode the base64-encoded public key
		agentPublicKey, err := base64.StdEncoding.DecodeString(agentPublicKeyBase64)
		if err != nil {
			logger.Info(fmt.Sprintf("failed to decode agent public key: %v", err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent public key encoding"})
			return
		}

		// Authenticate agent
		logger.Info(fmt.Sprintf("Handling authentication request for agent %s", c2request.AgentID))
		response, err := l.checkinController.Authenticate(agentPublicKey, string(authToken))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			logger.Info(fmt.Sprintf("failed to authenticate: %v", err))
			return
		}
		logger.Info(fmt.Sprintf("Agent %s authenticated", c2request.AgentID))
		ctx.JSON(http.StatusAccepted, gin.H{
			"public_key":    base64.StdEncoding.EncodeToString(response.PublicKey),
			"session_token": base64.StdEncoding.EncodeToString([]byte(response.SessionToken)),
		})

		return
	}

	// Get session AES key for agent, decrypt request body and unmarshal
	sessionToken, err := base64.StdEncoding.DecodeString(sessionTokenBase64)
	if err != nil {
		logger.Info(fmt.Sprintf("failed to base64 decode Session-Token header: %v", sessionTokenBase64))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Session-Token header"})
		return
	}
	key, exists := storage.KeyRegistry.GetKey(string(sessionToken))
	if !exists {
		// If agent received unauthorized response
		// it should send another authentication request right after
		logger.Info(fmt.Sprintf("Missing AES key for session: %s", sessionTokenBase64))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session token"})
		return
	}
	decryptedData, err := utils.AesDecrypt(key, body)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to decrypt request body: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "decryption failed"})
		return
	}
	if err := json.Unmarshal(decryptedData, &c2request); err != nil {
		logger.Info(fmt.Sprintf("Invalid request body: %v", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	switch c2request.Reason {

	case models.Task:
		logger.Info(fmt.Sprintf("Handling task request for agent %s", c2request.AgentID))
		taskResponse, err := l.checkinController.HandleTaskRequest(ctx.Request.Context(), c2request, string(sessionToken))
		if err != nil {
			switch err {
			case checkin.ErrDatabase:
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			case checkin.ErrInternalServer:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			case checkin.ErrUnauthorized:
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			default:
				panic("unsupported error returned")
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"message": taskResponse})
		return

	case models.Response:
		logger.Info(fmt.Sprintf("Handling response request for agent %s", c2request.AgentID))
		err := l.checkinController.HandleResponseRequest(ctx.Request.Context(), c2request)
		if err != nil {
			switch err {
			case checkin.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			case checkin.ErrInternalServer:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		ctx.JSON(http.StatusOK, "successfully recorded response")
		return

	case models.Register:
		logger.Info(fmt.Sprintf("Handling register request for agent %s", c2request.AgentID))
		err := l.checkinController.HandleRegisterRequest(ctx.Request.Context(), c2request)
		if err != nil {
			switch err {
			case checkin.ErrInvalidData:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			case checkin.ErrConflict:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			case checkin.ErrInternalServer:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.JSON(http.StatusCreated, nil)
		return
	}
}
