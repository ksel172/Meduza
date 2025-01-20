package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type ModuleConfig struct {
	Module interface{} `json:"module"`
}

func LoadModuleConfig(path string) (*ModuleConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ModuleConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadAssembly(filePath string) ([]byte, error) {
	assemblyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read assembly file: %w", err)
	}
	return assemblyBytes, nil
}
