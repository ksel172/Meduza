package auth

// AuthRequest is the auth request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse is the auth response
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
