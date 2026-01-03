package main

import (
	"context"
	"os"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	"github.com/mightymoud/sidekick-agent/internal/butler/providers/openrouter"
)

func main() {
	ctx := context.Background()

	prompt := "Read the file 'gemini.go' in folder of internal/butler/providers/gemini/ and tell me what is it missing?"

	openrouterProvider := openrouter.New(os.Getenv("OPENROUTER_API_KEY"))
	model := openrouterProvider.Model(ctx, "x-ai/grok-code-fast-1")
	runnableAgent := butler.NewAgent(model)
	runnableAgent.Run(ctx, prompt)
}
