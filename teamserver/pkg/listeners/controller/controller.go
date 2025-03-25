package controller

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// )

// type Controller struct {
// 	config  ControllerConfig // Controler config
// 	manager *Manager         // Handles managed listeners
// 	client  *http.Client     // For submitting requests to the C2 server
// 	server  server.Server    // Rest API router -- to receive requests from the C2 server

// 	stopChan chan bool

// 	// For caching key pairs
// 	// Will a cache even be useful really? It will reduce the amount of times the keys are sent over the network
// 	// but is there a point? the keys are the same for the same payload so maybe there will be many agents
// 	// reconnecting at once when the listener is restarted, so it could be useful.
// 	// KeyCache     *cache.Cache
// }

// // Important considerations
// //  1. Controller should register with server, if not succesfull, should clean up all of its resources and kill itself.
// //  2. C2 Server will expect a response with 30 seconds to 1 minute (??), if none is received, it will assume the controller could
// //     not be launched for some reason or failed to commnicate with the server and shut itself down
// func NewController(config ControllerConfig, manager *Manager) (*Controller, error) {
// 	// Setup controller
// 	controller := &Controller{
// 		config:  config,
// 		manager: manager,
// 		client:  &http.Client{Timeout: 30 * time.Second},
// 	}

// 	// Setup REST API server
// 	controller.setupServer(config.ListenAddr)

// 	// With everything finally setup, submit a register request to the c2 server
// 	if err := controller.RegisterWithServer(); err != nil {
// 		return nil, fmt.Errorf("failed to register with server: %w", err)
// 	}

// 	return controller, nil
// }

// func (c *Controller) Run() {

// 	// Start hearbeat loop
// 	go func() {
// 		c.heartbeatLoop()
// 	}()

// 	defer func() {
// 		if err := recover(); err != nil {
// 			// c.Shutdown()
// 		}
// 		panic("server error")
// 	}()
// 	if err := c.server.Run(); err != nil {
// 		log.Fatalf("Controller server failed: %v", err)
// 	}
// }

// func (c *Controller) GetManager() *Manager {
// 	return c.manager
// }

// func (c *Controller) SendTestHeartbeat(ctx context.Context) error {
// 	return c.sendHeartbeat(ctx)
// }

// // Stop will stop all listeners controller is responsible for
// // func (c *Controller) Stop() error {
// // 	var wg sync.WaitGroup
// // 	for id, listener := range c.manager.listeners {
// // 		wg.Add(1)
// // 		go func(l *Listener) {
// // 			defer wg.Done()
// // 			if err := l.Stop(); err != nil {
// // 				log.Printf("failed to stop listener with id: %s", id)
// // 			}
// // 		}(listener)
// // 	}
// // 	wg.Wait()
// // 	log.Println("All listeners shut down gracefully")
// // 	return nil
// // }

// // Shutdown will stop all listeners from running and the controller itself
// // func (c *Controller) Shutdown() {
// // 	var wg sync.WaitGroup
// // 	for _, listener := range c.manager.listeners {
// // 		wg.Add(1)
// // 		go func(l *Listener) {
// // 			defer wg.Done()
// // 			l.Shutdown()
// // 		}(listener)
// // 	}
// // 	wg.Wait()
// // 	log.Println("All listeners shut down gracefully")

// // 	if err := c.server.Shutdown(); err != nil {
// // 		log.Print("failed to shutdown server gracefully, forcing close: %v", err)
// // 		c.server.Kill()
// // 	}

// // 	os.Exit(1)
// // }
