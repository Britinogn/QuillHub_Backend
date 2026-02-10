package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/jackc/pgx/v5"
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
		INSERT INTO posts (title, content, image_url, tags, author_id, category)
		VALUES ($1, $2 , $3, $4, $5, $6)
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
		post.Category,    
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

func (r *PostRepository) GetAllPost(ctx context.Context, limit, offset int) ([]*model.Post, error){
	query := `
		SELECT id, title, content, author_id, image_url, tags, 
			category, is_published, view_count,
			likes, comments, created_at, updated_at		
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next(){
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.ImageURL,
			&post.Tags,
			&post.Category,
			&post.IsPublished,
			&post.ViewCount,
			&post.Likes,
			&post.Comments,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post %w", err)
		}

		posts = append(posts, &post)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}


func (r *PostRepository) CountPosts(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM posts"
	
	var count int64
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count posts: %w", err)
	}
	
	return count, nil
}

//FindByID - Get a post by ID
func (r *PostRepository) FindByID(ctx context.Context, postID string) (*model.Post, error) {
	query := `
		SELECT id, author_id, title, content, image_url, category, tags, 
		is_published, view_count, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	var post model.Post
	err := r.db.QueryRow(ctx, query, postID).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Content,
		&post.ImageURL,
		&post.Category,
		&post.Tags,
		&post.IsPublished,
		&post.ViewCount,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find post: %w", err)
	}

	return &post, nil
}


//FindByID - Get a auth by authorID
func (r *PostRepository) FindByAuthorID(ctx context.Context, authorID string) ([]*model.Post, error) {
	query := `
		SELECT id, author_id, title, content, image_url, category, tags, 
		is_published, view_count, created_at, updated_at
		FROM posts
		WHERE author_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user posts: %w", err)
	}
	defer rows.Close()
	var posts []*model.Post
	for rows.Next(){
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.ImageURL,
			&post.Category,
			&post.Tags,
			&post.IsPublished,
			&post.ViewCount,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}
	
	return posts, nil

}

//Update - Update a post
func (r *PostRepository) Update(ctx context.Context, post *model.Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, image_url = $3, category = $4, 
			tags = $5, is_published = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ImageURL,
		post.Category,
		post.Tags,
		post.IsPublished,
		post.ID,
	).Scan(&post.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("post not found")
		}
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

// Delete - Delete a post
func (r *PostRepository) Delete(ctx context.Context, postID string) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("post not found")
	}

	return nil
}

// IncrementViewCount - Increment view count
func (r *PostRepository) IncrementViewCount(ctx context.Context, postID string) error {
	query := `UPDATE posts SET view_count = view_count + 1 WHERE id = $1`

	_, err := r.db.Exec(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}

	return nil
}