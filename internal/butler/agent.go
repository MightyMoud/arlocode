package butler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/fatih/color"
	"github.com/mightymoud/arlocode/internal/butler/llm"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type Agent struct {
	llm           llm.LLM
	memory        []memory.MemoryEntry
	tools         []tools.Tool
	maxIterations int
}

func NewAgent(l llm.LLM) *Agent {
	return &Agent{
		llm:           l,
		memory:        []memory.MemoryEntry{},
		tools:         tools.StdToolset,
		maxIterations: 10, // Default max iterations as recommended by OpenRouter docs
	}
}

// Mods
func (a *Agent) WithMemory(memory []memory.MemoryEntry) *Agent {
	a.memory = memory
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

	// Maybe useful to debug later
	// if len(resultStr) > 100 {
	// 	color.Blue("Tool %s returned: %s", call.FunctionName, resultStr[:100])
	// } else {
	// 	color.Blue("Tool %s returned: %s", call.FunctionName, resultStr)
	// }

	return resultStr, nil
}

func (a *Agent) Run(ctx context.Context, prompt string) error {
	initMessage := memory.MemoryEntry{Message: prompt, Role: "user"}
	a.AddMemoryEntry(initMessage)

	iterationCount := 0
	for iterationCount < a.maxIterations {
		iterationCount++

		result, err := a.llm.Stream(ctx, a.memory, a.tools)
		if err != nil {
			log.Fatal("Error calling LLM Stream: ", err)
			return err
		}
		a.AddMemoryEntry(memory.MemoryEntry{Role: "model", Message: result.Text, ToolCalls: result.ToolCalls})

		if len(result.ToolCalls) == 0 {
			break
		}

		for _, call := range result.ToolCalls {
			output, _ := a.HandleToolCall(ctx, call)

			a.AddMemoryEntry(memory.MemoryEntry{
				Role:       "tool",
				Message:    output,
				ToolName:   call.FunctionName,
				ToolCallID: call.ID,
			})
		}
	}

	if iterationCount >= a.maxIterations {
		color.Yellow("\nWarning: Maximum iterations (%d) reached. The agent loop was terminated.\n", a.maxIterations)
	}

	return nil
}
