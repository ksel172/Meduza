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
	dal  dal.IUserDAL
	jwtS models.JWTServiceProvider
}

func NewAuthController(dal dal.IUserDAL, jwtS models.JWTServiceProvider) *AuthController {
	return &AuthController{
		dal:  dal,
		jwtS: jwtS,
	}
}

func (ac *AuthController) LoginController(ctx *gin.Context) {
	var loginR models.AuthRequest

	if err := ctx.ShouldBindJSON(&loginR); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	user, err := ac.dal.GetUserByUsername(ctx.Request.Context(), loginR.Username)
	if err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to authenticate user", "Invalid credentials")
		return
	}

	if !utils.CheckPasswordHash(loginR.Password, user.PasswordHash) {
		log.Print("failed password check: ", loginR.Password, user.PasswordHash)
		models.ResponseError(ctx, http.StatusUnauthorized, "Failed to authenticate user", "Invalid credentials")
		return
	}

	tokens, err := ac.jwtS.GenerateTokens(user.ID, user.Role)
	log.Printf("token errors: %v", err)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to generate authentication tokens", err.Error())
		return
	}

	ctx.SetCookie("refresh_token",
		tokens.RefreshToken,
		refresh_token_Age,
		cookie_path,
		cookie_domain,
		refresh_secure,
		refresh_http)

	models.ResponseSuccess(ctx, http.StatusOK, "User authenticated successfully", tokens)
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
			log.Printf("Refresh token invalid or already expired: %v", err)
		}
	}

	ctx.SetCookie("refresh_token",
		"",
		-1,
		cookie_path,
		cookie_domain,
		refresh_secure,
		refresh_http)

	models.ResponseSuccess(ctx, http.StatusOK, "User logged out successfully", nil)
}

func (ac *AuthController) RefreshTokenController(ctx *gin.Context) {
	prevRefreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		models.ResponseError(ctx, http.StatusUnauthorized, "Failed to refresh token", "No refresh token found")
		return
	}
	if prevRefreshToken == "" {
		models.ResponseError(ctx, http.StatusUnauthorized, "Failed to refresh token", "Empty refresh token")
		return
	}

	claims, err := ac.jwtS.ValidateToken(prevRefreshToken)
	if err != nil {
		models.ResponseError(ctx, http.StatusUnauthorized, "Failed to validate token", err.Error())
		return
	}

	if claims.ExpiresAt.Before(time.Now()) {
		models.ResponseError(ctx, http.StatusUnauthorized, "Failed to refresh token", "Refresh token has expired")
		return
	}

	tokens, err := ac.jwtS.GenerateTokens(claims.ID, claims.Role)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to generate tokens", err.Error())
		return
	}

	ctx.SetCookie("refresh_token",
		tokens.RefreshToken,
		refresh_token_Age,
		cookie_path,
		cookie_domain,
		refresh_secure,
		refresh_http)

	models.ResponseSuccess(ctx, http.StatusOK, "Token refreshed successfully", gin.H{
		"access_token": tokens.Token,
	})
}
