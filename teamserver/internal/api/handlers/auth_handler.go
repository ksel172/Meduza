package handlers

import (
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/services/auth"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var (
	refresh_token_Age = 30 * 24 * 60 * 60
	cookie_path       = utils.GetEnvString("COOKIE_PATH", "")
	cookie_domain     = utils.GetEnvString("COOKIE_DOMAIN", "")
	refresh_secure    = utils.GetEnvBool("REFRESH_SECURE", false)
	refresh_http      = utils.GetEnvBool("REFRESH_HTTP", false)
)

type AuthController struct {
	dal  *dal.UserDAL
	jwtS *auth.JWTService
}

func NewAuthController(dal *dal.UserDAL, jwtS *auth.JWTService) *AuthController {
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
			"Message": "Invalid Credentials or Request Error",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	if !utils.CheckPasswordHash(loginR.Password, user.PasswordHash) {
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"Message": "Invalid Credentials or Request Error",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}
	tokens, err := ac.jwtS.GenerateTokens(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"Details": "Token generation failed",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	ctx.SetCookie("refresh_token",
		tokens.RefreshToken,
		refresh_token_Age,
		cookie_path,
		cookie_domain,
		refresh_secure,
		refresh_http)

	ctx.JSONP(http.StatusOK, gin.H{
		"Key":     tokens,
		"Message": "User Authenticated Successfully",
		"Status":  "Success",
	})
}

func (ac *AuthController) RefreshTokenController(ctx *gin.Context) {
	prevRefreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"Error":   err.Error(),
			"Message": "Refresh Token Error",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}
	if prevRefreshToken == "" {
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"Message": "Empty Refresh token Cookie",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	claims, err := ac.jwtS.ValidateToken(prevRefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"Error":   err.Error(),
			"Message": "Token validation failed",
			"Status":  "Failed",
		})
		return
	}

	if claims.ExpiresAt.Before(time.Now()) {
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"Message": "Refresh token has expired",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	tokens, err := ac.jwtS.GenerateTokens(claims.ID, claims.Role)
	if err != nil {
		ctx.JSONP(http.StatusInternalServerError, gin.H{
			"Message": "Failed to generate tokens",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}
	ctx.SetCookie("refresh_token",
		tokens.RefreshToken,
		refresh_token_Age,
		cookie_path,
		cookie_domain,
		refresh_secure,
		refresh_http)

	ctx.JSONP(http.StatusOK, gin.H{
		"access_token": tokens.Token,
	})
}
