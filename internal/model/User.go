package model

import "time"

type User struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    Username    string     `json:"username"`
    Email       string     `json:"email"`
    Password    string     `json:"-"`                    // "-" means don't return in JSON responses
    Role        string     `json:"role"`
    Gender      *string    `json:"gender,omitempty"`     // pointer = nullable field
    ProfileURL  *string    `json:"profile_url,omitempty"`
    IsVerified  bool       `json:"is_verified"`
    IsActive    bool       `json:"is_active"`
    Bio         *string    `json:"bio,omitempty"`
    LastLogin   *time.Time `json:"last_login,omitempty"`
    UpdatedAt   time.Time  `json:"updated_at"`
    CreatedAt   time.Time  `json:"created_at"`
}

// CreateUserRequest - for registration
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest - for login (email OR username + password)
type LoginRequest struct {
    Identifier string `json:"identifier" binding:"required"` // can be email or username
    Password   string `json:"password" binding:"required"`
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