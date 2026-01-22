package providers

import (
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type ProviderResponse struct {
	Text      string
	ToolCalls []tools.ToolCall
	Metrics   *memory.Metrics
}
