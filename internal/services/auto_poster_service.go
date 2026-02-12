// internal/services/auto_poster_service.go
package services

import (
	"context"
	"log"
	"sync"
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
	mu          sync.Mutex  // ✅ Prevent concurrent posting
	isRunning   bool        // ✅ Track running state
}

func NewAutoPosterService(aiService *AIService, postRepo PostRepo, botUserID string) *AutoPosterService {
	return &AutoPosterService{
		aiService: aiService,
		postRepo:  postRepo,
		botUserID: botUserID,
		stopChan:  make(chan bool),
		isRunning: false,
	}
}

// Start - Start the auto-posting scheduler (every 25 minutes)
func (s *AutoPosterService) Start() {
	s.mu.Lock()
	if s.isRunning {
		log.Println("[AUTO-POSTER] ⚠️  Service already running")
		s.mu.Unlock()
		return
	}
	s.isRunning = true
	s.mu.Unlock()

	log.Println("[AUTO-POSTER] 🤖 Starting auto-poster service (posts every 25 minutes)")

	// Create ticker for 25 minutes
	s.ticker = time.NewTicker(25 * time.Minute)

	// ✅ Post immediately on start (optional - comment out if not needed)
	go s.createAndPostBlog()

	// Start the scheduler
	go func() {
		for {
			select {
			case <-s.ticker.C:
				log.Println("[AUTO-POSTER] ⏰ Timer triggered - creating new post")
				go s.createAndPostBlog()
			case <-s.stopChan:
				log.Println("[AUTO-POSTER] ⏹️  Stopping auto-poster service")
				s.ticker.Stop()
				s.mu.Lock()
				s.isRunning = false
				s.mu.Unlock()
				return
			}
		}
	}()

	log.Println("[AUTO-POSTER] ✅ Auto-poster service started successfully")
}

// Stop - Stop the auto-posting scheduler
func (s *AutoPosterService) Stop() {
	s.mu.Lock()
	if !s.isRunning {
		log.Println("[AUTO-POSTER] ⚠️  Service not running")
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	s.stopChan <- true
	log.Println("[AUTO-POSTER] 🛑 Stop signal sent")
}

// IsRunning - Check if auto-poster is currently running
func (s *AutoPosterService) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isRunning
}

// createAndPostBlog - Generate and publish a blog post
func (s *AutoPosterService) createAndPostBlog() {
	// ✅ Use timeout context to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

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
		ImageURL:   []string{}, // ✅ Changed from ImageURL to ImageURLs (match your model)
	}

	// Save to database
	if err := s.postRepo.Create(ctx, post); err != nil {
		log.Printf("[AUTO-POSTER] ❌ Failed to save post: %v", err)
		return
	}

	log.Printf("[AUTO-POSTER] ✅ Successfully posted: '%s' (ID: %s)", 
		post.Title, post.ID.String())
	log.Printf("[AUTO-POSTER] 🏷️  Tags: %v | Category: %s", post.Tags, *post.Category)
}

// PostNow - Manually trigger a post creation (for testing/admin)
func (s *AutoPosterService) PostNow() {
	log.Println("[AUTO-POSTER] 🚀 Manual post creation triggered")
	go s.createAndPostBlog()
}