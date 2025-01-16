package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

/*
	Agent authentication is used so agents can obtain the AES key for encrypted communication
		1. The server uses the Payload UUID to compile the payload, this UUID is stored within the compiled payload itself.
		2. On the first agent check-in, the agent will send the UUID of this payload and the server will verify it against its payload table in the db.
		3. Server will generate an AES key to be used for future communications with this agent and sends the key to the agent, who will store it in memory only.
		4. Agent is verified and communications are now encrypted
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
		ctx.Next() // move on to the next handler
	}

	var c2request models.C2Request
	if err := ctx.ShouldBindJSON(&c2request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedPayloadToken, err := a.payloadDAL.GetPayloadToken(ctx.Request.Context(), c2request.ConfigID)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	if payloadToken != storedPayloadToken {
		ctx.Status(http.StatusUnauthorized)
		ctx.Abort()
	}
	ctx.Set(AuthToken, payloadToken)
}
