package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/models"
)

// Entrypoint for external listener requests

// All external listeners call back to this server
// The server is responsible for receiving listener C2 requests
type ExternalServer struct {
	host              string
	port              int
	server            *gin.Engine
	registry          *listener.ListenerRegistry
	checkinController *checkin.CheckInController
}

func NewExternalServer(checkinController *checkin.CheckInController) *ExternalServer {
	return &ExternalServer{
		checkinController: checkinController,
	}
}

func (es *ExternalServer) RegisterListener(ctx *gin.Context) {
	var listenerModel models.Listener

	if err := ctx.ShouldBindJSON(&listenerModel); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Invalid request body", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener controller registered successfully", nil)
}

func (es *ExternalServer) HandleAuthentication(ctx *gin.Context) {

}

func (es *ExternalServer) HandleTaskRequest(ctx *gin.Context) {

}

func (es *ExternalServer) HandleResponseSubmission(ctx *gin.Context) {

}

func (es *ExternalServer) HandleAgentRegistration(ctx *gin.Context) {

}

func (es *ExternalServer) CheckInExternalListener(ctx *gin.Context) {

}
