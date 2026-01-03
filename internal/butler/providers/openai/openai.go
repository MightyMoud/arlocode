package openai_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/providers"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
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

type OpenAILLM struct {
	modelID string
	client  *openai.Client
}

// Model returns an llm that can generate and stream
func (p *OpenAIProvider) Model(ctx context.Context, modelID string) butler.LLM {
	return &OpenAILLM{modelID: modelID, client: p.client}
}

func (l OpenAILLM) Stream(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool) (providers.ProviderResponse, error) {
	openaiTools := makeOpenAITools(agentTools)
	messages := convertMemoryToOpenAIMessages(mem)

	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    l.modelID,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
	}
	if len(openaiTools) > 0 {
		params.Tools = openaiTools
	}

	stream := l.client.Chat.Completions.NewStreaming(ctx, params)

	var fullText strings.Builder

	type partialToolCall struct {
		ID   string
		Name string
		Args strings.Builder
	}
	pendingToolCalls := make(map[int64]*partialToolCall)

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta
			if delta.Content != "" {
				fmt.Print(delta.Content)
				fullText.WriteString(delta.Content)
			}

			for _, tc := range delta.ToolCalls {
				index := tc.Index
				if _, ok := pendingToolCalls[index]; !ok {
					pendingToolCalls[index] = &partialToolCall{}
				}
				ptc := pendingToolCalls[index]

				if tc.ID != "" {
					ptc.ID = tc.ID
				}
				if tc.Function.Name != "" {
					ptc.Name = tc.Function.Name
				}
				if tc.Function.Arguments != "" {
					ptc.Args.WriteString(tc.Function.Arguments)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		fmt.Print(err)
		return providers.ProviderResponse{}, err
	}

	var toolCalls []tools.ToolCall
	for _, ptc := range pendingToolCalls {
		var args map[string]any
		if ptc.Args.Len() > 0 {
			if err := json.Unmarshal([]byte(ptc.Args.String()), &args); err != nil {
				log.Printf("Error unmarshaling tool arguments: %v", err)
			}
		}

		toolCalls = append(toolCalls, tools.ToolCall{
			ID:           ptc.ID,
			FunctionName: ptc.Name,
			Arguments:    args,
		})
	}

	return providers.ProviderResponse{
		Text:      fullText.String(),
		ToolCalls: toolCalls,
	}, nil
}

func (l OpenAILLM) Generate(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool) error {
	openaiTools := makeOpenAITools(agentTools)
	messages := convertMemoryToOpenAIMessages(mem)

	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    l.modelID,
	}
	if len(openaiTools) > 0 {
		params.Tools = openaiTools
	}

	resp, err := l.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return err
	}

	if len(resp.Choices) > 0 {
		fmt.Print(resp.Choices[0].Message.Content)
	}
	return nil
}
