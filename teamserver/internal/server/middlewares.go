package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

func (s *Server) JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := s.dependencies.JwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func (s *Server) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Authorization header Error",
				"error":   "Authorization is empty",
				"Status":  "Empty",
			})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Invalid authorization format",
				"error":   "It should be a bearer token",
				"Status":  "Failed",
			})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := s.dependencies.JwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Invalid claims",
				"Error":   err.Error(),
				"Status":  "Failed",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)

		if claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Restricted Route. Only Admins are allowed.",
				"Status":  "Failed",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *Server) ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Authorization header Error",
				"error":   "Authorization is empty",
				"Status":  "Empty",
			})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Invalid authorization format",
				"error":   "It should be a bearer token",
				"Status":  "Failed",
			})
			c.Abort()
			return
		}

		tokenString := bearerToken[1]
		claims, err := s.dependencies.JwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Invalid claims",
				"Error":   err.Error(),
				"Status":  "Failed",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)

		if claims.Role != "moderator" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Message": "Restricted Route. Only Admins are allowed.",
				"Status":  "Failed",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
