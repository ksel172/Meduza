package utils

import (
	"fmt"
	"os"
)

// LoadAssembly loads a C# assembly from the given file path
func LoadAssembly(filePath string) ([]byte, error) {
	assemblyBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read assembly file: %w", err)
	}
	return assemblyBytes, nil
}
