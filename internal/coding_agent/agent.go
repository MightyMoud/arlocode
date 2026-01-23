package coding_agent

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/butler/providers/ollama"
)

var ctx = context.Background()

var provider = ollama.New(ctx)

// var provider = openai_provider.New(ctx,
// 	openai_provider.WithApiKey(os.Getenv("ZAI_API_KEY")),
// 	openai_provider.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
// )

var model = provider.Model(ctx, "qwen3:8b")

var Agent = agent.NewAgent(model).WithModelName("qwen3:14b").WithMaxIterations(300)
