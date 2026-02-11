package routes

import (
	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/signup", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
