package butler

import "github.com/mightymoud/arlocode/internal/butler/tools"

type OnTextChunkFunc func(string)
type OnThinkingChunkFunc func(string)
type OnToolCallFunc func(tools.ToolCall)

type EventHooks struct {
	OnTextChunk     func(string)
	OnThinkingChunk func(string)
	OnToolCall      func(tools.ToolCall)
}
