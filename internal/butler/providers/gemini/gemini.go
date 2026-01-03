package gemini

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/providers"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
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

type GeminiLLM struct {
	modelID string
	client  *genai.Client
}

func (l GeminiLLM) Stream(ctx context.Context, memory []memory.MemoryEntry, agentTools []tools.Tool) (providers.ProviderResponse, error) {
	geminiTools := makeGeminiTools(agentTools)

	config := &genai.GenerateContentConfig{
		Tools: geminiTools,
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
		},
	}

	history := convertMemoryToGeminiHistory(memory)
	resp := l.client.Models.GenerateContentStream(ctx, l.modelID, history, config)

	var currentResponseText []string
	var functionCalls []tools.ToolCall

	for chunk, err := range resp {
		if err != nil {
			log.Fatal(err)
		}

		for _, part := range chunk.Candidates[0].Content.Parts {

			if part.FunctionCall != nil {
				functionCalls = append(functionCalls, tools.ToolCall{
					ID:           part.FunctionCall.ID,
					FunctionName: part.FunctionCall.Name,
					Arguments:    part.FunctionCall.Args,
				})
			} else if part.Thought {
				// will be replaced with TUI integration or hooks
				fmt.Printf("\n[Thinking]: %s", part.Text)
			} else {
				// wil be replaced with TUI integration or hooks
				fmt.Printf("%s", part.Text)
				currentResponseText = append(currentResponseText, part.Text)
			}
		}
	}

	var textResponse strings.Builder
	for _, part := range currentResponseText {
		textResponse.WriteString(part)
	}
	return providers.ProviderResponse{
		Text:      textResponse.String(),
		ToolCalls: functionCalls,
	}, nil
}

func (l GeminiLLM) Generate(ctx context.Context, memory []memory.MemoryEntry, tools []tools.Tool) error {
	history := []*genai.Content{}
	for _, entry := range memory {
		genAIEntry := genai.Content{
			Role: entry.Role,
			Parts: []*genai.Part{
				{Text: entry.Message},
			},
		}
		history = append(history, &genAIEntry)
	}
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
		},
	}
	resp, err := l.client.Models.GenerateContent(ctx, l.modelID, history, config)

	fmt.Print(resp.Text())
	return err
}

// returns an llm that can generate and stream
func (p *GeminiProvider) Model(ctx context.Context, modelID string) butler.LLM {
	return &GeminiLLM{modelID: modelID, client: p.client}
}
