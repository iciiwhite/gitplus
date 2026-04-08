package ai

import (
	"context"
	"fmt"

	"github.com/iciwhite/gitplus/internal/config"
	"github.com/sashabaranov/go-openai"
)

type Assistant struct {
	client *openai.Client
	cfg    *config.Config
}

func NewAssistant(cfg *config.Config) *Assistant {
	client := openai.NewClient(cfg.OpenAIKey)
	return &Assistant{
		client: client,
		cfg:    cfg,
	}
}

func (a *Assistant) SuggestCommitMessage(diff string) (string, error) {
	prompt := fmt.Sprintf("Write a concise git commit message for this diff:\n%s", diff)
	resp, err := a.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful developer assistant that writes clear git commit messages.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 60,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (a *Assistant) GeneratePRDescription(title, body string) (string, error) {
	prompt := fmt.Sprintf("Write a clear pull request description for a PR titled '%s'. Additional context: %s", title, body)
	resp, err := a.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful developer assistant that writes concise pull request descriptions.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 300,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}