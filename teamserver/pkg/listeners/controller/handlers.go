package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	// err error `json:"-"`
}

// This is good for debug, for production definitely not.
// We will have to handle the possible error kinds and create a new error
// with the message to show the user at the handler level given the error that occurred
func ErrorResponse(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, Response{
		Status:  status,
		Message: err.Error(),
	})
}

func SuccessResponse(ctx *gin.Context, status int, message string, data any) {
	ctx.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// Add will setup a new listener in the controller but not start it
func (c *Controller) AddListener(ctx *gin.Context) {
	var listenerConfig ListenerConfig
	if err := ctx.ShouldBindJSON(&listenerConfig); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if err := c.manager.addListener(listenerConfig); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	SuccessResponse(ctx, http.StatusCreated, fmt.Sprintf("listener %s created", listenerConfig.ID), nil)
}

// Start will start a listener, given it has the ListenerReady status
func (c *Controller) StartListener(ctx *gin.Context) {
	listenerID, ok := ctx.Params.Get("listenerID")
	if !ok {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing listener ID in query params"))
		return
	}

	errChan := make(chan error)
	if err := c.manager.startListener(ctx.Request.Context(), listenerID, errChan); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	select {
	case err, ok := <-errChan:
		if !ok { // Channel closed, meaning successful completion
			SuccessResponse(ctx, http.StatusOK, fmt.Sprintf("listener %s started", listenerID), nil)
			return
		}
		log.Printf("Failed to start listener: %v", err)
		ErrorResponse(ctx, http.StatusInternalServerError, errors.New("listener failed to start"))
		return
	case <-ctx.Request.Context().Done():
		log.Printf("Listener start timed out")
		ErrorResponse(ctx, http.StatusRequestTimeout, errors.New("listener start timed out"))
		return
	}
}

// Stop will stop a listener, given it has the ListenerRunning status
func (c *Controller) StopListener(ctx *gin.Context) {
	listenerID, ok := ctx.Params.Get("listenerID")
	if !ok {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing listener ID in query params"))
		return
	}

	errChan := make(chan error)
	if err := c.manager.stopListener(ctx.Request.Context(), listenerID, errChan); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	select {
	case err, ok := <-errChan:
		if !ok { // Channel closed, meaning successful completion
			SuccessResponse(ctx, http.StatusOK, fmt.Sprintf("listener %s stopped", listenerID), nil)
			return
		}
		log.Printf("Failed to stop listener: %v", err)
		ErrorResponse(ctx, http.StatusInternalServerError, errors.New("listener failed to stop"))
		return
	case <-ctx.Request.Context().Done():
		log.Printf("Listener start timed out")
		ErrorResponse(ctx, http.StatusRequestTimeout, errors.New("listener stop timed out"))
		return
	}
}

// Terminate will fully remove a listener from the controller's jurisdiction
// if it is running, it will first stop it
// If it fails to be stopped, it will not remove the listener from the controller
func (c *Controller) TerminateListener(ctx *gin.Context) {
	listenerID, ok := ctx.Params.Get("listenerID")
	if !ok {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing listener ID in query params"))
		return
	}

	if err := c.manager.terminateListener(ctx.Request.Context(), listenerID); err != nil {
		log.Printf("failed to terminate listener '%s': %v", listenerID, err)
		ErrorResponse(ctx, http.StatusInternalServerError, errors.New("failed to terminate listener"))
		return
	}

	SuccessResponse(ctx, http.StatusOK, fmt.Sprintf("listener %s terminated", listenerID), nil)
}

// Updates a listener config, not all fields can be update
// Some will require the listener to be deleted and re-added
func (c *Controller) UpdateListener(ctx *gin.Context) {
	var listenerConfig ListenerConfig
	if err := ctx.ShouldBindJSON(&listenerConfig); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if err := c.manager.updateListener(ctx.Request.Context(), listenerConfig); err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	SuccessResponse(ctx, http.StatusCreated, fmt.Sprintf("listener %s created", listenerConfig.ID), nil)
}

// External listeners call this endpoint to update the status of the listener
// Something like: <HOST>/api/v1/
func (c *Controller) UpdateListenerStatus(ctx *gin.Context) {

	var statusUpdateRequest struct {
		ListenerID string `json:"listener_id"`
		Status     string `json:"status" validate:"oneof=pending ready starting running stopping terminating"`
	}

	if err := ctx.ShouldBindJSON(&statusUpdateRequest); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	c.manager.updateListenerStatus(ctx.Request.Context(),
		statusUpdateRequest.ListenerID, statusUpdateRequest.Status)

	ctx.Status(http.StatusOK)
}

// External listeners run a heartbeat + configuration synchronization loop
// works to know what listeners are still alive and to also re-send
func (c *Controller) SynchronizeConfig(ctx *gin.Context) {
	listenerID, ok := ctx.Params.Get("listenerID")
	if !ok {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing listener ID in query params"))
		return
	}

	config, err := c.manager.synchronize(listenerID)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, config)
}
