package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
    ID          pgtype.UUID    `json:"id" db:"id"`
    Name        string         `json:"name" db:"name"`
    Username    string         `json:"username" db:"username"`
    Email       string         `json:"email" db:"email"`
    Password    string         `json:"-" db:"password"`
    Role        string         `json:"role" db:"role"`
    Gender      *string        `json:"gender,omitempty" db:"gender"`
    ProfileURL  *string        `json:"profile_url,omitempty" db:"profile_url"`
    IsVerified  bool           `json:"is_verified" db:"is_verified"`
    IsActive    bool           `json:"is_active" db:"is_active"`
    Bio         *string        `json:"bio,omitempty" db:"bio"`
    LastLogin   *time.Time     `json:"last_login,omitempty" db:"last_login"`
    CreatedAt   time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest - for registration
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Username string `json:"username" binding:"required,min=3"` 
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Role     *string `json:"role"`  // ← Optional role field

}

// UserResponse - what to return to client (without sensitive data)
type UserResponse struct {
    ID         string     `json:"id"`
    Name       string     `json:"name"`
    Username   string     `json:"username"`
    Email      string     `json:"email"`
    Role       string     `json:"role"`
    ProfileURL *string    `json:"profile_url,omitempty"`
    CreatedAt  time.Time  `json:"created_at"`
}

// LoginRequest - for login (email OR username + password)
type LoginRequest struct {
    Identifier string `json:"identifier" binding:"required"` // can be email or username
    Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}