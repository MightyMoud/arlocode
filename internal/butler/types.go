package butler

import "github.com/mightymoud/arlocode/internal/butler/tools"

type OnTextChunkFunc func(string)
type OnStreamCompleteFunc func()
type OnThinkingChunkFunc func(string)
type OnThinkingCompleteFunc func()
type OnToolCallFunc func(tools.ToolCall)

type EventHooks struct {
	OnTextChunk        func(string)
	OnStreamComplete   func()
	OnThinkingChunk    func(string)
	OnThinkingComplete func()
	OnToolCall         func(tools.ToolCall)
}
