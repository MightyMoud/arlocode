package main

import (
	"context"
	"os"

	"github.com/mightymoud/sidekick-agent/internal/butler"
	gemini "github.com/mightymoud/sidekick-agent/internal/butler/providers/google"
)

func readFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func main() {
	ctx := context.Background()
	// client, err := genai.NewClient(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Define the tool
	// tools := []*genai.Tool{
	// 	{
	// 		FunctionDeclarations: []*genai.FunctionDeclaration{
	// 			{
	// 				Name:        "readFile",
	// 				Description: "Read the contents of a file",
	// 				Parameters: &genai.Schema{
	// 					Type: genai.TypeObject,
	// 					Properties: map[string]*genai.Schema{
	// 						"path": {
	// 							Type:        genai.TypeString,
	// 							Description: "The path to the file to read",
	// 						},
	// 					},
	// 					Required: []string{"path"},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// config := &genai.GenerateContentConfig{
	// 	Tools: tools,
	// 	ThinkingConfig: &genai.ThinkingConfig{
	// 		IncludeThoughts: true,
	// 	},
	// }

	prompt := "Read the file 'go.mod' and tell me how what version is used for go in this project"

	googleProvider := gemini.New(ctx)
	model := googleProvider.Model(ctx, "gemini-2.5-flash")
	runnableAgent := butler.NewAgent(model)
	runnableAgent.Run(ctx, prompt)

}
