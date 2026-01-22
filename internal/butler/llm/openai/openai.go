package openai_llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
	"github.com/openai/openai-go/v3"
)

type OpenAILLM struct {
	ModelID string
	Client  *openai.Client
}

func (l OpenAILLM) Stream(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool, hooks butler.EventHooks) (providers.ProviderResponse, error) {
	openaiTools := makeOpenAITools(agentTools)
	messages := convertMemoryToOpenAIMessages(mem)

	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    l.ModelID,
		StreamOptions: openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: openai.Bool(true),
		},
	}
	if len(openaiTools) > 0 {
		params.Tools = openaiTools
	}

	stream := l.Client.Chat.Completions.NewStreaming(ctx, params)

	var fullText strings.Builder
	var metrics *memory.Metrics

	type partialToolCall struct {
		ID   string
		Name string
		Args strings.Builder
	}
	pendingToolCalls := make(map[int64]*partialToolCall)
	var isThinking bool

	for stream.Next() {
		chunk := stream.Current()
		if chunk.Usage.PromptTokens > 0 || chunk.Usage.CompletionTokens > 0 {
			metrics = &memory.Metrics{
				PromptTokens:     int(chunk.Usage.PromptTokens),
				CompletionTokens: int(chunk.Usage.CompletionTokens),
			}
		}
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta

			if delta.JSON.ExtraFields != nil {
				for _, key := range []string{"reasoning_content", "reasoning"} {
					if field, ok := delta.JSON.ExtraFields[key]; ok {
						raw := field.Raw()
						if raw != "" && raw != "null" {
							var reasoning string
							if err := json.Unmarshal([]byte(raw), &reasoning); err == nil && reasoning != "" {
								if hooks.OnThinkingChunk != nil {
									hooks.OnThinkingChunk(reasoning)
								}
								isThinking = true
								break
							}
						}
					}
				}
			}

			if delta.Content != "" {
				if isThinking {
					if hooks.OnThinkingComplete != nil {
						hooks.OnThinkingComplete()
					}
					isThinking = false
				}
				if hooks.OnTextChunk != nil {
					hooks.OnTextChunk(delta.Content)
				}
				fullText.WriteString(delta.Content)
			}

			if len(delta.ToolCalls) > 0 && isThinking {
				if hooks.OnThinkingComplete != nil {
					hooks.OnThinkingComplete()
				}
				isThinking = false
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
				if hooks.OnToolCall != nil {
					hooks.OnToolCall(tools.ToolCall{
						ID:           ptc.ID,
						FunctionName: ptc.Name,
						Arguments:    nil, // partial, will be filled later
					})
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
		Metrics:   metrics,
	}, nil
}

func (l OpenAILLM) Generate(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool, hooks butler.EventHooks) error {
	openaiTools := makeOpenAITools(agentTools)
	messages := convertMemoryToOpenAIMessages(mem)

	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    l.ModelID,
	}
	if len(openaiTools) > 0 {
		params.Tools = openaiTools
	}

	resp, err := l.Client.Chat.Completions.New(ctx, params)
	if err != nil {
		return err
	}

	if len(resp.Choices) > 0 {
		fmt.Print(resp.Choices[0].Message.Content)
	}
	return nil
}
