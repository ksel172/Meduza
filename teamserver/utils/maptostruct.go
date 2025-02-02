package utils

import (
	"encoding/json"
)

// MapToStruct converts a map to a struct
func MapToStruct(input interface{}, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}
