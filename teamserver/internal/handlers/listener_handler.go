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
	var listener models.Listener
	if err := c.ShouldBindJSON(&listener); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.listenerDAL.CreateListener(c.Request.Context(), &listener); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, listener)
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
