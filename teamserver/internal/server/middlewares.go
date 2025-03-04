package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var JWTSecret string

func init() {
	JWTSecret = utils.GetEnvString("JWT_SECRET", "")
}

func (s *Server) HandleCors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle prefight OPTIONS request
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func (s *Server) UserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			models.ResponseError(ctx, http.StatusUnauthorized, "Authorization header is missing", nil)
			ctx.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			models.ResponseError(ctx, http.StatusUnauthorized, "Invalid authorization format", nil)
			ctx.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := s.dependencies.JwtService.ValidateToken(tokenString)
		if err != nil {
			models.ResponseError(ctx, http.StatusUnauthorized, "Invalid claims", err)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func (s *Server) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			models.ResponseError(ctx, http.StatusUnauthorized, "Authorization header is missing", nil)
			ctx.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			models.ResponseError(ctx, http.StatusUnauthorized, "Invalid authorization format", nil)
			ctx.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := s.dependencies.JwtService.ValidateToken(tokenString)
		if err != nil {
			models.ResponseError(ctx, http.StatusUnauthorized, "Invalid claims", err)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)

		if claims.Role != "admin" {
			models.ResponseError(ctx, http.StatusUnauthorized, "Restricted Route. Only Admins are allowed.", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
