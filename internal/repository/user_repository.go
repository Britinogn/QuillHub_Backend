package repository

import (
	"context"
	"log"
	// "database/sql"
	"errors"
	"fmt"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/britinogn/quillhub/pkg/utils"

	//"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Save, Get, Update, Delete

// Generate UUID for new user
//user.ID = uuid.New().String()

type UserRepository struct{
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users(name, username, email, password, role)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `


    // Hash password
    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // Execute query and scan the returned values
    err = u.db.QueryRow(
        ctx,
        query,
        user.Name,
        user.Username,
        user.Email,
        hashedPassword,
        user.Role,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }

    return nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    query := `
        SELECT id, name, email, password, role
        FROM users
        WHERE email = $1
    `

    var user model.User

    err := u.db.QueryRow(ctx, query, email).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.Role,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to find user by email: %w", err)
    }

    return &user, nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    query := ` 
        SELECT id, name, username, password, role
        FROM users 
        WHERE username = $1
    `

    var user model.User
    err := u.db.QueryRow(ctx, query, username).Scan(
        &user.ID, 
        &user.Name,
        &user.Username,
        &user.Password,
        &user.Role,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to find user by username: %w", err)
    }

    return &user, nil
}


// GetOrCreateAIBot - Get existing AI bot or create new one
func (u *UserRepository) GetOrCreateAIBot(ctx context.Context) (string, error) {
	// Check if AI bot user exists
	query := `SELECT id FROM users WHERE email = 'aibot@quillhub.com' LIMIT 1`
	var userID string
	err := u.db.QueryRow(ctx, query).Scan(&userID)
	
	if err == nil {
		log.Printf("[USER-REPO] AI Bot user already exists: %s", userID)
		return userID, nil
	}
	
	// Create AI bot user
	insertQuery := `
		INSERT INTO users (name, username, email, password, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	
	// Hash password
	hashedPassword, err := utils.HashPassword("secure-bot-password-12345")
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	
	err = u.db.QueryRow(ctx, insertQuery, 
		"QuillHub AI Bot", 
		"quillhub_ai", 
		"aibot@quillhub.com", 
		hashedPassword, 
		"bot",
	).Scan(&userID)
	
	if err != nil {
		return "", fmt.Errorf("failed to create AI bot user: %w", err)
	}
	
	log.Printf("[USER-REPO] ✅ Created new AI Bot user: %s", userID)
	return userID, nil
}