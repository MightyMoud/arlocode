package bridge

import (
	"context"

	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type AgentEvent interface {
	agentEvent()
}

type TextChunkEvent struct{ Text string }
type TextCompleteEvent struct{}
type ThinkingChunkEvent struct{ Text string }
type ThinkingCompleteEvent struct{}
type ToolCallEvent struct{ ToolCall tools.ToolCall }
type ErrorEvent struct{ Err error }
type TurnCompleteEvent struct{}

func (TextChunkEvent) agentEvent()        {}
func (TextCompleteEvent) agentEvent()     {}
func (ThinkingChunkEvent) agentEvent()    {}
func (ThinkingCompleteEvent) agentEvent() {}
func (ToolCallEvent) agentEvent()         {}
func (ErrorEvent) agentEvent()            {}
func (TurnCompleteEvent) agentEvent()     {}

type AgentBridge interface {
	Run(ctx context.Context, prompt string) error
	Events() <-chan AgentEvent
	Cancel() error
	Close() error
	IsResponding() bool
}
