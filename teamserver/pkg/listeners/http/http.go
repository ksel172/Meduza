package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

// HttpListener manages the HTTP server with configuration and security options.
type HTTPListenerController struct {
	Name        string
	Config      models.HTTPListenerConfig
	Server      *gin.Engine
	HTTPServe   *http.Server
	WhitelistMu sync.RWMutex
}

// type ICheckInController interface {
// 	CreateAgent(ctx *gin.Context)
// 	GetTasks(ctx *gin.Context)
// }

// Start begins the HTTP listener.
func (c *HTTPListenerController) Start() error {
	address := c.Config.HostBind + ":" + c.Config.PortBind
	c.HTTPServe = &http.Server{
		Addr:    address,
		Handler: c.Server,
	}

	errChan := make(chan error, 1)
	readyChan := make(chan struct{}, 1)

	go func() {
		readyChan <- struct{}{}

		var err error
		if c.Config.Secure {
			certPath := c.Config.Certificate.CertPath
			keyPath := c.Config.Certificate.KeyPath

			if err := validateCertificate(certPath, keyPath); err != nil {
				return
			}

			c.HTTPServe.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
			logger.Good("Starting HTTPS server on ", address)
			err = c.HTTPServe.ListenAndServeTLS(certPath, keyPath)
		} else {
			logger.Good("Starting HTTP server on ", address)
			err = c.HTTPServe.ListenAndServe()
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
func (c *HTTPListenerController) Stop(timeout time.Duration) error {
	if c.HTTPServe == nil {
		return fmt.Errorf("HTTP server is not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info("Stopping server:", c.Name)
	if err := c.HTTPServe.Shutdown(ctx); err != nil {
		logger.Error("Failed to gracefully shutdown server, forcing close:", err)
		return c.HTTPServe.Close()
	}
	return nil
}

func (c *HTTPListenerController) GetName() string {
	return c.Name
}

// validateCertificate checks if certificate files exist.
func validateCertificate(certPath, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("certificate file not found: %s", certPath)
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("key file not found: %s", keyPath)
	}
	return nil
}
