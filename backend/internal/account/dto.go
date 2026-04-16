package account

// CreateUserRequest represents the data required to create a new user.
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Password     string `json:"password" binding:"required,min=6,max=50"`
	Email        string `json:"email" binding:"required,email"`
	Active       bool   `json:"active"`
	CurrentLevel string `json:"current_level" binding:"omitempty,oneof=P1 P2 P3 LT1 LT2 LT3 LT4"`
}

// UpdateUserRequest represents the data required to update a user.
type UpdateUserRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Password     string `json:"password" binding:"omitempty,min=6,max=50"` // opcional no update
	Email        string `json:"email" binding:"required,email"`
	Active       bool   `json:"active"`
	CurrentLevel string `json:"current_level" binding:"omitempty,oneof=P1 P2 P3 LT1 LT2 LT3 LT4"`
	ManagerName  string `json:"manager_name" binding:"omitempty,max=100"`
	ManagerEmail string `json:"manager_email" binding:"omitempty,email"`
}

// UserResponse represents the user data returned in API responses.
type UserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Active       bool   `json:"active"`
	Role         string `json:"role"`
	CurrentLevel string `json:"current_level,omitempty"`
	ManagerName  string `json:"manager_name,omitempty"`
	ManagerEmail string `json:"manager_email,omitempty"`
}

// LoginRequest represents the credentials required for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response data after successful login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
