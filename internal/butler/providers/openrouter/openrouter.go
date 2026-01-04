package openrouter

import (
	"context"
	"log"
	"os"

	"github.com/iamwavecut/gopenrouter"
	"github.com/mightymoud/arlocode/internal/butler/llm"
	openrouter_llm "github.com/mightymoud/arlocode/internal/butler/llm/openrouter"
)

type OpenRouterProvider struct {
	client *gopenrouter.Client
}

// returns a general api client from that provider
func New(ctx context.Context) *OpenRouterProvider {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is not set")
	}
	client := gopenrouter.NewClient(apiKey)
	return &OpenRouterProvider{
		client: client,
	}
}

func (p *OpenRouterProvider) Model(ctx context.Context, modelID string) llm.LLM {
	return &openrouter_llm.OpenRouterLLM{
		ModelID:           modelID,
		Client:            p.client,
		ParallelToolCalls: nil, // True by default for most models
	}
}
