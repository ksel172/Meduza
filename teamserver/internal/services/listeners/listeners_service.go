package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/listeners"
)

type ListenersService struct {
	controller        listeners.ListenerController // controller is used internally to launch and manage listeners
	checkinController *CheckInController
	registry          *listeners.Registry
}

func NewListenerService(checkInController *CheckInController) *ListenersService {
	listenerService := &ListenersService{
		registry:          listeners.NewRegistry(),
		checkinController: checkInController,
	}

	// Create default http controller
	// TODO: add default listener HTTP Listener config
	listenerService.CreateListenerController("http", &models.HTTPListenerConfig{UserAgent: "default"})

	return listenerService
}

/*
	ListenerService should always maintain at least one controller open, by default it is HTTP
*/

// Given the desired listenerType, create a new controller
// The controller creates and manages the listener server.
// The latest controller created is then saved in the controller field.
// Following calls to the service will be passed into the controller implementation of the function
func (s *ListenersService) CreateListenerController(listenerType string, config any) error {
	parseConfig, err := ParseConfig(listenerType, config)
	if err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}
	switch listenerType {
	case "http", "https", "http3", "h2c":

		// Validate the config
		httpConfig, ok := parseConfig.(*models.HTTPListenerConfig)
		if !ok {
			return errors.New("parsed config is not of type *http.Config")
		}

		// Create HTTP controller and save to controller field
		controller, err := NewHTTPListenerController(httpConfig.UserAgent, *httpConfig, s.checkinController)
		if err != nil {
			return fmt.Errorf("failed to create HTTP listener controller: %v", err)
		}
		s.controller = controller

		return nil
	default:
		return fmt.Errorf("unsupported listener type: %s", listenerType)
	}
}

// Each controller implements its own Start function
func (s *ListenersService) Start(listener models.Listener) error {
	if err := s.controller.Start(); err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}
	// add to registry
	s.registry.AddListener(listener)
	return nil
}

func (s *ListenersService) Stop(listenerID string, timeout time.Duration) error {

	// check if listener is running
	listener, exists := s.registry.GetListener(listenerID)
	if !exists {
		return fmt.Errorf("listener with ID %s does not exist", listener.ID.String())
	}

	// stop server
	if err := s.controller.Stop(timeout); err != nil {
		return fmt.Errorf("failed to stop listener: %v", err)
	}

	// remove from registry
	s.registry.RemoveListener(listener.ID.String())
	return nil
}

func (s *ListenersService) GetListener(id string) (models.Listener, bool) {
	return s.registry.GetListener(id)
}

func (s *ListenersService) AddListener(listener models.Listener) {
	s.registry.AddListener(listener)
}

// Listener config helper functions below

// ConfigRegistry maps listener types to their corresponding struct types.
var ConfigRegistry = map[string]any{
	"http":    &models.HTTPListenerConfig{},
	"https":   &models.HTTPListenerConfig{},
	"h2c":     &models.HTTPListenerConfig{},
	"http2":   &models.HTTPListenerConfig{},
	"http3":   &models.HTTPListenerConfig{},
	"tcp":     &models.TCPListenerConfig{},
	"smb":     &models.SMBListenerConfig{},
	"foreign": &models.ForeignListenerConfig{},
}

// ValidateAndParseConfig validates and parses the raw config based on the listener type.
// Returns the parsed config or an error.
func ValidateAndParseConfig(listenerType string, rawConfig any) (any, error) {

	// Check if the listener type exists in the registry
	expectedType, ok := ConfigRegistry[listenerType]
	if !ok {
		return nil, fmt.Errorf("unsupported listener type: %s", listenerType)
	}

	// Clone the expected type for unmarshalling
	expectedConfig := cloneType(expectedType)

	// Convert the raw config to JSON and unmarshal into the expected type
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize raw config: %v", err)
	}
	// Decode JSON into the expected type, ensuring strict validation
	decoder := json.NewDecoder(bytes.NewReader(configBytes))
	decoder.DisallowUnknownFields() // Reject unknown fields
	if err := decoder.Decode(expectedConfig); err != nil {
		return nil, fmt.Errorf(
			"invalid config for listener type '%s': %v",
			listenerType, err,
		)
	}

	return expectedConfig, nil
}

// cloneType creates a new instance of the type pointed to by 'original'.
func cloneType(original any) any {
	if original == nil {
		return nil
	}
	return reflect.New(reflect.TypeOf(original).Elem()).Interface()
}

// ParseConfig parses the raw config and identifies its type based on the listener type.
// It validates and returns the parsed configuration or an error.
func ParseConfig(listenerType string, rawConfig any) (any, error) {
	// Check if the listener type exists in the registry
	expectedType, ok := ConfigRegistry[listenerType]
	if !ok {
		return nil, fmt.Errorf("unsupported listener type: %s", listenerType)
	}

	// Clone the expected type for unmarshalling
	expectedConfig := cloneType(expectedType)

	// Convert the raw config to JSON
	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize raw config: %v", err)
	}

	// Decode JSON into the expected type without strict validation
	if err := json.Unmarshal(configBytes, expectedConfig); err != nil {
		return nil, fmt.Errorf(
			"failed to parse config for listener type '%s': %v",
			listenerType, err,
		)
	}

	return expectedConfig, nil
}

// GetConfigDetails takes the parsed configuration and retrieves its details as a map.
func GetConfigDetails(parsedConfig any) (map[string]any, error) {
	if parsedConfig == nil {
		return nil, fmt.Errorf("parsedConfig is nil or invalid")
	}

	configBytes, err := json.Marshal(parsedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize parsed config: %v", err)
	}

	// Convert the JSON back into a map for detailed inspection
	var configDetails map[string]any
	if err := json.Unmarshal(configBytes, &configDetails); err != nil {
		return nil, fmt.Errorf("failed to parse config details: %v", err)
	}

	return configDetails, nil
}
