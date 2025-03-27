package http_listener

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPListener implements ListenerImplementation for HTTP/HTTPS servers
type HTTPListener struct {
	// Configuration
	Host string
	Port int

	// For now it loads from path, I am not sure about the other ways we can approach
	// this yet. I assume we can also load it from memory.
	EnableTLS bool
	CertPath  string
	KeyPath   string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Runtime state
	server         *http.Server
	router         *gin.Engine
	isRunning      bool
	mu             sync.RWMutex
	shutdownSignal chan struct{}
}

// NewHTTPListener creates and returns a new HTTP listener
func NewHTTPListener(host string, port int, enableTLS bool, certPath, keyPath string) *HTTPListener {
	return &HTTPListener{
		Host:         host,
		Port:         port,
		EnableTLS:    enableTLS,
		CertPath:     certPath,
		KeyPath:      keyPath,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// Start begins the HTTP listener
func (l *HTTPListener) Start(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.isRunning {
		return fmt.Errorf("listener is already running")
	}

	// Initialize the router if not already done
	if l.router == nil {
		l.router = gin.Default()
		l.router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})
		//We are going to need a checkin endpoint that acts like a redirector
		// for agent requests. The endpoint shouldn't do any sort of decryption
		// except maybe authenticating the agent session token. To make sure
		// redirection is appropriate.

		// l.router.GET("/checkin", func(c *gin.Context) {
		// 	c.JSON(http.StatusOK, gin.H{"status": "checked in"})
		// })
	}

	// Configure the server
	address := fmt.Sprintf("%s:%d", l.Host, l.Port)
	l.server = &http.Server{
		Addr:         address,
		Handler:      l.router,
		ReadTimeout:  l.ReadTimeout,
		WriteTimeout: l.WriteTimeout,
	}

	// Signal channel for startup completion
	errChan := make(chan error, 1)
	readyChan := make(chan struct{}, 1)
	l.shutdownSignal = make(chan struct{})

	// Start the server in a goroutine
	go func() {
		readyChan <- struct{}{}

		var err error
		if l.EnableTLS {
			// Validate certificates
			if err := l.validateCertificate(); err != nil {
				errChan <- err
				return
			}

			l.server.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
			err = l.server.ListenAndServeTLS(l.CertPath, l.KeyPath)
		} else {
			err = l.server.ListenAndServe()
		}

		// If server exits with an error other than shutdown, report it
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}

		close(l.shutdownSignal)
	}()

	// Wait for either ready signal or error
	select {
	case <-ctx.Done():
		go l.server.Close() // Force close since we're abandoning it
		return ctx.Err()
	case <-readyChan:
		l.isRunning = true
		return nil
	case err := <-errChan:
		return fmt.Errorf("failed to start HTTP listener: %w", err)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for HTTP listener to start")
	}
}

// Stop gracefully shuts down the HTTP listener
func (l *HTTPListener) Stop(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isRunning || l.server == nil {
		return fmt.Errorf("listener is not running")
	}

	// Create a timeout context if the provided context doesn't have one
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	err := l.server.Shutdown(shutdownCtx)
	if err != nil {
		l.server.Close()
		l.isRunning = false
		return fmt.Errorf("forced shutdown of HTTP listener: %w", err)
	}

	// Wait for server goroutine to finish or timeout
	select {
	case <-l.shutdownSignal:
		l.isRunning = false
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Terminate forcefully closes the HTTP listener
func (l *HTTPListener) Terminate(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isRunning || l.server == nil {
		return nil // Already terminated or never started
	}

	err := l.server.Close()
	l.isRunning = false
	return err
}

// UpdateConfig applies configuration changes to the HTTP listener
func (l *HTTPListener) UpdateConfig(ctx context.Context) error {
	// For configuration changes that require a restart
	isRunning := false

	l.mu.RLock()
	if l.isRunning {
		isRunning = true
	}
	l.mu.RUnlock()

	// If running, stop and restart to apply new config
	if isRunning {
		if err := l.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop listener for config update: %w", err)
		}

		if err := l.Start(ctx); err != nil {
			return fmt.Errorf("failed to restart listener after config update: %w", err)
		}
	}

	return nil
}

// validateCertificate checks if certificate files exist and are valid
func (l *HTTPListener) validateCertificate() error {
	if !l.EnableTLS {
		return nil
	}

	if l.CertPath == "" || l.KeyPath == "" {
		return fmt.Errorf("TLS enabled but certificate paths not specified")
	}

	if _, err := os.Stat(l.CertPath); os.IsNotExist(err) {
		return fmt.Errorf("certificate file not found: %s", l.CertPath)
	}

	if _, err := os.Stat(l.KeyPath); os.IsNotExist(err) {
		return fmt.Errorf("key file not found: %s", l.KeyPath)
	}

	return nil
}
