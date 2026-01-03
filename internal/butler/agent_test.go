package butler

import (
	"context"
	"testing"

	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/providers"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
)

// MockLLM is a mock implementation of the LLM interface
type MockLLM struct {
	StreamFunc   func(ctx context.Context, mem []memory.MemoryEntry, t []tools.Tool) (providers.ProviderResponse, error)
	GenerateFunc func(ctx context.Context, mem []memory.MemoryEntry, t []tools.Tool) error
}

func (m *MockLLM) Stream(ctx context.Context, mem []memory.MemoryEntry, t []tools.Tool) (providers.ProviderResponse, error) {
	if m.StreamFunc != nil {
		return m.StreamFunc(ctx, mem, t)
	}
	return providers.ProviderResponse{}, nil
}

func (m *MockLLM) Generate(ctx context.Context, mem []memory.MemoryEntry, t []tools.Tool) error {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, mem, t)
	}
	return nil
}

func TestNewAgent(t *testing.T) {
	mockLLM := &MockLLM{}
	agent := NewAgent(mockLLM)

	if agent == nil {
		t.Fatal("NewAgent returned nil")
	}
	if agent.llm != mockLLM {
		t.Errorf("Expected llm to be %v, got %v", mockLLM, agent.llm)
	}
	if len(agent.memory) != 0 {
		t.Errorf("Expected empty memory, got %v", agent.memory)
	}
	if len(agent.tools) == 0 {
		t.Error("Expected default tools, got empty")
	}
}

func TestAgent_WithMemory(t *testing.T) {
	mockLLM := &MockLLM{}
	agent := NewAgent(mockLLM)
	mem := []memory.MemoryEntry{{Message: "test", Role: "user"}}

	agent.WithMemory(mem)
	if len(agent.memory) != 1 || agent.memory[0].Message != "test" {
		t.Errorf("WithMemory failed")
	}
}

func TestAgent_WitTools(t *testing.T) {
	mockLLM := &MockLLM{}
	agent := NewAgent(mockLLM)
	newTools := []tools.Tool{}

	agent.WitTools(newTools)
	if len(agent.tools) != 0 {
		t.Errorf("WitTools failed")
	}
}

func TestAgent_WithNoTools(t *testing.T) {
	mockLLM := &MockLLM{}
	agent := NewAgent(mockLLM)

	agent.WithNoTools()
	if len(agent.tools) != 0 {
		t.Errorf("WithNoTools failed")
	}
}

func TestAgent_AddMemoryEntry(t *testing.T) {
	mockLLM := &MockLLM{}
	agent := NewAgent(mockLLM)
	entry := memory.MemoryEntry{Message: "test", Role: "user"}

	agent.AddMemoryEntry(entry)
	if len(agent.memory) != 1 {
		t.Errorf("AddMemoryEntry failed")
	}
}

func TestAgent_GetMemory(t *testing.T) {
	mockLLM := &MockLLM{}
	agent := NewAgent(mockLLM)
	entry := memory.MemoryEntry{Message: "test", Role: "user"}
	agent.AddMemoryEntry(entry)

	mem := agent.GetMemory()
	if len(mem) != 1 {
		t.Errorf("GetMemory failed")
	}
}

type MockToolArgs struct {
	Input string `json:"input"`
}

func MockToolHandler(args MockToolArgs) (string, error) {
	return "processed: " + args.Input, nil
}

func TestAgent_HandleToolCall(t *testing.T) {
	mockLLM := &MockLLM{}
	agent := NewAgent(mockLLM)

	mockTool := tools.NewButlerTool("mock_tool", "mock description", MockToolHandler)
	agent.WitTools([]tools.Tool{mockTool})

	call := tools.ToolCall{
		ID:           "call_1",
		FunctionName: "mock_tool",
		Arguments:    map[string]any{"input": "test"},
	}

	result, err := agent.HandleToolCall(context.Background(), call)
	if err != nil {
		t.Fatalf("HandleToolCall failed: %v", err)
	}
	if result != "processed: test" {
		t.Errorf("Expected 'processed: test', got '%s'", result)
	}
}

func TestAgent_Run(t *testing.T) {
	mockLLM := &MockLLM{
		StreamFunc: func(ctx context.Context, mem []memory.MemoryEntry, t []tools.Tool) (providers.ProviderResponse, error) {
			return providers.ProviderResponse{
				Text: "world",
			}, nil
		},
	}
	agent := NewAgent(mockLLM)
	agent.WithNoTools()

	ctx := context.Background()
	prompt := "hello"

	err := agent.Run(ctx, prompt)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	mem := agent.GetMemory()
	if len(mem) != 2 {
		t.Fatalf("Expected 2 memory entries, got %d", len(mem))
	}
	if mem[0].Message != "hello" {
		t.Errorf("Expected first message 'hello', got '%s'", mem[0].Message)
	}
	if mem[1].Message != "world" {
		t.Errorf("Expected second message 'world', got '%s'", mem[1].Message)
	}
}

func TestAgent_Run_WithToolCall(t *testing.T) {
	mockTool := tools.NewButlerTool("mock_tool", "mock description", MockToolHandler)

	callCount := 0
	mockLLM := &MockLLM{
		StreamFunc: func(ctx context.Context, mem []memory.MemoryEntry, t []tools.Tool) (providers.ProviderResponse, error) {
			callCount++
			if callCount == 1 {
				return providers.ProviderResponse{
					Text: "I will use the tool",
					ToolCalls: []tools.ToolCall{
						{
							ID:           "call_1",
							FunctionName: "mock_tool",
							Arguments:    map[string]any{"input": "test"},
						},
					},
				}, nil
			}
			return providers.ProviderResponse{
				Text: "Final answer",
			}, nil
		},
	}

	agent := NewAgent(mockLLM)
	agent.WitTools([]tools.Tool{mockTool})

	ctx := context.Background()
	prompt := "do something"

	err := agent.Run(ctx, prompt)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	mem := agent.GetMemory()
	// 1. User prompt
	// 2. Model response (tool call)
	// 3. Tool output
	// 4. Model response (final)
	if len(mem) != 4 {
		t.Fatalf("Expected 4 memory entries, got %d", len(mem))
	}
	if mem[2].Role != "tool" {
		t.Errorf("Expected 3rd message role 'tool', got '%s'", mem[2].Role)
	}
	if mem[2].Message != "processed: test" {
		t.Errorf("Expected tool output 'processed: test', got '%s'", mem[2].Message)
	}
}
