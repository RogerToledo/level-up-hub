package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" binding:"required,min=3,max=50"`
	Password  string    `json:"password,omitempty"`
	Email     string    `json:"email" binding:"required,email"`
	Active    bool      `json:"active"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
