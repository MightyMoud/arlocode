package main

import (
	"context"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	"github.com/mightymoud/sidekick-agent/internal/butler/providers/gemini"
)

func main() {
	ctx := context.Background()

	prompt := "Read the file 'gemini.go' in folder of internal/butler/providers/gemini/ and tell me what is it missing?"

	geminiProvider := gemini.New(ctx)
	model := geminiProvider.Model(ctx, "gemini-2.5-flash")
	runnableAgent := butler.NewAgent(model)
	runnableAgent.Run(ctx, prompt)
}
