package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var status = utils.Status

// Config represents the configuration settings for the service, including network, security, and proxy settings.
type Config struct {
	WorkingHours     string        `json:"working_hours"`
	Hosts            []string      `json:"hosts"`
	HostBind         string        `json:"host_bind"`
	HostRotation     string        `json:"host_rotation"`
	PortBind         string        `json:"port_bind"`
	PortConn         string        `json:"port_conn"`
	Secure           bool          `json:"secure"`
	HostHeader       string        `json:"host_header"`
	Headers          []Header      `json:"headers"`
	Uris             []string      `json:"uris"`
	Certificate      Certificate   `json:"certificate"`
	WhitelistEnabled bool          `json:"whitelist_enabled"`
	Whitelist        []string      `json:"whitelist"`
	BlacklistEnabled bool          `json:"blacklist_enabled"`
	ProxySettings    ProxySettings `json:"proxy_settings"`
}

// Certificate holds the paths to the SSL certificate and its corresponding private key.
type Certificate struct {
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
}

// ProxySettings represents the configuration for a proxy server.
type ProxySettings struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Header represents a key-value pair used in HTTP headers.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HttpListener struct {
	Name      string
	Config    Config
	Server    *gin.Engine
	httpServe *http.Server
}

func NewHttpListener(name string, config Config) *HttpListener {
	server := gin.Default()

	server.Use(func(ctx *gin.Context) {
		if config.WhitelistEnabled {
			clientIp := ctx.ClientIP()
			if !isIpBlocked(clientIp, config.Whitelist) {
				logger.Warn("Ip Not Allowed:", clientIp)
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"message": "ip not allowed",
					"status":  status.ERROR,
				})
				return
			}
		}
		ctx.Next()
	})

	for _, uri := range config.Uris {
		server.GET(uri, func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Ok",
				"status":  status.SUCCESS,
			})
		})
	}

	return &HttpListener{
		Name:      name,
		Config:    config,
		Server:    server,
		httpServe: nil,
	}
}

func (hlisten *HttpListener) Start() error {
	address := hlisten.Config.HostBind + ":" + hlisten.Config.PortBind
	hlisten.httpServe = &http.Server{
		Addr:    address,
		Handler: hlisten.Server,
	}
	if hlisten.Config.Secure {
		certPath := hlisten.Config.Certificate.CertPath
		keyPath := hlisten.Config.Certificate.KeyPath

		if certPath == "" || keyPath == "" {
			return fmt.Errorf("SSL certificate or key path is missing.")
		}

		hlisten.httpServe.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		logger.Good("Starting Https server on ", address)
		return hlisten.httpServe.ListenAndServeTLS(certPath, keyPath)
	}
	logger.Good("Starting Http server on ", address)
	return hlisten.httpServe.ListenAndServe()
}

func (hlisten *HttpListener) Stop(timeout time.Duration) error {
	if hlisten.httpServe == nil {
		return fmt.Errorf("HTTP server is not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info("stopping server: ", hlisten.Name)
	if err := hlisten.httpServe.Shutdown(ctx); err != nil {
		logger.Error("Failed to gracefully shutdown server: ", err)
		return err
	}
	return nil
}

func isIpBlocked(ip string, blacklist []string) bool {
	for _, blockedIp := range blacklist {
		if strings.TrimSpace(blockedIp) == ip {
			return true
		}
	}
	return false
}
