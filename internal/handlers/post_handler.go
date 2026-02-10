package handlers

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
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

	// Get uploaded files
	form, _ := c.MultipartForm()
	var files []*multipart.FileHeader
	if form != nil && form.File["images"] != nil {
		files = form.File["images"]
	}

	// Call service
	ctx := c.Request.Context()
	post, err := h.postService.CreatePost(ctx, &req, userId.(string), files)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


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

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	// Parse pagination params with validation
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Call service
	response, err := h.postService.GetPosts(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func(h *PostHandler) GetPostById(c *gin.Context){
	postID := c.Param("id")

	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Post ID is required",
		})
		return
	}

	// Call service to get post
	ctx := c.Request.Context()
	post, err := h.postService.GetPostByID(ctx, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if post exists
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Post not found",
		})
		return
	}

	// Return post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// GetPostsByAuthorID - HTTP handler for GET /posts/author/:authorId
func (h *PostHandler) GetPostsByAuthorID(c *gin.Context) {
	// Get author ID from URL parameter
	authorID := c.Param("authorId")

	if authorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author ID is required"})
		return
	}

	// Call service to get posts
	ctx := c.Request.Context()
	posts, err := h.postService.GetPostsByAuthorID(ctx, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return posts (empty array if no posts found)
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"count": len(posts),
	})
}

func (h *PostHandler) Update(c *gin.Context) {
	// Get post ID from URL
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	// Get authenticated user ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req model.UpdatePostRequest
	
	// Check content type
	contentType := c.GetHeader("Content-Type")
	
	// Handle based on content type
	if strings.Contains(contentType, "application/json") {
		// JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON", "details": err.Error()})
			return
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Form-data request
		if title := c.PostForm("title"); title != "" {
			req.Title = &title
		}
		if content := c.PostForm("content"); content != "" {
			req.Content = &content
		}
		if category := c.PostForm("category"); category != "" {
			req.Category = &category
		}
		if tagsString := c.PostForm("tags"); tagsString != "" {
			tags := []string{}
			for _, tag := range strings.Split(tagsString, ",") {
				tags = append(tags, strings.TrimSpace(tag))
			}
			req.Tags = tags
		}
		if isPublished := c.PostForm("is_published"); isPublished != "" {
			published := isPublished == "true"
			req.IsPublished = &published
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	log.Printf("[POST-HANDLER] Updating post %s by user: %s", postID, userId.(string))

	// Get uploaded files (only for form-data)
	form, _ := c.MultipartForm()
	var files []*multipart.FileHeader
	if form != nil && form.File["images"] != nil {
		files = form.File["images"]
		log.Printf("[POST-HANDLER] Received %d new image files", len(files))
	}

	// Call service
	ctx := c.Request.Context()
	post, err := h.postService.UpdatePost(ctx, &req, postID, userId.(string), files)
	if err != nil {
		log.Printf("[POST-HANDLER] Update error: %v", err)
		
		if errors.Is(err, services.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if errors.Is(err, services.ErrUnauthorizedPost) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this post"})
			return
		}
		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[POST-HANDLER] Post updated successfully: %s", postID)

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"data": model.PostResponse{
			ID:        post.ID.String(),
			AuthorID:  post.AuthorID.String(),
			Title:     post.Title,
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			Category:  post.Category,
			Tags:      post.Tags,
			IsPublished: post.IsPublished,
			ViewCount: post.ViewCount,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		},
	})
}


func (h *PostHandler) Delete(c *gin.Context) {
	// Get post ID from URL
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	// Get authenticated user ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("[POST-HANDLER] Deleting post %s by user: %s", postID, userId.(string))

	// Call service
	ctx := c.Request.Context()
	err := h.postService.DeletePost(ctx, postID, userId.(string))
	if err != nil {
		log.Printf("[POST-HANDLER] Delete error: %v", err)
		
		if errors.Is(err, services.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if errors.Is(err, services.ErrUnauthorizedPost) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this post"})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	log.Printf("[POST-HANDLER] Post deleted successfully: %s", postID)

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}