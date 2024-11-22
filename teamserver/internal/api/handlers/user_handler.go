package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/models"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/utils"
	"net/http"
)

type UserController struct {
	dal *dal.UserDAL
}

func NewUserController(dal *dal.UserDAL) *UserController {
	return &UserController{dal: dal}
}

func (uc *UserController) GetUsers(ctx *gin.Context) {

	getUsers, err := uc.dal.GetUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error retrieving getUsers: %s", err.Error())})
		return
	}

	ctx.JSON(http.StatusOK, getUsers)
}

func (uc *UserController) AddUsers(ctx *gin.Context) {

	var user models.ResUser

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSONP(http.StatusConflict, gin.H{
			"Error":   err.Error(),
			"Message": "Invalid Request Error.",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	hashPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Invalid Request Error.",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}
	user.PasswordHash = hashPassword

	validate := utils.NewValidatorService()
	if err := validate.ValidateStruct(user); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Format Validation Error.",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}
	err = uc.dal.AddUsers(ctx.Request.Context(), &user)
	if err != nil {
		ctx.JSONP(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"Message": "Error Adding Users.",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	ctx.JSONP(http.StatusCreated, gin.H{
		"Message": "User Created Successfully.",
		"Status":  "Success",
	})
}

func (uc *UserController) GetUsersController(ctx *gin.Context) {
	users, err := uc.dal.GetUsers(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Error retrieving users: %s", err.Error()),
		})
		ctx.Abort()
		return
	}

	ctx.JSONP(http.StatusOK, gin.H{
		"Data":    users,
		"Status":  "Success",
		"Message": "Users Fetched Successfully",
	})
}
