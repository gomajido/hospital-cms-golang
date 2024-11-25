package domain

// RegisterResponse represents the registration response
type RegisterResponse struct {
	User  *User    `json:"user"`
	Roles []string `json:"roles"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
