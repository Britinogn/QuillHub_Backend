package routes

import (
	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(
	public *gin.RouterGroup,
	protected *gin.RouterGroup,
	commentHandler *handlers.CommentHandler,
) {

	// Public
	publicComments := public.Group("/comments")
	{
		publicComments.GET("/:id", commentHandler.GetCommentsByPostID)
	}

	// Protected
	protectedComments := protected.Group("/comments")
	{
		protectedComments.DELETE("/:commentId", commentHandler.DeleteComment)
	}
}
