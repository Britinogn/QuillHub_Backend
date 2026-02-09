package handlers

import (
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
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get form fields
	title := c.PostForm("title")
	content := c.PostForm("content")
	tagsString := c.PostForm("tags")

	// Split tags by comma
	var tags []string
	if tagsString != "" {
		tags = strings.Split(tagsString, ",")
	}

	// Create request object
	req := &model.CreatePostRequest{
		Title:   title,
		Content: content,
		// ImageURL: ,
		Tags:    tags,
	}

	// Get uploaded file
	var fileHeader *multipart.FileHeader
	if uploadedFile, exists := c.Get("uploadedFile"); exists {
		fileHeader = uploadedFile.(*multipart.FileHeader)
	} else {
		file, err := c.FormFile("file")
		if err == nil {
			fileHeader = file
		}
	}

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Call service
	ctx := c.Request.Context()
	post, err := h.postService.CreatePost(ctx, req, userID.(string), fileHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}