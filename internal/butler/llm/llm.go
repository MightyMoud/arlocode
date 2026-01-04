package llm

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/providers"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type OnTextChunkFunc func(string)
type OnThinkingChunkFunc func(string)
type OnToolCallFunc func(tools.ToolCall)

type LLM interface {
	Stream(ctx context.Context, memory []memory.MemoryEntry, tools []tools.Tool) (providers.ProviderResponse, error)
	Generate(ctx context.Context, memory []memory.MemoryEntry, tools []tools.Tool) error
	WithOnThinkingChunk(f OnThinkingChunkFunc) LLM
	WithOnTextChunk(f OnTextChunkFunc) LLM
	WithOnToolCall(f OnToolCallFunc) LLM
}
