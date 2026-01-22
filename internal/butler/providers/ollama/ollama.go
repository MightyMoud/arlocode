package ollama

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/mightymoud/arlocode/internal/butler/llm"
	ollama_llm "github.com/mightymoud/arlocode/internal/butler/llm/ollama"
	"github.com/ollama/ollama/api"
)

type OllamaProvider struct {
	client *api.Client
}

// New returns a general api client from that provider
func New(ctx context.Context) *OllamaProvider {
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}

	u, err := url.Parse(ollamaHost)
	if err != nil {
		// Fallback to default if host is invalid
		u, _ = url.Parse("http://localhost:11434")
	}

	client := api.NewClient(u, http.DefaultClient)
	return &OllamaProvider{
		client: client,
	}
}

func (p *OllamaProvider) Model(ctx context.Context, modelID string) llm.LLM {
	return &ollama_llm.OllamaLLM{
		ModelID: modelID,
		Client:  p.client,
	}
}
