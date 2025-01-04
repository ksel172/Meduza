package listeners

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"

// 	"github.com/ksel172/Meduza/teamserver/pkg/listeners/http"
// )

// type ListenersService struct{}

// func (s *ListenersService) CreateListener(listenerType string, config any) (ListenerRegistry, error) {
// 	parseConfig, err := ParseConfig(listenerType, config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	switch listenerType {
// 	case "http", "https", "http3", "h2c":
// 		httpConfig, ok := parseConfig.(*http.Config)
// 		if !ok {
// 			return nil, errors.New("parsed config is not of type *http.Config")
// 		}
// 		return http.NewHttpListener(httpConfig.HostHeader, *httpConfig)
// 	default:
// 		return nil, fmt.Errorf("unsupported listener type: %s", listenerType)
// 	}
// }

// // ParseConfig parses the raw config and identifies its type based on the listener type.
// // It validates and returns the parsed configuration or an error.
// func ParseConfig(listenerType string, rawConfig any) (any, error) {
// 	// Check if the listener type exists in the registry
// 	expectedType, ok := ConfigRegistry[listenerType]
// 	if !ok {
// 		return nil, fmt.Errorf("unsupported listener type: %s", listenerType)
// 	}

// 	// Clone the expected type for unmarshalling
// 	expectedConfig := cloneType(expectedType)

// 	// Convert the raw config to JSON
// 	configBytes, err := json.Marshal(rawConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to serialize raw config: %v", err)
// 	}

// 	// Decode JSON into the expected type without strict validation
// 	if err := json.Unmarshal(configBytes, expectedConfig); err != nil {
// 		return nil, fmt.Errorf(
// 			"failed to parse config for listener type '%s': %v",
// 			listenerType, err,
// 		)
// 	}

// 	return expectedConfig, nil
// }
