package main

import (
	"context"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	"github.com/mightymoud/sidekick-agent/internal/butler/providers/openrouter"
)

func main() {
	ctx := context.Background()

	prompt := "Update the function listFolderContents and readFolderContentFn to ignore a .git folder"

	openrouterProvider := openrouter.New(ctx)
	model := openrouterProvider.Model(ctx, "x-ai/grok-code-fast-1")
	openrouterBasedAgent := butler.NewAgent(model)
	openrouterBasedAgent.Run(ctx, prompt)

	// geminiProvider := gemini.New(ctx)
	// model2 := geminiProvider.Model(ctx, "gemini-3-flash-preview")
	// geminiBasedAgent := butler.NewAgent(model2)
	// geminiBasedAgent.Run(ctx, prompt)
}
