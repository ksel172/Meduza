package http_listener

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
	// "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
)

// HTTPListener implements ListenerImplementation for HTTP/HTTPS servers
type HTTPListener struct {
	Config HTTPListenerConfig

	isRunning bool
	mu        sync.RWMutex

	server            *http.Server
	router            *gin.Engine
	checkinController *checkin.CheckInController
	shutdownSignal    chan struct{}
}

// NewHTTPListener creates and returns a new HTTP listener
func NewHTTPListener(config HTTPListenerConfig, checkinController *checkin.CheckInController) (*HTTPListener, error) {
	listener := &HTTPListener{
		Config:            config,
		checkinController: checkinController,
	}

	if err := listener.configure(); err != nil {
		return nil, err
	}

	return listener, nil
}

func (l *HTTPListener) configure() error {
	if l.isRunning {
		return fmt.Errorf("listener is running, cannot configure")
	}

	// Initialize the router if not already done
	if l.router == nil {
		l.router = gin.Default()
		l.router.POST("/", l.HandleCheckIn)
	}

	// Configure the server
	address := fmt.Sprintf("%s:%d", l.Config.Host, l.Config.Port)
	l.server = &http.Server{
		Addr:         address,
		Handler:      l.router,
		ReadTimeout:  l.Config.ReadTimeout,
		WriteTimeout: l.Config.WriteTimeout,
	}

	return nil
}

// Start begins the HTTP listener
func (l *HTTPListener) Start(ctx context.Context) error {
	utils.AssertNotNil(l.server)

	l.mu.Lock()
	defer l.mu.Unlock()

	// Signal channel for startup completion
	errChan := make(chan error, 1)
	readyChan := make(chan struct{}, 1)
	l.shutdownSignal = make(chan struct{})

	// Start the server in a goroutine
	go func() {
		readyChan <- struct{}{}

		var err error
		if l.Config.EnableTLS {
			// Validate certificates
			if err := l.validateCertificate(); err != nil {
				errChan <- err
				return
			}

			l.server.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
			err = l.server.ListenAndServeTLS(l.Config.CertPath, l.Config.KeyPath)
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
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
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

// TODO
func (l *HTTPListener) Validate() error {
	return nil
}

// validateCertificate checks if certificate files exist and are valid
func (l *HTTPListener) validateCertificate() error {
	if !l.Config.EnableTLS {
		return nil
	}

	if l.Config.CertPath == "" || l.Config.KeyPath == "" {
		return fmt.Errorf("TLS enabled but certificate paths not specified")
	}

	if _, err := os.Stat(l.Config.CertPath); os.IsNotExist(err) {
		return fmt.Errorf("certificate file not found: %s", l.Config.CertPath)
	}

	if _, err := os.Stat(l.Config.KeyPath); os.IsNotExist(err) {
		return fmt.Errorf("key file not found: %s", l.Config.KeyPath)
	}

	return nil
}

// Extract all required values from request for core checkin functionality
func (l *HTTPListener) HandleCheckIn(ctx *gin.Context) {
	// Read request body to unmarshal into c2request
	// body might be encrypted with AES or not
	// depending on whether the agent is already authenticated)
	body, _ := io.ReadAll(ctx.Request.Body)
	if len(body) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}
	var c2request models.C2Request

	// If a session token is sent, that means the agent is authenticated
	// Otherwise, treat is as an authentication request
	sessionTokenBase64 := ctx.GetHeader("Session-Token")
	if sessionTokenBase64 == "" {
		// Get base64 Auth-Token header and decode it
		authTokenBase64 := ctx.GetHeader("Auth-Token")
		if authTokenBase64 == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing Auth-Token header"})
			return
		}
		authToken, err := base64.StdEncoding.DecodeString(authTokenBase64)
		if err != nil {
			log.Printf("failed to base64 decode Auth-Token header: %v", authToken)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Auth-Token header"})
			return
		}

		// The C2Request should not be encrypted at this point, only base64 encoded, unmarshal and use it to authenticate
		// The c2 request should contain only a single message field with the agent public key:
		// BASE 64 ENCODED REQUEST BODY:
		// {"message": <AGENT_PUBLIC_KEY>}
		decodedC2Request, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			log.Printf("failed to decode c2request: %v, data: %v", err, decodedC2Request)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid encrypted data"})
			return
		}
		var c2request models.C2Request
		if err := json.Unmarshal(decodedC2Request, &c2request); err != nil {
			log.Printf("Invalid request body: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// Get the agent public key from the request message
		agentPublicKey := c2request.Message
		if agentPublicKey == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing agent public key"})
			return
		}

		// Authenticate agent
		log.Printf("Handling authentication request for agent %s", c2request.AgentID)
		response, err := l.checkinController.Authenticate(agentPublicKey, string(authToken))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			log.Printf("failed to authenticate: %v", err)

		}
		log.Printf("Agent %s authenticated", c2request.AgentID)
		ctx.JSON(http.StatusAccepted, gin.H{
			"public_key":    base64.StdEncoding.EncodeToString(response.PublicKey),
			"session_token": base64.StdEncoding.EncodeToString([]byte(response.SessionToken)),
		})

		return
	}

	// Get session AES key for agent, decrypt request body and unmarshal
	sessionToken, err := base64.StdEncoding.DecodeString(sessionTokenBase64)
	if err != nil {
		log.Printf("failed to base64 decode Session-Token header: %v", sessionToken)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Auth-Token header"})
		return
	}
	key, exists := storage.KeyRegistry.GetKey(string(sessionToken))
	if !exists {
		// If agent received unauthorized response
		// it should send another authentication request right after
		log.Printf("Missing AES key for session: %s", sessionToken)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session token"})
		return
	}
	decryptedData, err := utils.AesDecrypt(key, body)
	if err != nil {
		log.Printf("Failed to decrypt request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "decryption failed"})
		return
	}
	if err := json.Unmarshal(decryptedData, &c2request); err != nil {
		log.Printf("Invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	switch c2request.Reason {
	case models.Task:
		log.Printf("Handling task request for agent %s", c2request.AgentID)
		taskResponse, err := l.checkinController.HandleTaskRequest(ctx.Request.Context(), c2request, string(sessionToken))
		if err != nil {
			switch err {
			case checkin.ErrDatabase:
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			case checkin.ErrInternalServer:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			case checkin.ErrUnauthorized:
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"message": taskResponse})
		return
	case models.Response:
		// log.Printf("Handling response request for agent %s", c2request.AgentID)
		// l.checkinController.HandleResponseRequest(ctx, c2request)
		return
	case models.Register:
		// log.Printf("Handling register request for agent %s", c2request.AgentID)
		// l.checkinController.HandleRegisterRequest(ctx, c2request)
		return
	}
}
