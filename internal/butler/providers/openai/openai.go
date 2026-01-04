package openai_provider

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler/llm"
	openai_llm "github.com/mightymoud/arlocode/internal/butler/llm/openai"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIProvider struct {
	client *openai.Client
}

// New returns a general api client from that provider
func New(ctx context.Context, opts ...option.RequestOption) *OpenAIProvider {
	client := openai.NewClient(opts...)
	return &OpenAIProvider{
		client: &client,
	}
}

// WithApiKey returns config to get the APIkey from certain place
func WithApiKey(key string) option.RequestOption {
	return option.WithAPIKey(key)
}

// Model returns an llm that can generate and stream
func (p *OpenAIProvider) Model(ctx context.Context, modelID string) llm.LLM {
	return &openai_llm.OpenAILLM{ModelID: modelID, Client: p.client}
}
