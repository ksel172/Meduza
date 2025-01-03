package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var status = utils.Status

// Config represents the configuration settings for the service, including network, security, and proxy settings.
type Config struct {
	WorkingHours string   `json:"working_hours"`
	Hosts        []string `json:"hosts"`
	HostBind     string   `json:"host_bind"`
	HostRotation string   `json:"host_rotation"`
	PortBind     string   `json:"port_bind"`
	PortConn     string   `json:"port_conn"`
	Secure       bool     `json:"secure"`
	HostHeader   string   `json:"host_header"`
	Headers      []Header `json:"headers"`
	//Uris             []string      `json:"uris"`
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

// HttpListener manages the HTTP server with configuration and security options.
type HttpListener struct {
	Name        string
	Config      Config
	Server      *gin.Engine
	httpServe   *http.Server
	whitelistMu sync.RWMutex
}

var listenerRoute = "/check-in"

// NewHttpListener initializes a new HTTP listener.
func NewHttpListener(name string, config Config) (*HttpListener, error) {
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

	mux.Handle("GET", listenerRoute, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "OK")
	})
	mux.Handle("POST", listenerRoute, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Ok POST")
	})

	/*
			server.GET(uri, func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "Ok",
					"status":  status.SUCCESS,
				})
			})
		}
	*/

	return &HttpListener{
		Name:      name,
		Config:    config,
		Server:    mux,
		httpServe: nil,
	}, nil
}

// Validate ensures the configuration is valid before use.
func (config *Config) Validate() error {
	if config.HostBind == "" {
		return fmt.Errorf("HostBind is required")
	}
	if config.PortBind == "" {
		return fmt.Errorf("PortBind is required")
	}
	/*
		if len(config.Uris) == 0 {
			return fmt.Errorf("At least one URI is required")
		}
	*/
	if config.Secure {
		if config.Certificate.CertPath == "" || config.Certificate.KeyPath == "" {
			return fmt.Errorf("Certificate paths are required for secure mode")
		}
	}
	return nil
}

// Start begins the HTTP listener.
func (hlisten *HttpListener) Start() error {
	address := hlisten.Config.HostBind + ":" + hlisten.Config.PortBind
	hlisten.httpServe = &http.Server{
		Addr:    address,
		Handler: hlisten.Server,
	}

	errChan := make(chan error, 1)
	readyChan := make(chan struct{}, 1)

	go func() {
		readyChan <- struct{}{}

		var err error
		if hlisten.Config.Secure {
			certPath := hlisten.Config.Certificate.CertPath
			keyPath := hlisten.Config.Certificate.KeyPath

			if err := validateCertificate(certPath, keyPath); err != nil {
				return
			}

			hlisten.httpServe.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
			logger.Good("Starting HTTPS server on ", address)
			err = hlisten.httpServe.ListenAndServeTLS(certPath, keyPath)
		} else {
			logger.Good("Starting HTTP server on ", address)
			err = hlisten.httpServe.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	select {
	case <-readyChan:
		return nil
	case err := <-errChan:
		return fmt.Errorf("failed to start server: %v", err)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for server to start")
	}
}

// Stop gracefully shuts down the HTTP listener.
func (hlisten *HttpListener) Stop(timeout time.Duration) error {
	if hlisten.httpServe == nil {
		return fmt.Errorf("HTTP server is not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info("Stopping server:", hlisten.Name)
	if err := hlisten.httpServe.Shutdown(ctx); err != nil {
		logger.Error("Failed to gracefully shutdown server, forcing close:", err)
		return hlisten.httpServe.Close()
	}
	return nil
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

// validateCertificate checks if certificate files exist.
func validateCertificate(certPath, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("Certificate file not found: %s", certPath)
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("key file not found: %s", keyPath)
	}
	return nil
}
