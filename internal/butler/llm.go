package butler

import (
	"context"

	"github.com/mightymoud/sidekick-agent/internal/butler/common"
)

type LLM interface {
	Stream(ctx context.Context, memory []common.MemoryEntry) error
	Generate(ctx context.Context, memory []common.MemoryEntry) error
}
