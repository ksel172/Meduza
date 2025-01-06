package foreign

type Config struct {
	Endpoint       string         `json:"endpoint"`
	Authentication Authentication `json:"authentication"`
}

type Authentication struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
