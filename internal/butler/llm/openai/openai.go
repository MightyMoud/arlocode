package openai_llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mightymoud/arlocode/internal/butler/llm"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
	"github.com/openai/openai-go/v3"
)

type OpenAILLM struct {
	ModelID         string
	Client          *openai.Client
	OnThinkingChunk llm.OnThinkingChunkFunc
	OnTextChunk     llm.OnTextChunkFunc
	OnToolCall      llm.OnToolCallFunc
}

func (l OpenAILLM) Stream(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool) (providers.ProviderResponse, error) {
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
				if l.OnTextChunk != nil {
					l.OnTextChunk(delta.Content)
				}
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
				if l.OnToolCall != nil {
					l.OnToolCall(tools.ToolCall{
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
	}, nil
}

func (l OpenAILLM) Generate(ctx context.Context, mem []memory.MemoryEntry, agentTools []tools.Tool) error {
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

func (l OpenAILLM) WithOnThinkingChunk(f llm.OnThinkingChunkFunc) llm.LLM {
	l.OnThinkingChunk = f
	return l
}

func (l OpenAILLM) WithOnTextChunk(f llm.OnTextChunkFunc) llm.LLM {
	l.OnTextChunk = f
	return l
}

func (l OpenAILLM) WithOnToolCall(f llm.OnToolCallFunc) llm.LLM {
	l.OnToolCall = f
	return l
}
