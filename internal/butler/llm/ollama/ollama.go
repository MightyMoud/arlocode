package ollama_llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
	"github.com/ollama/ollama/api"
)

type OllamaLLM struct {
	ModelID string
	Client  *api.Client
}

func (l *OllamaLLM) Stream(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool, hooks butler.EventHooks) (providers.ProviderResponse, error) {
	ollamaTools := makeOllamaTools(agentTools)
	messages := convertMemoryToOllamaMessages(mem)

	stream := true
	req := &api.ChatRequest{
		Model:    l.ModelID,
		Messages: messages,
		Tools:    ollamaTools,
		Stream:   &stream,
	}

	var currentResponseText strings.Builder
	var toolCalls []tools.ToolCall
	isThinking := false
	var metrics *memory.Metrics

	err := l.Client.Chat(ctx, req, func(resp api.ChatResponse) error {
		metrics = &memory.Metrics{
			PromptTokens:     resp.Metrics.PromptEvalCount,
			CompletionTokens: resp.Metrics.EvalCount,
		}
		// Handle Thinking/Reasoning
		if resp.Message.Thinking != "" {
			isThinking = true
			if hooks.OnThinkingChunk != nil {
				hooks.OnThinkingChunk(resp.Message.Thinking)
			}
		}

		// Handle Content
		if resp.Message.Content != "" {
			if isThinking {
				isThinking = false
				if hooks.OnThinkingComplete != nil {
					hooks.OnThinkingComplete()
				}
			}
			if hooks.OnTextChunk != nil {
				hooks.OnTextChunk(resp.Message.Content)
			}
			currentResponseText.WriteString(resp.Message.Content)
		}

		// Handle Tool Calls
		if len(resp.Message.ToolCalls) > 0 {
			for _, tc := range resp.Message.ToolCalls {
				toolCalls = append(toolCalls, tools.ToolCall{
					ID:           fmt.Sprintf("call_%d", len(toolCalls)),
					FunctionName: tc.Function.Name,
					Arguments:    tc.Function.Arguments.ToMap(),
				})
			}
		}

		return nil
	})

	if err != nil {
		return providers.ProviderResponse{}, err
	}

	if hooks.OnTurnComplete != nil {
		hooks.OnTurnComplete()
	}

	return providers.ProviderResponse{
		Text:      currentResponseText.String(),
		ToolCalls: toolCalls,
		Metrics:   metrics,
	}, nil
}

func (l *OllamaLLM) Generate(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool, hooks butler.EventHooks) error {
	ollamaTools := makeOllamaTools(agentTools)
	messages := convertMemoryToOllamaMessages(mem)

	stream := false
	req := &api.ChatRequest{
		Model:    l.ModelID,
		Messages: messages,
		Tools:    ollamaTools,
		Stream:   &stream,
	}

	err := l.Client.Chat(ctx, req, func(resp api.ChatResponse) error {
		if resp.Message.Content != "" {
			fmt.Print(resp.Message.Content)
		}
		return nil
	})

	return err
}
