package http_listener

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// This goroutine always exists while the server is running
// Whenever a Stop/Terminate function is called, it will stop the underlying server
// That will cause this function to exit the listen and serve loop
// Pushing the shutdown signal to the Stop/Terminate functions
func (l *HTTPListener) defaultStartServer(errChan chan<- error) {
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
		logger.Info("launching server...")
		err = l.server.ListenAndServe()
	}

	// If server exits with an error other than shutdown, report it
	if err != nil && err != http.ErrServerClosed {
		errChan <- err
	}
}

// Start begins the HTTP listener
func (l *HTTPListener) Start(ctx context.Context) error {
	utils.AssertNotNil(l.server)

	l.mu.Lock()
	defer l.mu.Unlock()

	// Signal channel for startup completion
	errChan := make(chan error, 1)

	// Start server goroutine
	logger.Info("Launching HTTP listener server...")
	go l.startServerGoroutine(errChan)

	// Monitor the server start
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Try to connect to the server
			dialer := net.Dialer{
				Timeout: 1 * time.Second,
			}
			conn, err := dialer.DialContext(ctx, "tcp", l.server.Addr)
			logger.Info(fmt.Sprintf("attemped to dial server address %s: %v", l.server.Addr, err))
			if err == nil {
				conn.Close()
				l.isRunning = true
				return nil
			}
		case err := <-errChan:
			if err := l.shutdown(); err != nil {
				errChan <- fmt.Errorf("error during server start: %w", err)
			}
			return fmt.Errorf("failed to start server: %w", err)
		case <-ctx.Done():
			if err := l.shutdown(); err != nil {
				return fmt.Errorf("failed to shutdown server, should not have started: %w", err)
			}
			return fmt.Errorf("timed out waiting for server to start: %w", ctx.Err())
		}
	}
}

// Stop gracefully shuts down the HTTP listener
func (l *HTTPListener) Stop(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	utils.AssertNotNil(l.server)

	if err := l.shutdown(); err != nil {
		return err
	}

	return nil
}

// Terminate forcefully closes the HTTP listener
// Exists only for interface purposes, useful for external listeners/terminating the process
// For local listeners, it means it will get removed from the ListenerRegistry in the service level as well
func (l *HTTPListener) Terminate(ctx context.Context) error {
	return l.Stop(ctx)
}

//  TODO: This should not be implemented at this level like this tbh
// UpdateConfig applies configuration changes to the HTTP listener
// func (l *HTTPListener) UpdateConfig(ctx context.Context) error {
// 	// For configuration changes that require a restart
// 	isRunning := false

// 	l.mu.RLock()
// 	if l.isRunning {
// 		isRunning = true
// 	}
// 	l.mu.RUnlock()

// 	// If running, stop and restart to apply new config
// 	if isRunning {
// 		if err := l.Stop(ctx); err != nil {
// 			return fmt.Errorf("failed to stop listener for config update: %w", err)
// 		}

// 		if err := l.Start(ctx); err != nil {
// 			return fmt.Errorf("failed to restart listener after config update: %w", err)
// 		}
// 	}

//		return nil
//	}
func (l *HTTPListener) UpdateConfig(ctx context.Context) error { return nil }

func (l *HTTPListener) shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := l.server.Shutdown(ctx); err != nil {
		if err := l.server.Close(); err != nil {
			return fmt.Errorf("failed to close server: %w", err)
		}
	}

	l.isRunning = false
	return nil
}
