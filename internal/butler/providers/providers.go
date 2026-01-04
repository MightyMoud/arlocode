package providers

import (
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type ProviderResponse struct {
	Text      string
	ToolCalls []tools.ToolCall
}
