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
	// API Versioning + Grouping
	api := router.Group("/api")

	// ────────────────────────────────────────────────
	// Public routes (no authentication required)
	// ────────────────────────────────────────────────
	public := api.Group("")
	{
		// Health check
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"version": "1.0.0",
			})
		})

		// Root welcome
		public.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "QuillHub API is running 🪶",
			})
		})

		// Auth endpoints
		auth := public.Group("/auth")
		{
			auth.POST("/signup", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public post endpoints (anyone can view)
		posts := public.Group("/posts")
		{
			posts.GET("", postHandler.GetAllPosts)                     // List all posts
			posts.GET("/:id", postHandler.GetPostById)                 // Single post
			posts.GET("/author/:authorId", postHandler.GetPostsByAuthorID) // Posts by author

			// Public comments (anyone can view)
			posts.GET("/:id/comments", commentHandler.GetCommentsByPostID) 
			
		}

		comments := public.Group("/comments")
		{
			comments.GET("/:id", commentHandler.GetCommentsByPostID)
		}
	}

	// ────────────────────────────────────────────────
	// Protected routes (authentication required)
	// ────────────────────────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Protected post actions (need to be logged in)
		posts := protected.Group("/posts")
		{
			posts.POST("", postHandler.CreatePost)   // Create post
			posts.PUT("/:id", postHandler.Update)    // Update post
			posts.DELETE("/:id", postHandler.Delete) // Delete post

			// Protected comment actions
			posts.POST("/:postId/comments", commentHandler.CreateComment)
		}

		comments := protected.Group("/comments")
		{
			comments.DELETE("/:commentId", commentHandler.DeleteComment)
		}

		// User profile (authenticated only)
		// user := protected.Group("/users")
		// {
		//     user.GET("/me", authHandler.GetCurrentUser)
		//     user.PUT("/me", authHandler.UpdateProfile)
		// }

		// Future protected routes:
		// protected.Group("/comments")
		// protected.Group("/likes")
	}
}