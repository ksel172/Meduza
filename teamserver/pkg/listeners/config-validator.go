package listeners

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ksel172/Meduza/teamserver/pkg/listeners/foreign"
	"github.com/ksel172/Meduza/teamserver/pkg/listeners/http"
	"github.com/ksel172/Meduza/teamserver/pkg/listeners/smb"
	"github.com/ksel172/Meduza/teamserver/pkg/listeners/tcp"
)

// ConfigRegistry maps listener types to their corresponding struct types.
var ConfigRegistry = map[string]any{
	"http":    &http.Config{},
	"https":   &http.Config{},
	"h2c":     &http.Config{},
	"http2":   &http.Config{},
	"http3":   &http.Config{},
	"tcp":     &tcp.Config{},
	"smb":     &smb.Config{},
	"foreign": &foreign.Config{},
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
