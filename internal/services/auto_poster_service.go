// internal/services/auto_poster_service.go
package services

import (
	"context"
	"log"
	"time"

	"github.com/britinogn/quillhub/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type AutoPosterService struct {
	aiService   *AIService
	postRepo    PostRepo
	botUserID   string
	ticker      *time.Ticker
	stopChan    chan bool
}

func NewAutoPosterService(aiService *AIService, postRepo PostRepo, botUserID string) *AutoPosterService {
	return &AutoPosterService{
		aiService: aiService,
		postRepo:  postRepo,
		botUserID: botUserID,
		stopChan:  make(chan bool),
	}
}

// Start - Start the auto-posting scheduler (every 25 minutes)
func (s *AutoPosterService) Start() {
	log.Println("[AUTO-POSTER] 🤖 Starting auto-poster service (posts every 25 minutes)")

	// Create ticker for 25 minutes
	s.ticker = time.NewTicker(25 * time.Minute)

	// Post immediately on start (optional - remove if you don't want this)
	go s.createAndPostBlog()

	// Start the scheduler
	go func() {
		for {
			select {
			case <-s.ticker.C:
				go s.createAndPostBlog()
			case <-s.stopChan:
				log.Println("[AUTO-POSTER] ⏹️ Stopping auto-poster service")
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop - Stop the auto-posting scheduler
func (s *AutoPosterService) Stop() {
	s.stopChan <- true
}

// createAndPostBlog - Generate and publish a blog post
func (s *AutoPosterService) createAndPostBlog() {
	ctx := context.Background()

	log.Println("[AUTO-POSTER] 📝 Generating new AI blog post...")

	// Generate random topic and category
	topic := s.aiService.GenerateRandomTopic()
	category := s.aiService.GenerateRandomCategory()

	log.Printf("[AUTO-POSTER] 💡 Topic: %s | Category: %s", topic, category)

	// Generate blog post using AI
	generatedPost, err := s.aiService.GenerateBlogPost(ctx, topic)
	if err != nil {
		log.Printf("[AUTO-POSTER] ❌ Failed to generate post: %v", err)
		return
	}

	// Parse bot user UUID
	var botUserUUID pgtype.UUID
	if err := botUserUUID.Scan(s.botUserID); err != nil {
		log.Printf("[AUTO-POSTER] ❌ Invalid bot user ID: %v", err)
		return
	}

	// Create the post
	post := &model.Post{
		Title:       generatedPost.Title,
		Content:     generatedPost.Content,
		AuthorID:    botUserUUID,
		Tags:        generatedPost.Tags,
		Category:    &category,
		IsPublished: true,
		ImageURL:    []string{},
	}

	// Save to database
	if err := s.postRepo.Create(ctx, post); err != nil {
		log.Printf("[AUTO-POSTER] ❌ Failed to save post: %v", err)
		return
	}

	log.Printf("[AUTO-POSTER] ✅ Successfully posted: '%s' (ID: %s)", post.Title, post.ID.String())
}