// internal/model/ai_post.go
package model

type AIPostRequest struct {
	Topic      string   `json:"topic"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	MinLength  int      `json:"min_length"`  // Minimum word count
	MaxLength  int      `json:"max_length"`  // Maximum word count
}

type AIGeneratedPost struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}