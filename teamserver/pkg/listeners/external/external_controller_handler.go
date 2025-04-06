package external_listener

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
)

type ExternalControllerHandler struct {
}

func NewExternalListenerHandler() *ExternalControllerHandler {
	return &ExternalControllerHandler{}
}

func (ExternalControllerHandler) RegisterExternalListener(ctx *gin.Context) error {
	var listenerController models.ListenerController

	if err := ctx.ShouldBindJSON(&listenerController); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Invalid request body", err.Error())
		return nil
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener controller registered successfully", nil)
	return nil
}

func (ExternalControllerHandler) CheckInExternalListener(ctx *gin.Context) error {

	// Send listener related updates

	return nil
}
