package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get role from context (you must set it earlier)
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "User role not found",
			})
			return
		}

		role, ok := userRole.(string)
		if !ok || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			return
		}

		// if !ok || (role != "admin" && role != "moderator") {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		// 	return
		// }

		// Only admins reach here
		c.Next()
	}
}