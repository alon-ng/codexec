package ai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
)

type Client struct {
	client openai.Client
	model  shared.ChatModel
}

func NewClient(cfg Config) *Client {
	clientOptions := []option.RequestOption{
		option.WithAPIKey(cfg.APIKey),
	}

	if cfg.BaseURL != "https://api.openai.com/v1" {
		clientOptions = append(clientOptions, option.WithBaseURL(cfg.BaseURL))
	}

	client := openai.NewClient(clientOptions...)

	return &Client{
		client: client,
		model:  shared.ChatModel(cfg.Model),
	}
}

type SendMessageResponse struct {
	Content          string `json:"content"`
	PromptTokens     int64  `json:"prompt_tokens"`
	CompletionTokens int64  `json:"completion_tokens"`
}

func (c *Client) SendMessage(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion) (*SendMessageResponse, error) {
	chatReq := openai.ChatCompletionNewParams{
		Model:    c.model,
		Messages: messages,
	}

	chatResp, err := c.client.Chat.Completions.New(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	// Extract text content from the first choice
	choice := chatResp.Choices[0]
	content := choice.Message.Content

	if content == "" {
		return nil, fmt.Errorf("empty content in response")
	}

	return &SendMessageResponse{
		Content:          content,
		PromptTokens:     chatResp.Usage.PromptTokens,
		CompletionTokens: chatResp.Usage.CompletionTokens,
	}, nil
}
