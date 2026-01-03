package gemini

import (
	"context"
	"testing"
)

func TestWithApiKey(t *testing.T) {
	config := WithApiKey("test-key")
	if config.APIKey != "test-key" {
		t.Errorf("expected APIKey test-key, got %s", config.APIKey)
	}
}

func TestGeminiProvider_Model(t *testing.T) {
	p := &GeminiProvider{}
	llm := p.Model(context.Background(), "gemini-pro")
	if llm == nil {
		t.Fatal("expected LLM, got nil")
	}
}
