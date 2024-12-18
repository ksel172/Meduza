package tcp

type Config struct {
	PortBind   string `json:"port_bind"`
	HostBind   string `json:"host_bind"`
	BufferSize int    `json:"buffer_size"`
	Timeout    int    `json:"timeout"`
}
