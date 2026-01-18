package butler

import "github.com/mightymoud/arlocode/internal/butler/tools"

type OnTextChunkFunc func(string)
type OnTextCompleteFunc func()
type OnThinkingChunkFunc func(string)
type OnThinkingCompleteFunc func()
type OnToolCallFunc func(tools.ToolCall)
type OnTurnCompleteFunc func()

type EventHooks struct {
	OnTextChunk        func(string)
	OnTextComplete     func()
	OnThinkingChunk    func(string)
	OnThinkingComplete func()
	OnToolCall         func(tools.ToolCall)
	OnTurnComplete     func()
}
