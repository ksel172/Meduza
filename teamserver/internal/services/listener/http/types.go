package http_listener

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	// "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
)

// Function type for starting the server
type serverStartFunc func(errChan chan<- error)

// HTTPListener implements ListenerImplementation for HTTP/HTTPS servers
type HTTPListener struct {
	Config HTTPListenerConfig

	isRunning bool
	mu        sync.RWMutex

	server *http.Server
	// router            *gin.Engine
	checkinController checkin.ICheckInController

	startServerGoroutine serverStartFunc
}

// NewHTTPListener creates and returns a new HTTP listener
func NewHTTPListener(config HTTPListenerConfig, checkinController checkin.ICheckInController) (*HTTPListener, error) {
	listener := &HTTPListener{
		Config:            config,
		checkinController: checkinController,
	}

	if err := listener.configure(); err != nil {
		return nil, err
	}

	return listener, nil
}

func (l *HTTPListener) SetListenerContext(ctx context.Context, cancelFunc context.CancelFunc) {

}

func (l *HTTPListener) configure() error {
	if l.isRunning {
		return fmt.Errorf("listener is running, cannot configure")
	}

	// Initialize the router if not already done
	router := gin.Default()
	router.POST("/", l.HandleCheckIn)

	// Configure the server
	address := fmt.Sprintf("%s:%d", l.Config.Host, l.Config.Port)
	l.server = &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  l.Config.ReadTimeout,
		WriteTimeout: l.Config.WriteTimeout,
	}

	// Set default server start goroutine
	l.startServerGoroutine = l.defaultStartServer

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
