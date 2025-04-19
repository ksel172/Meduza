package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	listenerService "github.com/ksel172/Meduza/teamserver/internal/services/listener"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

type ListenerController struct {
	service *listenerService.ListenerService
	// listenerClient *external_listener.ListenerClient
	listenerDal dal.IListenerDAL
	// controllerDal  dal.IControllerDal
}

func NewListenersHandler(service *listenerService.ListenerService, listenerDAL dal.IListenerDAL) *ListenerController {
	return &ListenerController{
		service: service,
		// listenerClient: listenerClient,
		listenerDal: listenerDAL,
		// controllerDal:  controllerDAL,
	}
}

func (lc *ListenerController) GetAllListeners(ctx *gin.Context) {
	listeners, err := lc.listenerDal.GetAllListeners(ctx.Request.Context())
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error getting listeners", err.Error())
		return
	}
	models.ResponseSuccess(ctx, http.StatusOK, "Listeners retrieved successfully", listeners)
}

func (lc *ListenerController) GetListener(ctx *gin.Context) {
	listenerID := ctx.Param(models.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	listener, err := lc.listenerDal.GetListenerById(ctx.Request.Context(), listenerID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error getting listener", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener retrieved successfully", listener)
}

func (lc *ListenerController) GetListenerStatuses(ctx *gin.Context) {
	listeners, err := lc.listenerDal.GetAllListeners(ctx.Request.Context())
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error getting listeners", err.Error())
		return
	}

	listenerStatuses := make(map[string]string)
	for _, listener := range listeners {
		listenerStatuses[listener.ID] = listener.Status
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listeners retrieved succesfully", listenerStatuses)
}

func (lc *ListenerController) CreateListener(ctx *gin.Context) {
	var listenerModel models.Listener
	if err := ctx.ShouldBindJSON(&listenerModel); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to get listener from request", err.Error())
		return
	}

	// TODO: external listeners registration

	// Perhaps we will need to look into the registration process. The fact that it is created from
	// the handler does not mean that it isn't external. CreateListener only creates the config for the
	// listener. If external listeners are registered this is still relevant.

	// Either we use a different API when handling external listeners (since we need to fill params
	// dynamically anyways) or we use some other way to seperate concerns.
	if listenerModel.IsExternal {
		models.ResponseError(ctx, http.StatusBadRequest, "external listeners should register themselves", nil)
		return
	}

	listenerModel.ID = uuid.NewString()

	if err := lc.listenerDal.CreateListener(ctx.Request.Context(), &listenerModel); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create listener", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "Listener created", nil)
}

func (lc *ListenerController) StartListener(ctx *gin.Context) {
	listenerID := ctx.Param(models.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	if err := lc.service.StartListener(ctx.Request.Context(), listenerID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error starting listener", err.Error())
		return
	}
	models.ResponseSuccess(ctx, http.StatusOK, fmt.Sprintf("Successfully started listener with ID: %s", listenerID), nil)
}

func (lc *ListenerController) StopListener(ctx *gin.Context) {
	listenerID := ctx.Param(models.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	if err := lc.service.StopListener(ctx.Request.Context(), listenerID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error stopping listener", err.Error())
		return
	}
	models.ResponseSuccess(ctx, http.StatusOK, fmt.Sprintf("Successfully stopped listener with ID: %s", listenerID), nil)
}

func (lc *ListenerController) TerminateListener(ctx *gin.Context) {
	listenerID := ctx.Param(models.ParamListenerID)
	if listenerID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid listener ID", "Listener ID is required")
		return
	}

	err := lc.service.TerminateListener(ctx, listenerID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error deleting listener", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Listener deleted successfully", nil)
}

func (lc *ListenerController) UpdateListener(ctx *gin.Context) {
	var requestListener models.Listener
	if err := ctx.ShouldBindJSON(&requestListener); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	listener, err := lc.listenerDal.GetListenerByName(ctx, requestListener.Name)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "listener not found", err.Error())
		return
	}

	// Update the listener
	if err = lc.service.UpdateListener(ctx, listener); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError,
			fmt.Sprintf("Failed to update listener with ID: %s", listener.ID), err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, fmt.Sprintf("Successfully updated listener %s", listener.ID), nil)
}

func (lc *ListenerController) AutoStart(ctx context.Context) error {
	err := lc.service.AutoStart(ctx)
	if err != nil {
		return fmt.Errorf("error starting listeners: %v", err)
	}
	return nil
}

func (lc *ListenerController) GetListenerConfigurations(ctx *gin.Context) {

}
