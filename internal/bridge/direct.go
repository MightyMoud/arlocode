package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/butler/memory"
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

type atifAgent struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type atifTrajectory struct {
	SchemaVersion string           `json:"schema_version"`
	SessionID     string           `json:"session_id"`
	Agent         atifAgent        `json:"agent"`
	Steps         []atifStep       `json:"steps"`
	FinalMetrics  atifFinalMetrics `json:"final_metrics"`
}

type atifFinalMetrics struct {
	TotalPromptTokens     int                    `json:"total_prompt_tokens,omitempty"`
	TotalCompletionTokens int                    `json:"total_completion_tokens,omitempty"`
	TotalCachedTokens     int                    `json:"total_cached_tokens,omitempty"`
	TotalCostUSD          float64                `json:"total_cost_usd,omitempty"`
	TotalSteps            int                    `json:"total_steps"`
	Extra                 map[string]interface{} `json:"extra,omitempty"`
}

type atifStep struct {
	StepID           int                    `json:"step_id"`
	Timestamp        string                 `json:"timestamp,omitempty"`
	Source           string                 `json:"source"`
	ModelName        string                 `json:"model_name,omitempty"`
	ReasoningEffort  interface{}            `json:"reasoning_effort,omitempty"`
	Message          string                 `json:"message"`
	ReasoningContent string                 `json:"reasoning_content,omitempty"`
	ToolCalls        []memory.ToolCall      `json:"tool_calls,omitempty"`
	Observation      *memory.Observation    `json:"observation,omitempty"`
	Metrics          *memory.Metrics        `json:"metrics,omitempty"`
	Extra            map[string]interface{} `json:"extra,omitempty"`
}

func (db *DirectBridge) ExportATIF(path string) (string, error) {
	steps := db.agent.GetMemory()
	exportSteps := make([]atifStep, 0, len(steps))

	sessionID := fmt.Sprintf("session-%d", time.Now().UnixNano())
	final := atifFinalMetrics{
		TotalSteps: len(steps),
	}

	for _, step := range steps {
		exportStep := atifStep{
			StepID:           step.StepID,
			Timestamp:        step.Timestamp,
			Source:           step.Source,
			ModelName:        step.ModelName,
			ReasoningEffort:  step.ReasoningEffort,
			Message:          step.Message,
			ReasoningContent: step.ReasoningContent,
			ToolCalls:        step.ToolCalls,
			Extra:            step.Extra,
		}

		if hasObservation(step.Observation) {
			obs := step.Observation
			exportStep.Observation = &obs
		}
		if hasMetrics(step.Metrics) {
			m := step.Metrics
			exportStep.Metrics = &m

			final.TotalPromptTokens += m.PromptTokens
			final.TotalCompletionTokens += m.CompletionTokens
			final.TotalCachedTokens += m.CachedTokens
			final.TotalCostUSD += m.CostUSD
		}

		exportSteps = append(exportSteps, exportStep)
	}

	traj := atifTrajectory{
		SchemaVersion: "ATIF-v1.5",
		SessionID:     sessionID,
		Agent: atifAgent{
			Name:    "arlocode",
			Version: "unknown",
		},
		Steps:        exportSteps,
		FinalMetrics: final,
	}

	if path == "" || path == "." {
		path = fmt.Sprintf("atif_trajectory_%s.json", sessionID)
	}

	bytes, err := json.MarshalIndent(traj, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func hasObservation(obs memory.Observation) bool {
	return len(obs.Results) > 0
}

func hasMetrics(m memory.Metrics) bool {
	return m.PromptTokens != 0 ||
		m.CompletionTokens != 0 ||
		m.CachedTokens != 0 ||
		m.CostUSD != 0 ||
		len(m.PromptTokenIDs) > 0 ||
		len(m.CompletionTokenIDs) > 0 ||
		len(m.Logprobs) > 0 ||
		(len(m.Extra) > 0)
}
