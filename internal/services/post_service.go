package services

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostRepo interface{
	Create(ctx context.Context, post *model.Post) error 
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

//Get all posts 

//Create POSTS -  Business logic for creating a new post
func (s *PostService) CreatePost(ctx context.Context, req *model.CreatePostRequest, authorID string, fileHeader *multipart.FileHeader ) (*model.Post, error) {
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
	var imageURL *string
	if fileHeader != nil {
		//open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}
		defer file.Close()

		// Upload to Cloudinary
		uploadResult , err := s.cld.Upload.Upload(ctx, file , uploader.UploadParams{
			Folder: "posts",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to upload image to Cloudinary: %w", err)
		}

		imageURL = &uploadResult.SecureURL
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
		ImageURL: imageURL,
		Tags: processedTags,
	}

	// Save to database
	if err := s.repo.Create(ctx , post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}
