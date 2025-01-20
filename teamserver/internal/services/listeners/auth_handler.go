package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

/*
	Agent authentication is used so agents can obtain the AES key for encrypted communication
		1. The server uses an UUID token when compiling the payload, this UUID is stored within the compiled payload itself.
		2. On the first agent check-in, the agent will send this token to the server in the Authorization header, the server will verify it against its payload table in the db.
		3. Server will generate an AES key to be used for future communications with this agent and sends the key to the agent, who will store it in memory only.
		4. Agent is verified and communications (message field in the C2 request model) are now encrypted
*/

const AuthToken = "token"

type IAgentAuthController interface {
	AuthenticateAgent(*gin.Context)
}

type AgentAuthController struct {
	payloadDAL dal.IPayloadDAL
}

func NewAgentAuthController(payloadDAL dal.IPayloadDAL) *AgentAuthController {
	return &AgentAuthController{
		payloadDAL: payloadDAL,
	}
}

func (a *AgentAuthController) AuthenticateAgent(ctx *gin.Context) {
	payloadToken := ctx.GetHeader("Authorization")
	if payloadToken == "" {
		logger.Info("Skipping agent authentication")
		ctx.Next() // move on to the next handler
	}

	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	logger.Info(LogLevel, LogDetail, fmt.Sprintf("Authenticating agent %s", c2request.AgentID))
	storedPayloadToken, err := a.payloadDAL.GetPayloadToken(ctx.Request.Context(), c2request.ConfigID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed authentication: invalid config %s", c2request.ConfigID)
		logger.Error(LogLevel, LogDetail, errMsg)
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	if payloadToken != storedPayloadToken {
		errMsg := fmt.Sprintf("Failed authentication: provided agent token '%s' by agent %s does not match the token for config %s", payloadToken, c2request.AgentID, c2request.ConfigID)
		logger.Error(LogLevel, LogDetail, errMsg)
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
		return
	}

	ctx.Set(AuthToken, payloadToken)
}
