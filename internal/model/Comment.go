package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Comment - Database model
type Comment struct {
	ID        pgtype.UUID `json:"id" db:"id"`
	PostID    pgtype.UUID `json:"post_id" db:"post_id"`
	AuthorID  pgtype.UUID `json:"author_id" db:"author_id"`
	Text      string      `json:"text" db:"text"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateCommentRequest - For creating new comments
type CreateCommentRequest struct {
	Text string `json:"text" binding:"required,min=1,max=1000"`
}

// UpdateCommentRequest - For updating comments
type UpdateCommentRequest struct {
	Text *string `json:"text" binding:"omitempty,min=1,max=1000"`
}

// CommentResponse - What to return to client
type CommentResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	AuthorID  string    `json:"author_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentWithAuthor - Comment with author details (for rich responses)
// type CommentWithAuthor struct {
// 	ID         string    `json:"id"`
// 	PostID     string    `json:"post_id"`
// 	AuthorID   string    `json:"author_id"`
// 	//AuthorName string    `json:"author_name"`
// 	Text       string    `json:"text"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	UpdatedAt  time.Time `json:"updated_at"`
// }