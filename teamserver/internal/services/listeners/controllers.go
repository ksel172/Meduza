package services

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
	http_listener "github.com/ksel172/Meduza/teamserver/pkg/listeners/http"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

/*
	This file should hold all server initializations for any listener type
	If server is initialized in the listeners package, there will be a circular import
*/

var listenerRoute = "/" // Changed to "/" from "/check-in"
var status = utils.Status

// NewHTTPListenerController initializes a new HTTP listener controller.
// The controller is responsible for handling HTTP requests to the listener.
func NewHTTPListenerController(
	name string,
	config models.HTTPListenerConfig,
	checkInController ICheckInController,
) (*http_listener.HTTPListenerController, error) {

	if err := config.Validate(); err != nil {
		return nil, err
	}

	mux := gin.Default()

	mux.Use(func(ctx *gin.Context) {
		if config.WhitelistEnabled {
			clientIP := ctx.ClientIP()
			if !isIpWhitelisted(clientIP, config.Whitelist) {
				logger.Warn("IP not allowed:", clientIP)
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "IP not allowed",
					"status":  status.ERROR,
				})
				return
			}
		}
		ctx.Next()
	})

	// Handle the listener routes
	mux.Handle("GET", listenerRoute, checkInController.GetTasks)
	mux.Handle("POST", listenerRoute, checkInController.CreateAgent)

	return &http_listener.HTTPListenerController{
		Name:      name,
		Config:    config,
		Server:    mux,
		HTTPServe: nil,
	}, nil
}

// isIpWhitelisted checks if an IP is in the whitelist.
func isIpWhitelisted(ip string, whitelist []string) bool {
	for _, allowedIP := range whitelist {
		if strings.TrimSpace(allowedIP) == ip {
			return true
		}
	}
	return false
}
