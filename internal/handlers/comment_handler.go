package handlers

import (
	"errors"
	"net/http"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/britinogn/quillhub/internal/services"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment - HTTP handler for POST /posts/:postId/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	// Get post ID from URL parameter
	postID := c.Param("postId")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	// Parse request body
	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get authenticated user ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Call service to create comment
	ctx := c.Request.Context()
	comment, err := h.commentService.CreateComment(ctx, &req, postID, userId.(string))
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return success
	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// GetCommentsByPostID - HTTP handler for GET /posts/:postId/comments
func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	// Get post ID from URL parameter
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post ID is required"})
		return
	}

	// Call service to get comments
	comments, err := h.commentService.GetCommentsByPostID(c.Request.Context(), postID)
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return comments
	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"count":    len(comments),
	})
}

// DeleteComment - HTTP handler for DELETE /comments/:commentId
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	// Get comment ID from URL parameter
	commentID := c.Param("commentId")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment ID is required"})
		return
	}

	// Get authenticated user ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Call service to delete comment
	err := h.commentService.DeleteComment(c.Request.Context(), commentID, userId.(string))
	if err != nil {
		// Handle specific errors
		if errors.Is(err, services.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		if errors.Is(err, services.ErrUnauthorizedComment) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}