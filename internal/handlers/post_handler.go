package handlers

import (
	"log"
	"mime/multipart" 
	"net/http"
	"strings"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/britinogn/quillhub/internal/services"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePost - HTTP handler for POST /posts
func (h *PostHandler) CreatePost(c *gin.Context) {
	// Try JSON first
	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON fails, try form-data
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		
		// Parse form-data
		title := c.PostForm("title")
		content := c.PostForm("content")
		tagsString := c.PostForm("tags")
		
		var tags []string
		if tagsString != "" {
			for _, tag := range strings.Split(tagsString, ",") {
				tags = append(tags, strings.TrimSpace(tag))
			}
		}
		
		req = model.CreatePostRequest{
			Title:   title,
			Content: content,
			Tags:    tags,
		}
	}

	// Get authenticated user ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("[POST-HANDLER] Creating post by user: %s", userId.(string))

	// Get uploaded files (support multiple images)
	form, _ := c.MultipartForm()
	var files []*multipart.FileHeader
	if form != nil && form.File["images"] != nil {
		files = form.File["images"] // ← Changed to "images" (plural) to accept multiple
	}

	// Call service
	ctx := c.Request.Context()
	post, err := h.postService.CreatePost(ctx, &req, userId.(string), files)
	if err != nil {
		log.Printf("[POST-HANDLER] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[POST-HANDLER] Post created successfully: %s", post.ID.String())

	// Return response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"data": model.PostResponse{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Title:     post.Title,
			Content:   post.Content,
			ImageURL:  post.ImageURL, // ← Now array
			Tags:      post.Tags,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
	})
}