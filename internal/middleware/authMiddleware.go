package middleware

import (
	"net/http"
	"strings"

	"github.com/britinogn/quillhub/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc{
	return func (c *gin.Context)  {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return 
		}

		//Check if its a bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
			return 
		}

		token := parts[1]
		
		//verify token - this returns *Claims, not string
		claims, err := utils.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return 
		}
	
		// Store the USER ID STRING from claims, not the whole claims object
		c.Set("userId", claims.UserID)  // ← Extract UserID from claims
		c.Set("userRole", claims.Role)
		
		// Optionally store other useful info
		// c.Set("userEmail", claims.Email)
		// c.Set("userRole", claims.Role)
		
		// Continue to next handler
		c.Next()
	}
}