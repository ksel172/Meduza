package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type UserController struct {
	dal *storage.UserDAL
}

func NewUserController(dal *storage.UserDAL) *UserController {
	return &UserController{dal: dal}
}

// AddUsersController - Add Users.
// Handles validation , hashing password and request Error
func (uc *UserController) AddUsersController(ctx *gin.Context) {
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

// GetUsers returns All the users present the database.
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
