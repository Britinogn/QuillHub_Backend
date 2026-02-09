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
	var req model.CreatePostRequest
	
	// Check content type
	contentType := c.GetHeader("Content-Type")
	log.Printf("[POST-HANDLER] Content-Type: %s", contentType)
	
	// Handle based on content type
	if strings.Contains(contentType, "application/json") {
		// JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
			return
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Form-data request
		title := c.PostForm("title")
		content := c.PostForm("content")
		tagsString := c.PostForm("tags")
		
		log.Printf("[POST-HANDLER] Form data - title: '%s', content: '%s', tags: '%s'", title, content, tagsString)
		
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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	// Get authenticated user ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("[POST-HANDLER] Creating post - Title: '%s', Content: '%s'", req.Title, req.Content)

	// Get uploaded files
	form, _ := c.MultipartForm()
	var files []*multipart.FileHeader
	if form != nil && form.File["images"] != nil {
		files = form.File["images"]
		log.Printf("[POST-HANDLER] Received %d files", len(files))
	}

	// Call service
	ctx := c.Request.Context()
	post, err := h.postService.CreatePost(ctx, &req, userId.(string), files)
	if err != nil {
		log.Printf("[POST-HANDLER] Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[POST-HANDLER] Post created: %s", post.ID.String())

	// Return response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"data": model.PostResponse{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Title:     post.Title,
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			Tags:      post.Tags,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
	})
}