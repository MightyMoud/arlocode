package openrouter_llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/iamwavecut/gopenrouter"
	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type OpenRouterLLM struct {
	ModelID           string
	Client            *gopenrouter.Client
	ParallelToolCalls *bool // unused atm
}

func (l OpenRouterLLM) Stream(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool, hooks butler.EventHooks) (providers.ProviderResponse, error) {
	openRouterTools := makeOpenRouterTools(agentTools)
	messages := convertMemoryToOpenRouterMessages(mem)

	req := gopenrouter.ChatCompletionRequest{
		Model:    l.ModelID,
		Messages: messages,
		// Tools must be included in every request (initial and follow-ups)
		// per OpenRouter documentation
		Tools: openRouterTools,
		Reasoning: &gopenrouter.ReasoningParams{
			MaxTokens: 1000,
		},
	}

	// Apply parallel tool calls configuration if set
	if l.ParallelToolCalls != nil {
		req.ParallelToolCalls = l.ParallelToolCalls
	}

	stream, err := l.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return providers.ProviderResponse{}, err
	}
	defer stream.Close()

	var currentResponseText strings.Builder
	isThinking := false

	type streamToolCall struct {
		Index int
		ID    string
		Name  string
		Args  strings.Builder
	}
	pendingToolCalls := make(map[int]*streamToolCall)

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("Stream error: %v", err)
			break
		}

		if len(response.Choices) > 0 {
			delta := response.Choices[0].Delta

			if delta.Reasoning != "" {
				isThinking = true
				if hooks.OnThinkingChunk != nil {
					hooks.OnThinkingChunk(delta.Reasoning)
				}
				// currentResponseText.WriteString(delta.Reasoning)
			}

			// Handle Content
			if delta.Content != "" {
				if isThinking {
					isThinking = false
					if hooks.OnThinkingComplete != nil {
						hooks.OnThinkingComplete()
					}
				}
				if hooks.OnTextChunk != nil {
					hooks.OnTextChunk(delta.Content)
				}
				currentResponseText.WriteString(delta.Content)
			}

			// Handle Tool Calls
			if len(delta.ToolCalls) > 0 {
				for _, tc := range delta.ToolCalls {
					idx := tc.Index
					if _, exists := pendingToolCalls[idx]; !exists {
						pendingToolCalls[idx] = &streamToolCall{Index: idx}
					}

					if tc.ID != "" {
						pendingToolCalls[idx].ID = tc.ID
					}
					if tc.Function.Name != "" {
						pendingToolCalls[idx].Name = tc.Function.Name
					}
					if tc.Function.Arguments != "" {
						pendingToolCalls[idx].Args.WriteString(tc.Function.Arguments)
					}
				}
			}
		}

		// Accumulate usage stats -> for later
		// if response.Usage != nil {
		// 	fmt.Printf("\n[Usage]: Prompt: %d, Completion: %d, Total: %d, Cost: %f\n",
		// 		response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens, response.Usage.Cost)
		// }
	}

	var toolCalls []tools.ToolCall
	for _, ptc := range pendingToolCalls {
		var args map[string]any
		if ptc.Args.Len() > 0 {
			// Try to unmarshal, if fails (e.g. partial json), we might have an issue,
			// but at end of stream it should be complete.
			if err := json.Unmarshal([]byte(ptc.Args.String()), &args); err != nil {
				log.Printf("Error unmarshaling tool args: %v", err)
			}
		}
		toolCalls = append(toolCalls, tools.ToolCall{
			ID:           ptc.ID,
			FunctionName: ptc.Name,
			Arguments:    args,
		})
		if hooks.OnToolCall != nil {
			hooks.OnToolCall(tools.ToolCall{
				ID:           ptc.ID,
				FunctionName: ptc.Name,
				Arguments:    args,
			})
		}
	}
	if hooks.OnStreamComplete != nil {
		hooks.OnStreamComplete()
	}

	return providers.ProviderResponse{
		Text:      currentResponseText.String(),
		ToolCalls: toolCalls,
	}, nil
}

func (l OpenRouterLLM) Generate(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool, hooks butler.EventHooks) error {
	openRouterTools := makeOpenRouterTools(agentTools)
	messages := convertMemoryToOpenRouterMessages(mem)

	req := gopenrouter.ChatCompletionRequest{
		Model:    l.ModelID,
		Messages: messages,
		// Tools must be included in every request
		Tools: openRouterTools,
	}

	// Apply parallel tool calls configuration if set
	if l.ParallelToolCalls != nil {
		req.ParallelToolCalls = l.ParallelToolCalls
	}

	resp, err := l.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}

	if len(resp.Choices) > 0 {
		fmt.Print(resp.Choices[0].Message.Content)
		fmt.Printf("\n[Usage]: Prompt: %d, Completion: %d, Total: %d",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}
	return nil
}

// WithParallelToolCalls sets whether multiple tools can be called simultaneously
func (l *OpenRouterLLM) WithParallelToolCalls(enabled bool) *OpenRouterLLM {
	l.ParallelToolCalls = &enabled
	return l
}
