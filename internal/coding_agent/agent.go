package coding_agent

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
)

var ctx = context.Background()
var provider = openrouter.New(ctx)
var model = provider.Model(ctx, "z-ai/glm-4.7")

var Agent = agent.NewAgent(model)
