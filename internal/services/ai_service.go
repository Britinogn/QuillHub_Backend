// internal/services/ai_service.go
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/britinogn/quillhub/internal/model"
)

type AIService struct {
	apiKey     string
	httpClient *http.Client
}

func NewAIService() *AIService {
	return &AIService{
		apiKey: os.Getenv("ANTHROPIC_API_KEY"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ClaudeRequest - Request structure for Claude API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse - Response structure from Claude API
type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// GenerateRandomTopic - Generate a random tech topic
func (s *AIService) GenerateRandomTopic() string {
	topics := []string{
		"Go programming best practices",
		"Building scalable microservices",
		"Docker and containerization",
		"Kubernetes deployment strategies",
		"RESTful API design patterns",
		"PostgreSQL performance optimization",
		"JWT authentication implementation",
		"GraphQL vs REST APIs",
		"CI/CD pipeline setup",
		"Cloud architecture patterns",
		"Git workflow strategies",
		"Test-driven development",
		"Database indexing strategies",
		"Caching strategies with Redis",
		"Message queues and async processing",
		"gRPC vs HTTP/REST",
		"Serverless architecture",
		"API rate limiting techniques",
		"Monitoring and observability",
		"Security best practices in web development",
	}

	rand.Seed(time.Now().UnixNano())
	return topics[rand.Intn(len(topics))]
}

// GenerateRandomCategory - Generate a random category
func (s *AIService) GenerateRandomCategory() string {
	categories := []string{
		"Backend Development",
		"DevOps",
		"Architecture",
		"Security",
		"Databases",
		"Cloud Computing",
		"Best Practices",
	}

	rand.Seed(time.Now().UnixNano())
	return categories[rand.Intn(len(categories))]
}

// GenerateBlogPost - Generate a blog post using Claude AI
func (s *AIService) GenerateBlogPost(ctx context.Context, topic string) (*model.AIGeneratedPost, error) {
	if s.apiKey == "" {
		return nil, errors.New("ANTHROPIC_API_KEY not set")
	}

	log.Printf("[AI-SERVICE] Generating blog post for topic: %s", topic)

	// Create prompt for Claude
	prompt := fmt.Sprintf(`You are a technical blog writer. Write a complete, professional blog post about: "%s"

Requirements:
- Write a compelling, SEO-friendly title
- Write 800-1200 words of well-structured content
- Include an introduction, main body with sections, and conclusion
- Use clear, technical language
- Include code examples if relevant (use markdown code blocks)
- Make it engaging and informative
- Don't use first person ("I", "we")

Format your response EXACTLY as JSON with this structure:
{
  "title": "Your catchy title here",
  "content": "Your full blog post content here (with markdown formatting)",
  "tags": ["tag1", "tag2", "tag3", "tag4", "tag5"]
}

Only return valid JSON, nothing else.`, topic)

	// Prepare Claude API request
	reqBody := ClaudeRequest{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4000,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make API request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Claude API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return nil, errors.New("empty response from Claude API")
	}

	// Extract JSON from Claude's response
	responseText := claudeResp.Content[0].Text
	
	// Parse the generated post
	var generatedPost model.AIGeneratedPost
	if err := json.Unmarshal([]byte(responseText), &generatedPost); err != nil {
		return nil, fmt.Errorf("failed to parse generated post: %w", err)
	}

	log.Printf("[AI-SERVICE] Successfully generated post: %s", generatedPost.Title)
	return &generatedPost, nil
}