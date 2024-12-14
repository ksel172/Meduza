package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	u "github.com/ksel172/Meduza/teamserver/utils"
)

type ListenersHandler struct {
	dal *dal.ListenersDAL
}

func NewListenersHandler(dal *dal.ListenersDAL) *ListenersHandler {
	return &ListenersHandler{dal: dal}
}

func (h *ListenersHandler) AddListener(ctx *gin.Context) {
	var listener models.Listener

	if err := ctx.ShouldBindJSON(&listener); err != nil {
		logger.Error("Request Body Error while bind the json:", err)
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "Invalid Request body.Please type correct input",
			"status":  "error",
		})
		ctx.Abort()
		return
	}
	listener.Whitelist = u.NormalizeToSlice(listener.Whitelist)
	listener.Blacklist = u.NormalizeToSlice(listener.Blacklist)

	err := h.dal.CreateListeners(ctx.Request.Context(), &listener)
	if err != nil {
		logger.Error("Error Occured while Adding Data to listener:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Unable to create a listener.",
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "listener created successfully",
	})

}
