package providers

import (
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
)

type ProviderResponse struct {
	Text      string
	ToolCalls []tools.ToolCall
}
