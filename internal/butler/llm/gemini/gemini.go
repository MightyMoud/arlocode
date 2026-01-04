package gemini_llm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/mightymoud/arlocode/internal/butler/llm"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
	"google.golang.org/genai"
)

type GeminiLLM struct {
	ModelID         string
	Client          *genai.Client
	OnTextChunk     llm.OnTextChunkFunc
	OnThinkingChunk llm.OnThinkingChunkFunc
	OnToolCall      llm.OnToolCallFunc
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
	resp := l.Client.Models.GenerateContentStream(ctx, l.ModelID, history, config)

	var currentResponseText []string
	var functionCalls []tools.ToolCall

	for chunk, err := range resp {
		if err != nil {
			log.Fatal(err)
		}

		for _, part := range chunk.Candidates[0].Content.Parts {

			if part.FunctionCall != nil {
				functionCalls = append(functionCalls, tools.ToolCall{
					ID:               part.FunctionCall.ID,
					FunctionName:     part.FunctionCall.Name,
					Arguments:        part.FunctionCall.Args,
					ThoughtSignature: part.ThoughtSignature,
				})
			} else if part.Thought {
				// will be replaced with TUI integration or hooks
				// fmt.Printf("\n[Thinking]: %s", part.Text)
				color.RGB(230, 42, 42).Printf("%s", part.Text)
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
	resp, err := l.Client.Models.GenerateContent(ctx, l.ModelID, history, config)

	fmt.Print(resp.Text())
	return err
}

func (l GeminiLLM) WithOnThinkingChunk(f llm.OnThinkingChunkFunc) llm.LLM {
	l.OnThinkingChunk = f
	return l
}

func (l GeminiLLM) WithOnTextChunk(f llm.OnTextChunkFunc) llm.LLM {
	l.OnTextChunk = f
	return l
}

func (l GeminiLLM) WithOnToolCall(f llm.OnToolCallFunc) llm.LLM {
	l.OnToolCall = f
	return l
}
