package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//var db *pgxpool.Pool

// RunMigrations - Create tables if they don't exist
func RunMigrations(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("🔄 Running database migrations...")

	migrations := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(250) NOT NULL,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'user',
		gender VARCHAR(25),
		profile_url VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT check_role CHECK (role IN ('user', 'admin', 'moderator'))
	);

	-- Posts table
	CREATE TABLE IF NOT EXISTS posts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		image_url TEXT[],
		tags TEXT[],
		category VARCHAR(100),
		is_published BOOLEAN DEFAULT true,
		view_count INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Comments table
	CREATE TABLE IF NOT EXISTS comments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		text TEXT NOT NULL,
		post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Likes table
	CREATE TABLE IF NOT EXISTS likes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
		comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
		author_id UUID REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT one_like_per_user_per_post UNIQUE (author_id, post_id),
		CONSTRAINT one_like_per_user_per_comment UNIQUE (author_id, comment_id)
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_posts_author ON posts(author_id);
	CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id);
	CREATE INDEX IF NOT EXISTS idx_comments_user ON comments(author_id);
	CREATE INDEX IF NOT EXISTS idx_likes_post ON likes(post_id);
	CREATE INDEX IF NOT EXISTS idx_likes_user ON likes(author_id);
	`

	_, err := db.Exec(ctx, migrations)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}