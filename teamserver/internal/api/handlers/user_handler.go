package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/app/users"
	"net/http"
)

type UserController struct {
	service *users.Service
}

func NewUserController(service *users.Service) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) GetUsers(ctx *gin.Context) {

	getUsers, err := uc.service.GetUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error retrieving getUsers: %s", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, getUsers)
}
