package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"mime/multipart"
	"strings"

	// "time"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrPostNotFound      = errors.New("post not found")
	ErrUnauthorizedPost  = errors.New("unauthorized to modify this post")
)

type PostRepo interface{
	Create(ctx context.Context, post *model.Post) error 
	GetAllPost(ctx context.Context, limit, offset int) ([]*model.Post, error)
	CountPosts(ctx context.Context) (int64, error)
	FindByID(ctx context.Context, postID string) (*model.Post, error)
	FindByAuthorID(ctx context.Context, authorID string) ([]*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, postID string) error
	IncrementViewCount(ctx context.Context, postID string) error
}

type PostService struct {
	repo PostRepo
	cld *cloudinary.Cloudinary
}

func NewPostService(repo PostRepo, cld *cloudinary.Cloudinary) *PostService {
	return  &PostService{
		repo: repo,
		cld:  cld,
	}
}

type PaginatedPostsResponse struct {
	TotalPages      int            `json:"totalPages"`
	TotalDocuments  int64          `json:"totalDocuments"`
	Page            int            `json:"page"`
	Limit           int            `json:"limit"`
	Posts          []*model.Post  `json:"posts"`
}

//Get all posts 

//Create POSTS -  Business logic for creating a new post
func (s *PostService) CreatePost(ctx context.Context, req *model.CreatePostRequest, authorID string, fileHeaders []*multipart.FileHeader ) (*model.Post, error) {
	//ValidATE ALL Required fields
	if strings.TrimSpace(req.Title) == "" ||
		strings.TrimSpace(req.Content) == "" ||
		strings.TrimSpace(authorID) == ""{
		return nil , errors.New("all required fields (tittle , content) must be provided")
	}

	//Validate title length
	if len(req.Title) < 3 {
		return nil, errors.New("title must be at least 3 characters long")
	}

	if len(req.Title) > 200 {
		return nil, errors.New("title must not exceed 200 characters")
	}

	// Validate content length
	if len(req.Content) < 10 {
		return nil, errors.New("content must be at least 10 characters long")
	}

	//normalize data
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	// Process tags - split by comma and trim whitespace
	var processedTags []string
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			trimmedTag := strings.TrimSpace(tag)
			if trimmedTag != ""{
				// Convert to lowercase for consistency
				processedTags = append(processedTags, strings.ToLower(trimmedTag))
			}
		}
	}

	// Handle image upload to Cloudinary
	var imageURLs []string 
	if len(fileHeaders) > 0 {
		
		for i, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open uploaded file %d: %w", i, err)
			}
			defer file.Close()

			uploadResult, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
				Folder: "posts",
			})
			if err != nil {
				return nil, fmt.Errorf("failed to upload image %d to Cloudinary: %w", i, err)
			}

			imageURLs = append(imageURLs, uploadResult.SecureURL)
		}
	}	
	// Parse author UUID
	var authorUUID pgtype.UUID
	if err := authorUUID.Scan(authorID); err != nil {
		return nil, fmt.Errorf("invalid author ID format: %w", err)
	}

	// Create post model
	post := &model.Post{
		Title: req.Title,
		Content: req.Content,
		AuthorID: authorUUID,
		ImageURL: imageURLs,
		Tags: processedTags,
	}

	// Save to database
	if err := s.repo.Create(ctx , post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (s *PostService) GetPosts(ctx context.Context, page, limit int)(*PaginatedPostsResponse, error){
	// set default
	if page < 1{
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	totalDocuments, err := s.repo.CountPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	totalPages := int(math.Ceil(float64(totalDocuments) / float64(limit)))

	posts, err := s.repo.GetAllPost(ctx, limit , offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}

	
	return &PaginatedPostsResponse{
		TotalPages: totalPages,
		TotalDocuments: totalDocuments,
		Page: page,
		Limit: limit,
		Posts: posts,
	}, nil

}

func (s *PostService) GetPostByID(ctx context.Context, postID string) (*model.Post, error) {
	// Validate input
	if strings.TrimSpace(postID)  == "" {
		return nil, errors.New("post Id is required")
	}
	
	
	post, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)

	}

	// Check if post exists
	if post == nil {
		return nil, ErrPostNotFound
	}

	// Increment view count
	// _ = s.repo.IncrementViewCount(ctx, postID)
	// Increment view count (async, don't fail if this errors)
	go func() {
		// Use background context since original request may end
		bgCtx := context.Background()
		if err := s.repo.IncrementViewCount(bgCtx, postID); err != nil {
			// Log error but don't fail the request
			log.Printf("Failed to increment view count for post %s: %v", postID, err)
		}
	}()

	return post, nil
}

//GetPostsByAuthorID - Get all posts by author
func (s *PostService) GetPostsByAuthorID(ctx context.Context, authorID string) ([]*model.Post, error) {
	if strings.TrimSpace(authorID) == "" {
		return nil, errors.New("author ID is required")
	}

	posts, err := s.repo.FindByAuthorID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by author: %w", err)
	}

	return posts, nil
}

//update 
func (s *PostService) UpdatePost(ctx context.Context, post *model.Post, userID string) error {
	existing, err := s.repo.FindByID(ctx, post.ID.String())
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPostNotFound
	}

	// Check ownership
	if existing.AuthorID.String() != userID {
		return ErrUnauthorizedPost
	}

	return s.repo.Update(ctx, post)
}

//delete
func (s *PostService) DeletePost(ctx context.Context, postID, userID string) error {
	existing, err := s.repo.FindByID(ctx, postID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPostNotFound
	}

	// Check ownership
	if existing.AuthorID.String() != userID {
		return ErrUnauthorizedPost
	}

	// Delete image from Cloudinary if exists
	// if existing.ImageURL != nil && *existing.ImageURL != "" {
	// 	_ = s.cld.DeleteImage(ctx, *existing.ImageURL)
	// }

	return s.repo.Delete(ctx, postID)
}