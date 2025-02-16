package models

import (
	"github.com/gin-gonic/gin"
)

type ServerResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success sends a successful response
func ResponseSuccess(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, ServerResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func ResponseError(ctx *gin.Context, status int, message string, err interface{}) {
	ctx.JSON(status, ServerResponse{
		Status:  status,
		Message: message,
		Error:   err,
	})
}
