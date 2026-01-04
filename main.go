package main

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
)

func main() {
	ctx := context.Background()

	prompt := "Some tool tests are missing in the tools package. Please add them."

	openrouterProvider := openrouter.New(ctx)
	model := openrouterProvider.Model(ctx, "anthropic/claude-sonnet-4.5")
	openrouterBasedAgent := butler.NewAgent(model).WithMaxIterations(50)
	openrouterBasedAgent.Run(ctx, prompt)

	// geminiProvider := gemini.New(ctx)
	// model2 := geminiProvider.Model(ctx, "gemini-3-flash-preview")
	// geminiBasedAgent := butler.NewAgent(model2)
	// geminiBasedAgent.Run(ctx, prompt)
}
