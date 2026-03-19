package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	maxTokens      = 1000 // Reduced for faster CPU processing
	analysisPrompt = `Analyze this log and provide:
1. Summary
2. Key errors/warnings
3. Recommendations

Log:
%s`
)

func AnalyzeLogWithAI(logContent string) (string, error) {
	// Get Ollama base URL from environment or use default
	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434/v1"
	}

	// Get model from environment or use default
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.2:latest" // Default Ollama model
	}

	// Truncate log if too long (optimized for CPU performance)
	maxChars := 3000 // Significantly reduced for CPU performance
	if len(logContent) > maxChars {
		logContent = logContent[:maxChars] + "\n\n... (truncated)"
	}

	// Create Ollama client (Ollama doesn't require API key, but library needs one)
	config := openai.DefaultConfig("ollama")
	config.BaseURL = ollamaURL
	client := openai.NewClientWithConfig(config)

	// Prepare prompt
	prompt := fmt.Sprintf(analysisPrompt, logContent)

	// Create context with timeout (generous timeout for CPU processing)
	// Can be overridden with OLLAMA_TIMEOUT environment variable
	timeoutSeconds := 600 // 10 minutes default for CPU
	if timeoutEnv := os.Getenv("OLLAMA_TIMEOUT"); timeoutEnv != "" {
		if parsed, err := time.ParseDuration(timeoutEnv + "s"); err == nil {
			timeoutSeconds = int(parsed.Seconds())
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// Call Ollama API (using OpenAI-compatible endpoint)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a log analyzer. Be concise.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   maxTokens,
			Temperature: 0.3, // Lower temperature for more focused analysis
		},
	)

	if err != nil {
		return "", fmt.Errorf("Ollama API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from Ollama")
	}

	analysis := strings.TrimSpace(resp.Choices[0].Message.Content)
	return analysis, nil
}
