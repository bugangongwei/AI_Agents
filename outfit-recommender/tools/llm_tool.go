package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func GetLLMRecommendation(prompt string) (recommendation string, err error) {
	err = godotenv.Load("outfit-recommender/.env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY not set")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.deepseek.com/v1"
	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "deepseek-chat",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		recommendation = resp.Choices[0].Message.Content
	}

	return recommendation, nil
}
