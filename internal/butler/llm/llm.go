package llm

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type LLM interface {
	Stream(ctx context.Context, memory []memory.MemoryEntry, tools []tools.Tool, hooks butler.EventHooks) (providers.ProviderResponse, error)
	Generate(ctx context.Context, memory []memory.MemoryEntry, tools []tools.Tool, hooks butler.EventHooks) error
}
