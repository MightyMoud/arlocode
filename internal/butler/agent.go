package butler

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
)

type Agent struct {
	llm    LLM
	memory []memory.MemoryEntry
	tools  []tools.Tool
}

func NewAgent(l LLM) *Agent {
	return &Agent{llm: l, memory: []memory.MemoryEntry{}, tools: []tools.Tool{}}
}

// Mock for Memory stuff later this is where Agent will use it
func (a *Agent) AddMemoryEntry(entry memory.MemoryEntry) {
	a.memory = append(a.memory, entry)
}

// Memory stuff later
func (a *Agent) GetMemory() []memory.MemoryEntry {
	return a.memory
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

func (a *Agent) HandleToolCall(ctx context.Context, call tools.ToolCall) (string, error) {
	var tool tools.Tool
	for _, t := range a.tools {
		if t.Name == call.FunctionName {
			tool = t
			break
		}
	}

	// 2. Create a fresh instance of the specific Arg struct (e.g., ReadFileArgs)
	// reflect.New returns a pointer to the type
	argsPtr := reflect.New(tool.ArgType).Interface()

	// 3. Convert the map[string]any into the concrete struct.
	// This works for ANY provider because we go through JSON bytes first.
	bytes, _ := json.Marshal(call.Arguments)
	if err := json.Unmarshal(bytes, argsPtr); err != nil {
		return "", fmt.Errorf("failed to unmarshal tool args: %w", err)
	}

	// 4. Execute the function.
	// .Elem() dereferences the pointer so we pass the struct value.
	results := tool.Handler.Call([]reflect.Value{
		reflect.ValueOf(argsPtr).Elem(),
	})

	// 5. Handle potential errors from the tool itself
	if len(results) > 1 && !results[1].IsNil() {
		return "", results[1].Interface().(error)
	}

	return results[0].String(), nil
}

func (a *Agent) Run(ctx context.Context, prompt string) error {
	initMessage := memory.MemoryEntry{Message: prompt, Role: "user"}
	a.AddMemoryEntry(initMessage)

	for {
		result, err := a.llm.Stream(ctx, a.memory, a.tools)
		if err != nil {
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
	return nil
}
