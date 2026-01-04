package main

import (
	"context"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	"github.com/mightymoud/sidekick-agent/internal/butler/providers/openrouter"
)

func main() {
	ctx := context.Background()

	prompt := "read this page https://openrouter.ai/docs/guides/features/tool-calling and then read the codebase and check if you can improve the way I call the tools in openrouter provider implementation"

	openrouterProvider := openrouter.New(ctx)
	model := openrouterProvider.Model(ctx, "anthropic/claude-sonnet-4.5")
	openrouterBasedAgent := butler.NewAgent(model)
	openrouterBasedAgent.Run(ctx, prompt)

	// geminiProvider := gemini.New(ctx)
	// model2 := geminiProvider.Model(ctx, "gemini-3-flash-preview")
	// geminiBasedAgent := butler.NewAgent(model2)
	// geminiBasedAgent.Run(ctx, prompt)
}
