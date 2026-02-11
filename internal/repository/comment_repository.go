package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

// Create - Create a new comment
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (text, post_id, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	log.Printf("[COMMENT-REPO] Creating comment for post: %s", comment.PostID.String())

	err := r.db.QueryRow(
		ctx,
		query,
		comment.Text,
		comment.PostID,
		comment.AuthorID,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// GetAllComments - Get all comments (optionally by post_id)
func (r *CommentRepository) GetAllComments(ctx context.Context, postID string) ([]*model.Comment, error) {
	query := `
		SELECT id, text, post_id, author_id, created_at, updated_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var comment model.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.Text,
			&comment.PostID,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return comments, nil
}


// FindByID - Get a single comment by ID
func (r *CommentRepository) FindByID(ctx context.Context, commentID string) (*model.Comment, error) {
	query := `
		SELECT id, text, post_id, author_id, created_at, updated_at
		FROM comments 
		WHERE id = $1
	`

	var comment model.Comment
	err := r.db.QueryRow(ctx, query, commentID).Scan(
		&comment.ID,
		&comment.Text,
		&comment.PostID,
		&comment.AuthorID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	return &comment, nil
}

// GetCommentsByPostID - Get all comments for a specific post
func (r *CommentRepository) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	query := `
		SELECT id, text, post_id, author_id, created_at, updated_at
		FROM comments 
		WHERE post_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Text,
			&comment.PostID,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

// Update - Update a comment
func (r *CommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	query := `
		UPDATE comments
		SET text = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		comment.Text,
		comment.ID,
	).Scan(&comment.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("comment not found")
		}
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// Delete - Delete a comment
func (r *CommentRepository) Delete(ctx context.Context, commentID string) error {
	query := `DELETE FROM comments WHERE id = $1`

	result, err := r.db.Exec(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("comment not found")
	}

	log.Printf("[COMMENT-REPO] Comment deleted successfully: %s", commentID)
	return nil
}

// CountCommentsByPostID - Count total comments for a post
func (r *CommentRepository) CountCommentsByPostID(ctx context.Context, postID string) (int64, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var count int64
	err := r.db.QueryRow(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}

	return count, nil
}

// GetCommentsByPostIDWithAuthor - Get comments with author details
func (r *CommentRepository) GetCommentsByPostIDWithAuthor(ctx context.Context, postID string) ([]*model.CommentWithAuthor, error) {
	query := `
		SELECT 
			c.id, c.text, c.post_id, c.author_id, c.created_at, c.updated_at,
			u.name as author_name, u.username as author_username
		FROM comments c
		INNER JOIN users u ON c.author_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []*model.CommentWithAuthor
	for rows.Next() {
		var comment model.CommentWithAuthor
		err := rows.Scan(
			&comment.ID,
			&comment.Text,
			&comment.PostID,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.AuthorName,
			&comment.AuthorUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}