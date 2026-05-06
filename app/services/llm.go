package services

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/anthropics/anthropic-sdk-go"
)

// LLMClient supports multiple LLM providers
type LLMClient interface {
	Chat(ctx context.Context, messages []ChatMessage) (string, error)
	Complete(ctx context.Context, prompt string) (string, error)
}

// ChatMessage for LLM input
type ChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"` // message content
}

// LLM Provider types
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderLocal   Provider = "local"
)

// NewLLMClient creates an LLM client based on configuration
func NewLLMClient() LLMClient {
	provider := Provider(os.Getenv("LLM_PROVIDER"))
	apiKey := os.Getenv("OPENAI_API_KEY"))

	switch provider {
	case ProviderAnthropic:
		return NewAnthropicClient(os.Getenv("ANTHROPIC_API_KEY"))
	case ProviderLocal:
		return NewLocalClient(os.Getenv("LOCAL_LLM_URL"))
	default:
		return NewOpenAIClient(apiKey)
	}
}

// OpenAI Client
type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
	}
}

func (c *OpenAIClient) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	msgs := make([]openai.ChatMessage, len(messages))
	for i, m := range messages {
		msgs[i] = openai.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4,
		Messages: msgs,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) Complete(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateCompletion(ctx, openai.CompletionRequest{
		Model:       openai.GPT4,
		Prompt:      prompt,
		MaxTokens:   2000,
		Temperature: 0.7,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}

// Anthropic Client
type AnthropicClient struct {
	client *anthropic.Client
}

func NewAnthropicClient(apiKey string) *AnthropicClient {
	return &AnthropicClient{
		client: anthropic.NewClient(apiKey),
	}
}

func (c *AnthropicClient) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	var system string
	var userMsgs []string
	for _, m := range messages {
		if m.Role == "system" {
			system = m.Content
		} else {
			userMsgs = append(userMsgs, m.Content)
		}
	}

	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:      anthropic.ModelClaude3_5Sonnet20241022,
		MaxTokens:  anthropic.Int(2000),
		System:     []anthropic.TextBlock{{Text: system}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(userMsgs[len(userMsgs)-1]),
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Content[0].Text, nil
}

func (c *AnthropicClient) Complete(ctx context.Context, prompt string) (string, error) {
	return c.Chat(ctx, []ChatMessage{
		{Role: "user", Content: prompt},
	})
}

// Local/Other LLM Client (Ollama, LM Studio, etc.)
type LocalClient struct {
	url string
}

func NewLocalClient(url string) *LocalClient {
	return &LocalClient{
		url: url,
	}
}

func (c *LocalClient) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	// Implement local LLM API call
	return "", fmt.Errorf("not implemented")
}

func (c *LocalClient) Complete(ctx context.Context, prompt string) (string, error) {
	// Implement local LLM API call
	return "", fmt.Errorf("not implemented")
}