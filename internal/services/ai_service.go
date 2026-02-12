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

		// MORE

		"Weird food combinations people actually love",
		"Hidden gems in your city you probably never noticed",
		"The psychology behind why we procrastinate (and how to stop)",
		"Best cheap date ideas that actually work",
		"Why adults still love cartoons and nostalgic games",
		"The most underrated travel destinations in 2025",
		"How to travel the world on a tight budget",
		"Scary travel stories that actually happened",
		"Beautiful small towns you should visit before they become tourist traps",
		"The world's strangest museums and why they're worth seeing",
		"Dark humor jokes that are too good to be true",
		"Dad jokes so bad they're actually good",
		"Memes that perfectly describe adult life in 2025",
		"The funniest autocorrect fails of all time",
		"Why Gen Z humor is completely different from Millennials",
		"Mind-blowing historical facts most people don't know",
		"How the human brain actually learns new things",
		"The real story behind everyday inventions",
		"Why some people are naturally good at languages",
		"The science of why music makes us feel emotions",
		"Countries with the weirdest laws still in effect",
		"The world's happiest countries — and why they’re happy",
		"Secret traditions only locals know about",
		"Countries where time moves differently (literally)",
		"The most polite and most rude countries according to travelers",
		"Why we yawn — and why it’s contagious",
		"The weirdest things doctors have found inside people",
		"How your body changes when you fall in love",
		"Foods that are secretly good for your brain",
		"Why some people never get sick (genetics or luck?)",
		"Animals with the strangest superpowers",
		"Why cats are secretly plotting world domination",
		"The psychology of why we love true crime",
		"Conspiracy theories that turned out to be true",
		"The most satisfying optical illusions ever",

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

		// New related categories
		"Microservices",
		"API Design",
		"System Design",
		"Performance Optimization",
		"Scalability",
		"Distributed Systems",
		"Containerization",
		"Orchestration",
		"Infrastructure as Code",
		"CI/CD Pipelines",
		"Observability & Monitoring",
		"Logging & Tracing",
		"Testing Strategies",

		// New related categories
		"Fun & Weird Facts",
		"Travel & Hidden Places",
		"Jokes & Humor",
		"Psychology & Mind",
		"Countries & Cultures",
		"Health & Body Science",
		"Animals & Nature",
		"Life Hacks & Tips",
		"Memes & Internet Culture",
		"History & True Stories",
		"Love & Relationships",
		"Food & Eating Habits",
		"Daily Life & Adulthood",
		"Conspiracy & Mysteries",
		"Optical Illusions & Brain Tricks",
		"Education & Learning Hacks",
		"Budget & Money Saving",
		"Strange Traditions",
		"Dark Humor & Edgy Jokes",
		"Personal Development",
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
	prompt := fmt.Sprintf(`You are a skilled blogger who adapts tone to the topic.

Topic: "%s"

Rules:
- Create a catchy, engaging title
- Write 150–350 words of well-structured content
- If the topic is technical/programming → use professional tone, clear explanations, code examples in markdown blocks
- If the topic is fun/travel/jokes/psychology/countries/health/life → use light, relatable, humorous tone
- Structure: short intro, main body (sections or bullets), quick conclusion
- End with a question or call-to-action to engage readers
- Use markdown for formatting (headings, bold, code blocks, lists)
- No first person ("I", "we") — neutral voice

Return ONLY valid JSON in this exact format, nothing else:
{
	"title": "Your title here",
	"content": "Full post text here (markdown ok)",
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