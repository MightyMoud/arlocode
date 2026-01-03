package main

import (
	"context"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	openai_provider "github.com/mightymoud/sidekick-agent/internal/butler/providers/openai"
)

func main() {
	ctx := context.Background()

	prompt := "Read the file 'gemini.go' in folder of internal/butler/providers/gemini/ and tell me what is it missing?"

	openaiProvider := openai_provider.New(ctx)
	model := openaiProvider.Model(ctx, "gpt-5.1")
	runnableAgent := butler.NewAgent(model)
	runnableAgent.Run(ctx, prompt)
}
