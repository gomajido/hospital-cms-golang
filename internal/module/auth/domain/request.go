package domain

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=8"`
	Name      string   `json:"name" validate:"required"`
	Phone     string   `json:"phone" validate:"required"`
	RoleNames []string `json:"role_names" validate:"required,min=1"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
