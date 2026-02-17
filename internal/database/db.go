package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// DB is the global database connection pool
var DB *pgxpool.Pool

func init() {
	paths := []string{".", ".env", "../.env", "../../.env"}
	
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("✓ Loaded .env from: %s", path)
			return
		}
	}
	
	log.Println("Warning: .env file not found in any location")
}

func ConnectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	// Read env vars
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Set default for sslmode if not provided
	if sslmode == "" {
		sslmode = "require"
	}

	// Validate required env vars
	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	// Build DSN
	// dsn := fmt.Sprintf(
	// 	"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
	// 	user,
	// 	password,
	// 	host,
	// 	port,
	// 	dbName,
	// 	sslmode,
	// )

	// Replace your DSN build with this
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslmode,
	)

	// Parse pgx config
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	// //Force IPv4
	// config.ConnConfig.Config.Host = host
	// config.ConnConfig.Config.Port = 5432

	// Configure connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Store in global variable
	DB = pool

	log.Println("✓ Postgres connected successfully")
	return pool, nil
}

// Close gracefully closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("✓ Database connection closed")
	}
}