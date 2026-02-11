package routes

import (
	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPostRoutes(
	public *gin.RouterGroup,
	protected *gin.RouterGroup,
	postHandler *handlers.PostHandler,
	commentHandler *handlers.CommentHandler,
) {

	// Public
	publicPosts := public.Group("/posts")
	{
		publicPosts.GET("", postHandler.GetAllPosts)
		publicPosts.GET("/:id", postHandler.GetPostById)
		publicPosts.GET("/author/:authorId", postHandler.GetPostsByAuthorID)
		publicPosts.GET("/:id/comments", commentHandler.GetCommentsByPostID)
	}

	// Protected
	protectedPosts := protected.Group("/posts")
	{
		protectedPosts.POST("", postHandler.CreatePost)
		protectedPosts.PUT("/:id", postHandler.Update)
		protectedPosts.DELETE("/:id", postHandler.Delete)

		protectedPosts.POST("/:postId/comments", commentHandler.CreateComment)
	}
}
