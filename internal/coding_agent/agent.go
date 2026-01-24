package coding_agent

import (
	"context"
	"os"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	openai_provider "github.com/mightymoud/arlocode/internal/butler/providers/openai"
)

var systemPrompt = `You are ArloCode, an AI coding agent designed to assist developers with complex coding tasks.
You have access to various tools to help you accomplish your goals.
You should think carefully about which tools to use and when to use them.
You should also provide reasoning for your actions to help the user understand your thought process.

When you need to use a tool, you must specify the tool name and provide the necessary arguments in JSON format.
After executing a tool, you should analyze the results and decide on the next steps.

Always aim to provide clear and concise explanations for your actions and decisions.
Your ultimate goal is to assist the user in completing their coding tasks effectively and efficiently.
You are in a terminal environment and have access to local resources and files. You are working in a coding project. Your tools run in the same folder as the project.
`

var ctx = context.Background()

// var provider = ollama.New(ctx)

var provider = openai_provider.New(ctx,
	openai_provider.WithApiKey(os.Getenv("ZAI_API_KEY")),
	openai_provider.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
)

var model = provider.Model(ctx, "glm-4.7")

var Agent = agent.NewAgent(model).WithModelName("glm-4.7").WithMaxIterations(300).WithMemory([]memory.MemoryEntry{
	{
		StepID:  1,
		Source:  "system",
		Message: systemPrompt,
	},
})
