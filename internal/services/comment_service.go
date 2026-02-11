package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrCommentNotFound     = errors.New("comment not found")
	ErrUnauthorizedComment = errors.New("unauthorized to modify this comment")
)

type CommentRepo interface {
	Create(ctx context.Context, comment *model.Comment) error
	FindByID(ctx context.Context, commentID string) (*model.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error)
	Delete(ctx context.Context, commentID string) error
	CountCommentsByPostID(ctx context.Context, postID string) (int64, error)
	GetAllComments(ctx context.Context, postID string) ([]*model.Comment, error)
}

type CommentService struct{
	commentRepo 	CommentRepo
	postRepo 		PostRepo
}

func NewCommentService(commentRepo CommentRepo, postRepo PostRepo) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// CreateComment - Create a new comment on a post
func (s *CommentService) CreateComment(ctx context.Context, req *model.CreateCommentRequest, postID, authorID  string)(*model.Comment, error){
	//Validate text
	text := strings.TrimSpace(req.Text)
	if text == ""{
		return nil , errors.New("comment text is required")
	}

	if len(text) < 1 {
		return nil, errors.New("comment must be at least 1 character long")
	}

	if len(text) > 1000 {
		return nil, errors.New("comment must not exceed 1000 characters")
	}

	// Verify post exists
	post, err := s.postRepo.FindByID(ctx , postID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify post: %w", err)
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	// Parse UUIDs
	var postUUID pgtype.UUID
	if err := postUUID.Scan(postID);err != nil {
		return nil, fmt.Errorf("invalid post ID format: %w", err)
	}

	var authorUUID pgtype.UUID
	if err := authorUUID.Scan(authorID); err != nil {
		return nil, fmt.Errorf("invalid author ID format: %w", err)
	}

	// Create comment model
	comment := &model.Comment{
		Text:     text,
		PostID:   postUUID,
		AuthorID: authorUUID,
	}

	// Save to database
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	log.Printf("[COMMENT-SERVICE] Comment created successfully: %s", comment.ID.String())
	return comment, nil

}

//GetCommentsByPostID - Get all comments for a specific post
func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID string) ([]*model.Comment, error) {
	// Verify post exists
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify post: %w", err)
	}
	if post == nil {
		return nil, ErrPostNotFound
	}

	// Get comments
	comments, err := s.commentRepo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	log.Printf("[COMMENT-SERVICE] Found %d comments for post: %s", len(comments), postID)
	return comments, nil
}

// DeleteComment - Delete a comment (only by author)
func (s *CommentService) DeleteComment(ctx context.Context, commentID string, userID string) error {
	// Find existing comment
	existing, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCommentNotFound
	}

	// Check ownership - only author can delete their comment
	if existing.AuthorID.String() != userID {
		return ErrUnauthorizedComment
	}

	// Delete from database
	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	log.Printf("[COMMENT-SERVICE] Comment deleted successfully: %s", commentID)
	return nil
}

func (s *CommentService) GetAllComments(ctx context.Context, postID string) ([]*model.Comment, error) {

	if postID == "" {
		return nil, errors.New("post ID is required")
	}

	comments, err := s.commentRepo.GetAllComments(ctx, postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// func (s *CommentService) GetCommentsByPostIDWithAuthor(ctx context.Context, postID string) ([]*model.CommentWithAuthor, error) {
// 	// Verify post exists
// 	post, err := s.postRepo.FindByID(ctx, postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if post == nil {
// 		return nil, ErrPostNotFound
// 	}

// 	// Get comments with author info
// 	return s.commentRepo.GetCommentsByPostID(ctx, postID)
// }