package coding_agent

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
)

var ctx = context.Background()

var provider = openrouter.New(ctx)

// var provider = openai_provider.New(ctx,
// 	openai_provider.WithApiKey(os.Getenv("ZAI_API_KEY")),
// 	openai_provider.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
// )

var model = provider.Model(ctx, "z-ai/glm-4.7")

var Agent = agent.NewAgent(model)
