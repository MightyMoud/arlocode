package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/mightymoud/arlocode/internal/butler"
	"github.com/mightymoud/arlocode/internal/butler/llm"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type Agent struct {
	llm                llm.LLM
	memory             []memory.MemoryEntry
	tools              []tools.Tool
	maxIterations      int
	ModelName          string
	OnTextChunk        butler.OnTextChunkFunc
	OnTextComplete     butler.OnTextCompleteFunc
	OnThinkingChunk    butler.OnThinkingChunkFunc
	OnThinkingComplete butler.OnThinkingCompleteFunc
	OnToolCall         butler.OnToolCallFunc
	OnTurnComplete     butler.OnTurnCompleteFunc
}

func NewAgent(l llm.LLM) *Agent {
	return &Agent{
		llm:           l,
		memory:        []memory.MemoryEntry{},
		tools:         tools.StdToolset,
		maxIterations: 10, // Default max iterations as recommended by OpenRouter docs
		ModelName:     "unknown",
	}
}

// Mods
func (a *Agent) WithMemory(memory []memory.MemoryEntry) *Agent {
	a.memory = memory
	return a
}

func (a *Agent) WithModelName(modelName string) *Agent {
	a.ModelName = modelName
	return a
}

func (a *Agent) WitTools(tools []tools.Tool) *Agent {
	a.tools = tools
	return a
}

func (a *Agent) WithNoTools() *Agent {
	a.tools = []tools.Tool{}
	return a
}

func (a *Agent) WithMaxIterations(max int) *Agent {
	a.maxIterations = max
	return a
}

func (l *Agent) WithOnThinkingChunk(f butler.OnThinkingChunkFunc) *Agent {
	l.OnThinkingChunk = f
	return l
}

func (l *Agent) WithOnThinkingComplete(f butler.OnThinkingCompleteFunc) *Agent {
	l.OnThinkingComplete = f
	return l
}

func (l *Agent) WithOnTextChunk(f butler.OnTextChunkFunc) *Agent {
	l.OnTextChunk = f
	return l
}

func (l *Agent) WithOnTextComplete(f butler.OnTextCompleteFunc) *Agent {
	l.OnTextComplete = f
	return l
}

func (l *Agent) WithOnToolCall(f butler.OnToolCallFunc) *Agent {
	l.OnToolCall = f
	return l
}

func (l *Agent) WithOnTurnComplete(f butler.OnTurnCompleteFunc) *Agent {
	l.OnTurnComplete = f
	return l
}

// Mock for Memory stuff later this is where Agent will use it
func (a *Agent) AddMemoryEntry(entry memory.MemoryEntry) {
	a.memory = append(a.memory, entry)
}

// Memory stuff later
func (a *Agent) GetMemory() []memory.MemoryEntry {
	return a.memory
}

func (a *Agent) HandleToolCall(ctx context.Context, call tools.ToolCall) (string, error) {
	var tool tools.Tool
	for _, t := range a.tools {
		if t.Name == call.FunctionName {
			tool = t
			break
		}
	}

	argsPtr := reflect.New(tool.ArgType).Interface()

	// Convert the map[string]any into the concrete struct.
	// This works for ANY provider because we go through JSON bytes first.
	bytes, _ := json.Marshal(call.Arguments)
	if err := json.Unmarshal(bytes, argsPtr); err != nil {
		return "", fmt.Errorf("failed to unmarshal tool args: %w", err)
	}

	results := tool.Handler.Call([]reflect.Value{
		reflect.ValueOf(argsPtr).Elem(),
	})

	if len(results) > 1 && !results[1].IsNil() {
		return "", results[1].Interface().(error)
	}
	resultStr := results[0].String()

	if a.OnToolCall != nil {
		a.OnToolCall(tools.ToolCall{
			ID:           "toolID", // TODO: generate unique ID for tool call
			FunctionName: tool.Name,
			Arguments:    call.Arguments,
		})
	}
	// Maybe useful to debug later
	// if len(resultStr) > 100 {
	// 	color.Blue("Tool %s returned: %s", call.FunctionName, resultStr[:100])
	// } else {
	// 	color.Blue("Tool %s returned: %s", call.FunctionName, resultStr)
	// }

	return resultStr, nil
}

