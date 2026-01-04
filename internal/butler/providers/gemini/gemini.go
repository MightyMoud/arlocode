package gemini

import (
	"context"
	"log"

	"github.com/mightymoud/arlocode/internal/butler/llm"
	gemini_llm "github.com/mightymoud/arlocode/internal/butler/llm/gemini"
	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
}

// returns a general api client from that provider
func New(ctx context.Context) *GeminiProvider {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &GeminiProvider{
		client: client,
	}
}

// returns config to get the APIkey from certain place
func WithApiKey(key string) *genai.ClientConfig {
	return &genai.ClientConfig{
		APIKey: key,
	}
}

// returns an llm that can generate and stream
func (p *GeminiProvider) Model(ctx context.Context, modelID string) llm.LLM {
	return &gemini_llm.GeminiLLM{ModelID: modelID, Client: p.client}
}
