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
	users, err := uc.dal.GetUsers(ctx.Request.Context())
	if err != nil {
		logger.Error("Error retrieving users:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Users retrieved successfully", users)
}

func (uc *UserController) AddUsers(ctx *gin.Context) {
	var user models.ResUser

	if err := ctx.ShouldBindJSON(&user); err != nil {
		logger.Error("Invalid request body:", err)
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	hashPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		logger.Error("Failed to hash password:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}
	user.PasswordHash = hashPassword

	validate := utils.NewValidatorService()
	if err := validate.ValidateStruct(user); err != nil {
		logger.Error("Validation error:", err)
		models.ResponseError(ctx, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	if err = uc.dal.AddUsers(ctx.Request.Context(), &user); err != nil {
		logger.Error("Failed to add user:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to add user", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "User created successfully", user)
}

func (uc *UserController) GetUsersController(ctx *gin.Context) {
	users, err := uc.dal.GetUsers(ctx.Request.Context())
	if err != nil {
		logger.Error("Error retrieving users:", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Users retrieved successfully", users)
}
