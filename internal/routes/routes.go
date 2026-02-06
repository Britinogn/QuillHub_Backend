// internal/routes/routes.go
package routes

import (
	"net/http"

	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/britinogn/quillhub/internal/middleware" 
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	//   API Versioning + Grouping
	api := router.Group("/api")

	//   Public routes (no authentication required)
	public := api.Group("")
	{
		// Health / ping endpoint (very useful)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"version": "1.0.0", 
			})
		})

		// Auth endpoints
		auth := public.Group("/auth")
		{
			auth.POST("/signup", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// Future additions:
			// auth.POST("/forgot-password", authHandler.ForgotPassword)
			// auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Root welcome message (optional)
		public.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "QuillHub API is running 🪶",
			})
		})
	}

	// ────────────────────────────────────────────────
	//   Authenticated / Protected routes
	// ────────────────────────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware()) 
	{
		// User-related (profile, settings, etc.)
		// user := protected.Group("/users")
		// {
		// 	user.GET("/me", authHandler.GetCurrentUser) // you'll add this later
		// 	// user.PUT("/me", authHandler.UpdateProfile)
		// 	// user.DELETE("/me", authHandler.DeleteAccount)
		// }

		// Events / posts / quills (whatever your main feature is)
		//events := protected.Group("/posts") // or /posts, /quills, etc.
		// {
		// 	events.GET("", getEvents)                // list all
		// 	events.GET("/:id", getEventById)         // single event

		// 	events.POST("", createEvent)             // create new
		// 	events.PUT("/:id", updateEvent)          // edit
		// 	events.DELETE("/:id", deleteEvent)       // delete

		// 	// Participation
		// 	events.POST("/:id/register", registerForEvent)
		// 	events.DELETE("/:id/register", cancelRegistration)
		// }

		// Future groups:
		// protected.Group("/comments")
		// protected.Group("/likes")
		// protected.Group("/search")
	}
}




// In your routes.go
// protected := api.Group("")
// protected.Use(middleware.AuthMiddleware())
// {
//     // Admin-only routes
//     admin := protected.Group("/admin")
//     admin.Use(middleware.AdminOnly()) // You'll need to create this middleware
//     {
//         admin.POST("/register-admin", authHandler.RegisterAdmin)
//     }
//}