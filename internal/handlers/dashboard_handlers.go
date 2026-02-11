// internal/handlers/dashboard_handler.go
package handlers

import (
	"log"
	"net/http"

	"github.com/britinogn/quillhub/internal/services"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GetAdminDashboard - GET /api/admin/dashboard
func (h *DashboardHandler) GetAdminDashboard(c *gin.Context) {
	log.Printf("[DASHBOARD-HANDLER] Admin dashboard requested")

	ctx := c.Request.Context()
	dashboard, err := h.dashboardService.GetAdminDashboard(ctx)
	if err != nil {
		log.Printf("[DASHBOARD-HANDLER] Error fetching admin dashboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin dashboard fetched successfully",
		"data":    dashboard,
	})
}

// GetUserDashboard - GET /api/dashboard
func (h *DashboardHandler) GetUserDashboard(c *gin.Context) {
	// Get authenticated user ID from middleware
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("[DASHBOARD-HANDLER] User dashboard requested for: %s", userId.(string))

	ctx := c.Request.Context()
	dashboard, err := h.dashboardService.GetUserDashboard(ctx, userId.(string))
	if err != nil {
		log.Printf("[DASHBOARD-HANDLER] Error fetching user dashboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User dashboard fetched successfully",
		"data":    dashboard,
	})
}