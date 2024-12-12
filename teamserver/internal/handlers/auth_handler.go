package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"

	"github.com/gin-gonic/gin"
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
	jwtS *models.JWTService
}

func NewAuthController(dal *dal.UserDAL, jwtS *models.JWTService) *AuthController {
	return &AuthController{
		dal:  dal,
		jwtS: jwtS,
	}
}

func (ac *AuthController) LoginController(ctx *gin.Context) {
	var loginR models.AuthRequest

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
func (ac *AuthController) LogoutController(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	var accessToken string
	if len(header) > 7 && header[:7] == "Bearer " {
		accessToken = header[7:]
	}

	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		refreshToken = ""
	}

	if accessToken != "" {
		claims, err := ac.jwtS.ValidateToken(accessToken)
		if err == nil {
			expiry := time.Unix(claims.ExpiresAt.Unix(), 0)
			ac.jwtS.RevokeToken(accessToken, expiry)
		} else {
			log.Printf("Access token invalid or already expired: %v", err)
		}
	}

	if refreshToken != "" {
		claims, err := ac.jwtS.ValidateToken(refreshToken)
		if err == nil {
			expiry := time.Unix(claims.ExpiresAt.Unix(), 0)
			ac.jwtS.RevokeToken(refreshToken, expiry)
		} else {
			log.Printf("Access token invalid or already expired: %v", err)
		}
	}

	ctx.SetCookie("refresh_token",
		"",
		-1,
		cookie_path,
		cookie_domain,
		refresh_secure,
		refresh_http)

	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Successfully logged Out",
		"Status":  "Sucess",
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
