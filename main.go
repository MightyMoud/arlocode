package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

func main() {
	ctx := context.Background()

	// prompt := "Some tool tests are missing in the tools package. Please add them."
	// prompt := "Add a README.md file to the internal package named butler that will explain what the butler package is doing and how to use it."
	// prompt := "Add a github action workflow file that will run tests for all the packages on every push to the repository."
	prompt := "Read all the code in the butler package and tell me what it does:"

	openrouterProvider := openrouter.New(ctx)
	model := openrouterProvider.Model(ctx, "z-ai/glm-4.7").
		WithOnThinkingChunk(func(chunk string) {
			color.RGB(255, 128, 0).Printf("%s", chunk)

		}).
		WithOnTextChunk(func(chunk string) {
			fmt.Printf("%s", chunk)

		}).
		WithOnToolCall(func(t tools.ToolCall) {
			color.Blue("\n[Tool Call] %s - with Arguments: %+v", t.FunctionName, t.Arguments)
		})
	openrouterBasedAgent := butler.NewAgent(model).WithMaxIterations(50)
	openrouterBasedAgent.Run(ctx, prompt)

	// geminiProvider := gemini.New(ctx)
	// model2 := geminiProvider.Model(ctx, "gemini-3-flash-preview")
	// geminiBasedAgent := butler.NewAgent(model2)
	// geminiBasedAgent.Run(ctx, prompt)
}
