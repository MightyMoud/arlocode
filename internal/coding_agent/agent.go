package coding_agent

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
)

var ctx = context.Background()
var provider = openrouter.New(ctx)
var model = provider.Model(ctx, "anthropic/claude-sonnet-4.5")

var Agent = agent.NewAgent(model)
