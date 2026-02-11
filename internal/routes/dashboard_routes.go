package routes

import (
	"github.com/britinogn/quillhub/internal/handlers"
	"github.com/britinogn/quillhub/internal/middleware"
	"github.com/gin-gonic/gin"
)

// func RegisterDashboardRoutes(
// 	protected *gin.RouterGroup, 
// 	dashboardHandler *handlers.DashboardHandler,
// 	) {
// 	protectedDashboard := protected.Group("/dashboard")
// 	{
// 		protectedDashboard.GET("/admin", dashboardHandler.GetAdminDashboard)
// 		protectedDashboard.GET("/user", dashboardHandler.GetUserDashboard)
// 	}
// }


func RegisterDashboardRoutes(protected *gin.RouterGroup, dashboardHandler *handlers.DashboardHandler) {
	// Admin dashboard — strict admin only
	admin := protected.Group("/dashboard/admin")
	admin.Use(middleware.AdminOnly())
	admin.GET("", dashboardHandler.GetAdminDashboard)

	// User dashboard — any authenticated user
	protected.GET("/dashboard/user", dashboardHandler.GetUserDashboard)
}