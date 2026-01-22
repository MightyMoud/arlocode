package coding_agent

import (
	"context"
	"os"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	openai_provider "github.com/mightymoud/arlocode/internal/butler/providers/openai"
)

var ctx = context.Background()

// var provider = openrouter.New(ctx)
var provider = openai_provider.New(ctx,
	openai_provider.WithApiKey(os.Getenv("ZAI_API_KEY")),
	openai_provider.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
)
var model = provider.Model(ctx, "glm-4.7-flash")

var Agent = agent.NewAgent(model)
