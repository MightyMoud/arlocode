package memory

import "github.com/mightymoud/sidekick-agent/internal/butler/tools"

type MemoryEntry struct {
	Message    string
	Role       string
	ToolName   string
	ToolCallID string
	ToolCalls  []tools.ToolCall
}
