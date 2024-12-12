package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

type ListenerHandler struct {
	listenerDAL *dal.ListenerDAL
}

func NewListenerHandler(listenerService *dal.ListenerDAL) *ListenerHandler {
	return &ListenerHandler{listenerDAL: listenerService}
}

func (h *ListenerHandler) CreateListener(c *gin.Context) {
	var listenerRequest models.ListenerRequest
	if err := c.ShouldBindJSON(&listenerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listener := &models.Listener{
		Type:        listenerRequest.Type,
		Name:        listenerRequest.Name,
		Description: listenerRequest.Description,
		Config:      listenerRequest.Config,
	}

	createdListener, err := h.listenerDAL.CreateListener(c.Request.Context(), listener)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdListener)
}

func (h *ListenerHandler) GetListener(c *gin.Context) {
	id := c.Param("id")
	listener, err := h.listenerDAL.GetListener(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, listener)
}

func (h *ListenerHandler) GetListenerConfig(c *gin.Context) {
	id := c.Param("id")
	config, err := h.listenerDAL.GetListenerConfig(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ListenerHandler) UpdateListener(c *gin.Context) {
	var listener models.Listener
	if err := c.ShouldBindJSON(&listener); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.listenerDAL.UpdateListener(c.Request.Context(), &listener); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, listener)
}

func (h *ListenerHandler) DeleteListener(c *gin.Context) {
	id := c.Param("id")
	if err := h.listenerDAL.DeleteListener(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listener deleted"})
}

func (h *ListenerHandler) StartListener(c *gin.Context) {

	// Listener Start logic here

	id := c.Param("id")
	listener, err := h.listenerDAL.GetListener(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Listener not found": err.Error()})
		return
	}
	if listener.Type == models.ListenerTypeHTTP {
		// HTTP Logic

		c.JSON(http.StatusOK, gin.H{"message": "Http listener successfully started"})
	} else if listener.Type == models.ListenerTypeSMB {
		// SMB Logic

		c.JSON(http.StatusOK, gin.H{"message": "Smb listener successfully started"})
	} else if listener.Type == models.ListenerTypeTCP {
		// TCP Logic

		c.JSON(http.StatusOK, gin.H{"message": "Tcp listener successfully started"})
	} else if listener.Type == models.ListenerTypeForeign {
		// Foreign Logic

		c.JSON(http.StatusOK, gin.H{"message": "Foreign listener successfully started"})
	}
}

func (h *ListenerHandler) StopListener(c *gin.Context) {

	// Listener Stop logic here

	//id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Listener stopped"})
}