func (a *Agent) Run(ctx context.Context, prompt string) error {
	initMessage := memory.MemoryEntry{
		StepID:  a.nextStepID(),
		Source:  "user",
		Message: prompt,
	}
	a.AddMemoryEntry(initMessage)

	iterationCount := 0
	for iterationCount < a.maxIterations {
		iterationCount++

		var reasoning strings.Builder
		hooks := butler.EventHooks{
			OnTextChunk:    a.OnTextChunk,
			OnTextComplete: a.OnTextComplete,
			OnToolCall:     a.OnToolCall,
			OnTurnComplete: a.OnTurnComplete,
			OnThinkingChunk: func(chunk string) {
				reasoning.WriteString(chunk)
				if a.OnThinkingChunk != nil {
					a.OnThinkingChunk(chunk)
				}
			},
			OnThinkingComplete: func() {
				if a.OnThinkingComplete != nil {
					a.OnThinkingComplete()
				}
			},
		}

		result, err := a.llm.Stream(ctx, a.memory, a.tools, hooks)
		if err != nil {
			log.Fatal("Error calling LLM Stream: ", err)
			return err
		}

		entry := memory.MemoryEntry{
			StepID:  a.nextStepID(),
			Source:  "agent",
			Message: result.Text,
		}
		if reasoning.Len() > 0 {
			entry.ReasoningContent = reasoning.String()
		}
		if result.Metrics != nil {
			entry.Metrics = *result.Metrics
		}
		if len(result.ToolCalls) > 0 {
			entry.ToolCalls = mapToolCalls(result.ToolCalls)
		}

		if len(result.ToolCalls) > 0 {
			var results []memory.ObservationResult
			for _, call := range result.ToolCalls {
				// Ask user for confirmation before executing tool
				// var confirm bool
				// if err := huh.NewConfirm().
				// 	Title("Execute tool?").
				// 	Description(fmt.Sprintf("Do you want to execute %s?", call.FunctionName)).
				// 	Affirmative("Yes").
				// 	Negative("No").
				// 	Value(&confirm).
				// 	WithTheme(huh.ThemeCatppuccin()).
				// 	Run(); err != nil {
				// 	color.Red("Error getting confirmation: %v", err)
				// 	continue
				// }

				// if !confirm {
				// 	color.Red("Tool call %s was cancelled by the user.", call.FunctionName)
				// 	a.AddMemoryEntry(memory.MemoryEntry{
				// 		Role:       "tool",
				// 		Message:    "Tool call cancelled by user",
				// 		ToolName:   call.FunctionName,
				// 		ToolCallID: call.ID,
				// 	})
				// 	continue
				// }

				output, _ := a.HandleToolCall(ctx, call)
				results = append(results, memory.ObservationResult{
					SourceCallID: call.ID,
					Content:      output,
				})
			}
			entry.Observation = memory.Observation{Results: results}
		}

		a.AddMemoryEntry(entry)

		if len(result.ToolCalls) == 0 {
			break
		}
	}

	if iterationCount >= a.maxIterations {
		color.Yellow("\nWarning: Maximum iterations (%d) reached. The agent loop was terminated.\n", a.maxIterations)
	}

	return nil
}

func (a *Agent) nextStepID() int {
	return len(a.memory) + 1
}

func mapToolCalls(calls []tools.ToolCall) []memory.ToolCall {
	if len(calls) == 0 {
		return nil
	}
	toolCalls := make([]memory.ToolCall, 0, len(calls))
	for _, call := range calls {
		toolCalls = append(toolCalls, memory.ToolCall{
			ToolCallID:   call.ID,
			FunctionName: call.FunctionName,
			Arguments:    call.Arguments,
		})
	}
	return toolCalls
}
