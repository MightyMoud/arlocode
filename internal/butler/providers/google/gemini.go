package gemini

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	"github.com/mightymoud/sidekick-agent/internal/butler/common"
	"google.golang.org/genai"
)

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

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

func (l GeminiLLM) Stream(ctx context.Context, memory []common.MemoryEntry) error {
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
		},
	}
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
	fmt.Println(history)

	for {
		resp := l.client.Models.GenerateContentStream(ctx, l.modelID, history, config)

		var currentResponseParts []*genai.Part
		var functionCalls []*genai.FunctionCall

		for chunk, err := range resp {
			if err != nil {
				log.Fatal(err)
			}

			for _, part := range chunk.Candidates[0].Content.Parts {
				currentResponseParts = append(currentResponseParts, part)

				if part.FunctionCall != nil {
					// Ignore function calls for now
				} else if part.Thought {
					fmt.Printf("\n[Thinking]: %s", part.Text)
				} else {
					fmt.Printf("%s", part.Text)
				}
			}
		}
		modelResponseHistory := &genai.Content{
			Role:  "model",
			Parts: currentResponseParts,
		}
		history = append(history, modelResponseHistory)

		if len(functionCalls) == 0 {
			// No more functions to call, we are done!
			break
		}

		for _, fn := range functionCalls {
			fmt.Printf("Calling tool: %s\n", fn.Name)

			var content string
			if fn.Name == "readFile" {
				path, ok := fn.Args["path"].(string)
				if !ok {
					content = "Error: invalid argument 'path'"
				} else {
					c, err := readFile(path)
					if err != nil {
						content = fmt.Sprintf("Error: %v", err)
					} else {
						content = c
					}
				}
			} else {
				content = fmt.Sprintf("Error: unknown tool %s", fn.Name)
			}
			functionResponseHistory := &genai.Content{
				Role: "user",
				Parts: []*genai.Part{{
					FunctionResponse: &genai.FunctionResponse{
						Name:     fn.Name,
						Response: map[string]any{"content": content},
					},
				}},
			}
			history = append(history, functionResponseHistory)
		}
	}
	return nil
}

func (l GeminiLLM) Generate(ctx context.Context, memory []common.MemoryEntry) error {
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
