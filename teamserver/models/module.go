package models

type Module struct {
	Id          string    `json:"module_id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	FileName    string    `json:"file_name"`
	Usage       string    `json:"usage"`
	Commands    []Command `json:"commands"`
}

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ModuleConfig struct {
	Module Module `json:"Module"`
}
