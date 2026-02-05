package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc{
	return func(c *gin.Context) {
		// Get user role from context (you'll set this after fetching user from DB)
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "User role not found",
			})
			return 
		}



		role := userRole.(string)

		//Check if user's role is in allowed roles 
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Insufficient permissions",
			})
			return
		}

		c.Next()
	}
}