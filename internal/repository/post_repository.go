package repository


import (
	"context"
	"fmt"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepository struct {
    db *pgxpool.Pool
}

func NewPostRepository(db *pgxpool.Pool) *PostRepository{
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (title, content, image_url, tags, author_id)
		VALUES ($1, $2 , $3, $4, $5)
		RETURNING id , created_at , updated_at
	`
	// Execute query and scan the returned values
	err := r.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ImageURL,
		post.Tags,
		post.AuthorID,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}