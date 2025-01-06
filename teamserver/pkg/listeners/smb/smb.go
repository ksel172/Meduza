package smb

type Config struct {
	PipeName     string `json:"pipe_name"`
	MaxInstances int    `json:"max_instances"`
	KillDate     int64  `json:"kill_date"`
}
