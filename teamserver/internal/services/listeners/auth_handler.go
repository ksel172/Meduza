package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
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

	payloadUUID := ctx.GetHeader("Authorization")
	if payloadUUID == "" {
		ctx.Next() // move on to the next handler
	}

	payloadConfigs, err := a.payloadDAL.GetAllPayloads(ctx.Request.Context())
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	for _, pConfig := range payloadConfigs {
		if pConfig.PayloadID == payloadUUID {
			ctx.Set(AuthToken, payloadUUID)
			ctx.Next()
		}
	}

	ctx.Status(http.StatusUnauthorized)
	ctx.Abort()
}
