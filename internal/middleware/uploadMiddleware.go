package middleware

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ValidateUpload validates file uploads
func ValidateUpload(maxSizeMB int64, allowedExtensions []string) gin.HandlerFunc {
	// normalize allowed extensions once
	allowed := make(map[string]struct{})
	for _, ext := range allowedExtensions {
		allowed[strings.ToLower(ext)] = struct{}{}
	}

	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file uploaded",
			})
			return
		}

		// size check
		maxSizeBytes := maxSizeMB * 1024 * 1024
		if file.Size > maxSizeBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"message":     "File too large",
				"max_size_mb": maxSizeMB,
			})
			return
		}

		// extension check
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if _, ok := allowed[ext]; !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message":       "Invalid file type",
				"allowed_types": allowedExtensions,
			})
			return
		}

		// MIME type check (extra safety)
		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Invalid file content type",
			})
			return
		}

		// pass file to handler
		c.Set("uploadedFile", file)
		c.Next()
	}
}
