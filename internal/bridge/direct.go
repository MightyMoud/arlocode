package bridge

import (
	"context"
	"sync/atomic"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type DirectBridge struct {
	agent      *agent.Agent
	events     chan AgentEvent
	cancelFn   context.CancelFunc
	responding atomic.Bool
}

func NewDirectBridge(a *agent.Agent) *DirectBridge {
	db := &DirectBridge{
		agent:  a,
		events: make(chan AgentEvent, 100),
	}
	db.agent.
		WithOnTextChunk(func(s string) {
			db.events <- TextChunkEvent{Text: s}
		}).
		WithOnTextComplete(func() {
			db.events <- TextCompleteEvent{}
		}).
		WithOnThinkingChunk(func(s string) {
			db.events <- ThinkingChunkEvent{Text: s}
		}).
		WithOnThinkingComplete(func() {
			db.events <- ThinkingCompleteEvent{}
		}).
		WithOnToolCall(func(tc tools.ToolCall) {
			db.events <- ToolCallEvent{ToolCall: tc}
		}).
		WithOnTurnComplete(func() {
			db.events <- TurnCompleteEvent{}
		})

	return db
}

func (db *DirectBridge) Run(ctx context.Context, prompt string) error {
	ctx, db.cancelFn = context.WithCancel(ctx)
	db.responding.Store(true)

	go func() {
		err := db.agent.Run(ctx, prompt)
		db.responding.Store(false)
		if err != nil {
			db.events <- ErrorEvent{Err: err}
		}
		db.events <- TurnCompleteEvent{}
	}()

	return nil
}

func (db *DirectBridge) IsResponding() bool {
	return db.responding.Load()
}

func (db *DirectBridge) Events() <-chan AgentEvent {
	return db.events
}

func (db *DirectBridge) Cancel() error {
	if db.cancelFn != nil {
		db.cancelFn()
	}
	return nil
}
func (db *DirectBridge) Close() error {
	close(db.events)
	return nil
}
