package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type UserController struct {
	dal dal.IUserDAL
}

func NewUserController(dal dal.IUserDAL) *UserController {
	return &UserController{dal: dal}
}

func (uc *UserController) GetUsers(ctx *gin.Context) {

	getUsers, err := uc.dal.GetUsers(ctx.Request.Context())
	if err != nil {
		logger.Error("Error retrieving Users from GetUsers Function :", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error retrieving users",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":   getUsers,
		"status": "success",
	})
}

func (uc *UserController) AddUsers(ctx *gin.Context) {

	var user models.ResUser

	if err := ctx.ShouldBindJSON(&user); err != nil {
		logger.Error("Invalid Request error From AddUsers function: ", err)
		ctx.JSONP(http.StatusConflict, gin.H{
			"message": "invalid request error.",
			"status":  "failed",
		})
		ctx.Abort()
		return
	}

	hashPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		logger.Error("Invalid Request Error From hashing: ", err)
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"message": "invalid request error.",
			"status":  "error",
		})
		ctx.Abort()
		return
	}
	user.PasswordHash = hashPassword

	validate := utils.NewValidatorService()
	if err := validate.ValidateStruct(user); err != nil {
		logger.Error("Format Validation Error while creating/adding users:", err)
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"message": "format validation error.",
			"status":  "failed",
		})
		ctx.Abort()
		return
	}
	err = uc.dal.AddUsers(ctx.Request.Context(), &user)
	if err != nil {
		logger.Error("Error Adding Users :", err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{
			"message": "error adding users.",
			"status":  "failed",
		})
		ctx.Abort()
		return
	}

	ctx.JSONP(http.StatusCreated, gin.H{
		"message": "user created successfully.",
		"status":  "success",
	})
}

func (uc *UserController) GetUsersController(ctx *gin.Context) {
	users, err := uc.dal.GetUsers(ctx.Request.Context())
	if err != nil {
		logger.Error("Error Retrieveing Users: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error retrieving users",
		})
		ctx.Abort()
		return
	}

	ctx.JSONP(http.StatusOK, gin.H{
		"data":    users,
		"status":  "success",
		"message": "Users Fetched Successfully",
	})
}
