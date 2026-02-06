package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/britinogn/quillhub/pkg/utils"
)

var (
    ErrEmailAlreadyRegistered = errors.New("email already registered")
    ErrInvalidInput      = errors.New("invalid registration data")
    ErrWeakPassword      = errors.New("password is too weak")
    ErrUsernameTaken     = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid email/username or password")
	ErrInvalidToken = errors.New("invalid token")
    ErrDatabaseOperation = errors.New("database operation failed")
)
type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error) 
}

type AuthService struct {
	repo UserRepo
}

func NewAuthService(repo UserRepo) *AuthService{
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(ctx context.Context, user *model.User) error {
	if user == nil {
		return ErrInvalidInput
	}

	// Normalize users
	name := strings.TrimSpace(user.Name)
	username := strings.TrimSpace(user.Username)
	email    := strings.ToLower(strings.TrimSpace(user.Email))
	password := strings.TrimSpace(user.Password)


	//   ONE combined check for all required fields
	if name == "" ||
		username == "" ||
		email == "" ||
		password == "" {
		return errors.New("name, username, email, and password are required")
	}

	// Optional: very basic extra rules 
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")  // ← Fixed typo
	}
	if !strings.Contains(email, "@"){
		return errors.New("invalid email format")
	}

	// Check if username already exists
	existingUsername, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to check username: %w", err)
	}
	if existingUsername != nil {
		return ErrUsernameTaken
	}

	// Check if email already exists
	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if existingUser != nil {
		return ErrEmailAlreadyRegistered
	}

	// Update the passed user object directly
	user.Name = name
	user.Username = username
	user.Email = email

	if user.Role == "" {
		user.Role = "user"
	}

	// Create user - this will populate ID and CreatedAt
	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, identifier, password string) (*model.User, string, error) {
	if identifier == "" || password == "" {
		return nil, "",  ErrInvalidCredentials
	}

	identifier = strings.TrimSpace(identifier)

	var user *model.User
	var err error

	//Try as email first
	if strings.Contains(identifier, "@"){
		user, err = s.repo.FindByEmail(ctx , strings.ToLower(identifier))
	}

	// If not found or not email → try as username
	if user == nil {
		user, err = s.repo.FindByUsername(ctx , identifier)
	}

	if err != nil {
        return nil, "", fmt.Errorf("login failed: %w", err)
    }

	if user == nil {
        return nil, "", ErrInvalidCredentials
    }

	//Check hash password
	if !utils.CheckPasswordHash(password, user.Password){
		return nil, "", ErrInvalidCredentials
	}
	user.Password = ""


	// Generate token 
    // token, err := utils.GenerateToken(user.ID, user.Email, user.Username, user.Role)
    // if err != nil {
    //     return nil, "", fmt.Errorf("failed to generate token: %w", err)
    // }

	// Generate token 
    token, err := utils.GenerateToken(user.ID.String(), user.Email, user.Username, user.Role)
    if err != nil {
        return nil, "", fmt.Errorf("failed to generate token: %w", err)
    }


	return user, token, nil

}

// RegisterAdmin - Only callable by existing admins
func (s *AuthService) RegisterAdmin(ctx context.Context, user *model.User, requestingUserRole string) error {
    // Check if the requesting user is an admin
    if requestingUserRole != "admin" {
        return errors.New("unauthorized: only admins can create admin users")
    }

    // Same validation as Register
    if user == nil {
        return ErrInvalidInput
    }

    name := strings.TrimSpace(user.Name)
    username := strings.TrimSpace(user.Username)
    email := strings.ToLower(strings.TrimSpace(user.Email))
    password := strings.TrimSpace(user.Password)

    if name == "" || username == "" || email == "" || password == "" {
        return errors.New("name, username, email, and password are required")
    }

    if len(username) < 3 {
        return errors.New("username must be at least 3 characters")
    }
    if len(password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }

    // Check username
    existingUsername, err := s.repo.FindByUsername(ctx, username)
    if err != nil {
        return fmt.Errorf("failed to check username: %w", err)
    }
    if existingUsername != nil {
        return ErrUsernameTaken
    }

    // Check email
    existingUser, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        return fmt.Errorf("failed to check email: %w", err)
    }
    if existingUser != nil {
        return ErrEmailAlreadyRegistered
    }

    // Set role to admin
    user.Name = name
    user.Username = username
    user.Email = email
    user.Role = "admin"  // ← Force admin role

    if err := s.repo.Create(ctx, user); err != nil {
        return fmt.Errorf("failed to create admin user: %w", err)
    }

    return nil
}