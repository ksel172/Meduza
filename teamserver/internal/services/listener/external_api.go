package listener

// Move this to the handlers and manage routes with middlewares?

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
)

// Entrypoint for external listener requests

// All external listeners call back to this server
// The server is responsible for receiving listener C2 requests
type ExternalServer struct {
	host     string
	port     int
	server   *gin.Engine
	registry *ListenerRegistry
}

func NewExternalServer() *ExternalServer {
	return &ExternalServer{}
}

func (es *ExternalServer) RegisterListener(ctx *gin.Context) error {
	var listenerModel models.Listener

	if err := ctx.ShouldBindJSON(&listenerModel); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Invalid request body", err.Error())
		return nil
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener controller registered successfully", nil)
	return nil
}

func (es *ExternalServer) CheckInExternalListener(ctx *gin.Context) error {
	return nil
}
