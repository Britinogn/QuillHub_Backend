// internal/services/ai_service.go
package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/britinogn/quillhub/internal/model"
	"cloud.google.com/go/ai/generativelanguage/apiv1beta/generativelanguagepb"
	"google.golang.org/api/option"
	generativelanguage "cloud.google.com/go/ai/generativelanguage/apiv1beta"
)

type AIService struct {
	apiKey string
	client *generativelanguage.GenerativeClient
}

func NewAIService() *AIService {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	
	client, err := generativelanguage.NewGenerativeClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("[AI-SERVICE] Failed to create Gemini client: %v", err)
		return &AIService{apiKey: apiKey}
	}

	return &AIService{
		apiKey: apiKey,
		client: client,
	}
}

// Close - Clean up the client
func (s *AIService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
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

// GenerateBlogPost - Generate a blog post using Gemini AI SDK
func (s *AIService) GenerateBlogPost(ctx context.Context, topic string) (*model.AIGeneratedPost, error) {
	if s.apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY not set")
	}

	if s.client == nil {
		return nil, errors.New("Gemini client not initialized")
	}

	log.Printf("[AI-SERVICE] Generating blog post for topic: %s", topic)

	// Create prompt
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

	// Call Gemini API using SDK
	req := &generativelanguagepb.GenerateContentRequest{
		Model: "models/gemini-flash-latest", // Free tier model
		Contents: []*generativelanguagepb.Content{
			{
				Parts: []*generativelanguagepb.Part{
					{
						Data: &generativelanguagepb.Part_Text{
							Text: prompt,
						},
					},
				},
			},
		},
	}

	resp, err := s.client.GenerateContent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract response text
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("empty response from Gemini API")
	}

	responseText := resp.Candidates[0].Content.Parts[0].GetText()
	
	// Clean markdown code fences if present
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")
	responseText = strings.TrimSpace(responseText)

	// Parse JSON response
	var generatedPost model.AIGeneratedPost
	if err := json.Unmarshal([]byte(responseText), &generatedPost); err != nil {
		log.Printf("[AI-SERVICE] Failed to parse JSON. Raw response: %s", responseText)
		return nil, fmt.Errorf("failed to parse generated post: %w", err)
	}

	log.Printf("[AI-SERVICE] Successfully generated post: %s", generatedPost.Title)
	return &generatedPost, nil
}