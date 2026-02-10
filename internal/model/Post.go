package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Post - Database model
type Post struct {
	ID         	pgtype.UUID  	`json:"id" db:"id"`
	AuthorID   	pgtype.UUID  	`json:"author_id" db:"author_id"`
	Title      	string       	`json:"title" db:"title"`
	Content    	string       	`json:"content" db:"content"`
	ImageURL   	[]string     	`json:"image_url,omitempty" db:"image_url"`
	Tags       	[]string     	`json:"tags" db:"tags"`
	Likes    	*pgtype.UUID 	`json:"likes,omitempty" db:"likes"`      // Reference to likes table
	Comments 	*pgtype.UUID 	`json:"comments,omitempty" db:"comments"` // Reference to comments table
	Category 	*string 			`json:"category" db:"category"`
	IsPublished bool         	`json:"is_published" db:"is_published"`
	ViewCount 	int64   			`json:"view_count" db:"view_count"`
	CreatedAt 	time.Time    	`json:"created_at" db:"created_at"`
	UpdatedAt  	time.Time    	`json:"updated_at" db:"updated_at"`
}

// CreatePostRequest - For creating new posts
type CreatePostRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	//ImageURL *[]string  `json:"image_url,omitempty"`
	Category  *string 	`json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	AuthorID   string `json:"author_id"`
}

// UpdatePostRequest - For updating posts
type UpdatePostRequest struct {
	Title    *string  `json:"title,omitempty"`
	Content  *string  `json:"content,omitempty"`
	ImageURL []string  `json:"image_url,omitempty"`
	Category    *string  `json:"category"`
	Tags     []string `json:"tags,omitempty"`
	IsPublished *bool    `json:"is_published"`
}

// PostResponse - What to return to client
type PostResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	ImageURL  []string   `json:"image_url,omitempty"`
	Tags      []string  `json:"tags"`
	Category  *string 	`json:"category,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
