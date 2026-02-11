// internal/handlers/ai_handler.go
package handlers

// import (
// 	"log"
// 	"net/http"

// 	"github.com/britinogn/quillhub/internal/services"
// 	"github.com/gin-gonic/gin"
// )

// type AIHandler struct {
// 	autoPoster *services.AutoPosterService
// }

// func NewAIHandler(autoPoster *services.AutoPosterService) *AIHandler {
// 	return &AIHandler{autoPoster: autoPoster}
// }

// // TriggerAIPost - Manually trigger AI post generation (admin only)
// func (h *AIHandler) TriggerAIPost(c *gin.Context) {
// 	log.Printf("[AI-HANDLER] Manual AI post generation triggered")
	
// 	// Trigger post creation
// 	go h.autoPoster.createAndPostBlog()
	
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "AI post generation triggered successfully",
// 	})
// }