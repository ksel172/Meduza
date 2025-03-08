package models

const (
	// URL parameter constants
	ParamModuleID string = "module_id"
)

type Module struct {
	Id          string    `json:"module_id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	FileName    string    `json:"file_name"`
	Commands    []Command `json:"commands"`
}

type ModuleBytes struct {
	ModuleBytes     []byte            `json:"module_bytes"`
	DependencyBytes map[string][]byte `json:"dependency_bytes"`
}

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModuleConfig struct {
	Module Module `json:"Module"`
}
