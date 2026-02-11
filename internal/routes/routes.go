package routes

import (
	"net/http"

	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/britinogn/quillhub/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	postHandler *handlers.PostHandler,
	commentHandler *handlers.CommentHandler,
) {

	api := router.Group("/api")

	// Health
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "QuillHub API is running",
		})
	})

	// Public
	public := api.Group("")

	// Protected
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	// Register separated routes
	RegisterAuthRoutes(public, authHandler)
	RegisterPostRoutes(public, protected, postHandler, commentHandler)
	RegisterCommentRoutes(public, protected, commentHandler)
}
