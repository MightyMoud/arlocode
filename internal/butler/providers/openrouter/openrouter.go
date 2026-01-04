package openrouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/iamwavecut/gopenrouter"
	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
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

type OpenRouterLLM struct {
	modelID           string
	client            *gopenrouter.Client
	parallelToolCalls *bool // unused atm
}

func (l OpenRouterLLM) Stream(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool) (providers.ProviderResponse, error) {
	openRouterTools := makeOpenRouterTools(agentTools)
	messages := convertMemoryToOpenRouterMessages(mem)

	req := gopenrouter.ChatCompletionRequest{
		Model:    l.modelID,
		Messages: messages,
		// Tools must be included in every request (initial and follow-ups)
		// per OpenRouter documentation
		Tools: openRouterTools,
		Reasoning: &gopenrouter.ReasoningParams{
			MaxTokens: 1000,
		},
	}

	// Apply parallel tool calls configuration if set
	if l.parallelToolCalls != nil {
		req.ParallelToolCalls = l.parallelToolCalls
	}

	stream, err := l.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return providers.ProviderResponse{}, err
	}
	defer stream.Close()

	var currentResponseText strings.Builder

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

			// Handle Content
			if delta.Content != "" {
				fmt.Printf("%s", delta.Content)
				currentResponseText.WriteString(delta.Content)
			}
			if delta.Reasoning != "" {
				// fmt.Printf("%s", delta.Reasoning)
				// color.Magenta("%s", delta.Reasoning)
				color.RGB(255, 128, 0).Printf("%s", delta.Reasoning)
				// currentResponseText.WriteString(delta.Reasoning)
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
	}

	return providers.ProviderResponse{
		Text:      currentResponseText.String(),
		ToolCalls: toolCalls,
	}, nil
}

func (l OpenRouterLLM) Generate(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool) error {
	openRouterTools := makeOpenRouterTools(agentTools)
	messages := convertMemoryToOpenRouterMessages(mem)

	req := gopenrouter.ChatCompletionRequest{
		Model:    l.modelID,
		Messages: messages,
		// Tools must be included in every request
		Tools: openRouterTools,
	}

	// Apply parallel tool calls configuration if set
	if l.parallelToolCalls != nil {
		req.ParallelToolCalls = l.parallelToolCalls
	}

	resp, err := l.client.CreateChatCompletion(ctx, req)
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

func (p *OpenRouterProvider) Model(ctx context.Context, modelID string) butler.LLM {
	return &OpenRouterLLM{
		modelID:           modelID,
		client:            p.client,
		parallelToolCalls: nil, // True by default for most models
	}
}

// WithParallelToolCalls sets whether multiple tools can be called simultaneously
func (l *OpenRouterLLM) WithParallelToolCalls(enabled bool) *OpenRouterLLM {
	l.parallelToolCalls = &enabled
	return l
}
