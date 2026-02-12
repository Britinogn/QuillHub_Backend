package routes

import (
	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/britinogn/quillhub/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/signup", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	auth.Use(middleware.AdminOnly())
	auth.POST("/admins", authHandler.RegisterAdmin)
}
