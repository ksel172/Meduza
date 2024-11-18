package api

import (
	/* "encoding/json"
	"fmt"
	"log" */
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/services/auth"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type AuthController struct {
	dal  *storage.UserDAL
	jwtS *auth.JWTService
}

func NewAuthController(dal *storage.UserDAL, jwtS *auth.JWTService) *AuthController {
	return &AuthController{
		dal:  dal,
		jwtS: jwtS,
	}
}

func (ac *AuthController) LoginController(ctx *gin.Context) {
	var loginR auth.AuthRequest

	if err := ctx.ShouldBindJSON(&loginR); err != nil {
		ctx.JSONP(http.StatusConflict, gin.H{
			"Error":   err.Error(),
			"Message": "Invalid Request",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	user, err := ac.dal.GetUserByUsername(ctx.Request.Context(), loginR.Username)
	if err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Invalid Credentials or Request Error.",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	if !utils.CheckPasswordHash(loginR.Password, user.PasswordHash) {
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"Error":   "Invalid Credentials Error",
			"Message": "Invalid Credentials or Request Error.",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}
	tokens, err := ac.jwtS.GenerateTokens(user.ID, user.Role)
	if err != nil {
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"Error":   err.Error(),
			"Message": "Cannot Generate Token",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	ctx.JSONP(http.StatusOK, gin.H{
		"Key":     tokens,
		"Message": "User Authenticated Successfully",
		"Status":  "Success",
	})
}
