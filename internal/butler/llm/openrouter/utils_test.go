package openrouter_llm

import (
	"reflect"
	"testing"

	"github.com/iamwavecut/gopenrouter"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

func TestMakeOpenRouterTools(t *testing.T) {
	handler := func(args struct{ Name string }) (string, error) { return "", nil }
	toolList := []tools.Tool{
		tools.NewButlerTool("test_tool", "description", handler),
	}

	openRouterTools := makeOpenRouterTools(toolList)
	if len(openRouterTools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(openRouterTools))
	}

	tool := openRouterTools[0]
	if tool.Type != "function" {
		t.Errorf("expected type function, got %s", tool.Type)
	}
	if tool.Function.Name != "test_tool" {
		t.Errorf("expected name test_tool, got %s", tool.Function.Name)
	}
	if tool.Function.Description != "description" {
		t.Errorf("expected description description, got %s", tool.Function.Description)
	}
}

func TestGenerateJSONSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected map[string]any
	}{
		{"String", "test", map[string]any{"type": "string"}},
		{"Int", 123, map[string]any{"type": "integer"}},
		{"Float", 1.23, map[string]any{"type": "number"}},
		{"Bool", true, map[string]any{"type": "boolean"}},
		{"Slice", []string{"a", "b"}, map[string]any{"type": "array", "items": map[string]any{"type": "string"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := generateJSONSchema(reflect.TypeOf(tt.input))
			if schema["type"] != tt.expected["type"] {
				t.Errorf("expected type %v, got %v", tt.expected["type"], schema["type"])
			}
			if tt.name == "Slice" {
				if items, ok := schema["items"].(map[string]any); ok {
					if items["type"] != "string" {
						t.Errorf("expected items type string, got %v", items["type"])
					}
				} else {
					t.Error("expected items to be a map")
				}
			}
		})
	}
}

func TestGenerateJSONSchema_Struct(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2,omitempty"`
		Field3 string `json:"field3" description:"A description"`
	}

	schema := generateJSONSchema(reflect.TypeOf(TestStruct{}))
	if schema["type"] != "object" {
		t.Errorf("expected type object, got %v", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties to be a map")
	}
	if len(properties) != 3 {
		t.Errorf("expected 3 properties, got %d", len(properties))
	}

	if field1, ok := properties["field1"].(map[string]any); ok {
		if field1["type"] != "string" {
			t.Errorf("expected field1 type string, got %v", field1["type"])
		}
	} else {
		t.Error("expected field1 to be a map")
	}

	if field3, ok := properties["field3"].(map[string]any); ok {
		if field3["description"] != "A description" {
			t.Errorf("expected description 'A description', got %v", field3["description"])
		}
	} else {
		t.Error("expected field3 to be a map")
	}

	required, ok := schema["required"].([]string)
	if !ok {
		t.Fatal("expected required to be a slice")
	}
	if len(required) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(required))
	}
	found := false
	for _, r := range required {
		if r == "field1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected field1 to be required")
	}
}

func TestGetRoleFromMemoryEntry(t *testing.T) {
	tests := []struct {
		input    string
		expected gopenrouter.ChatCompletionMessageRole
	}{
		{"user", gopenrouter.RoleUser},
		{"assistant", gopenrouter.RoleAssistant},
		{"system", gopenrouter.RoleSystem},
		{"unknown", gopenrouter.RoleUser}, // default case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			entry := memory.MemoryEntry{Role: tt.input}
			role := getRoleFromMemoryEntry(entry)
			if role != tt.expected {
				t.Errorf("expected role %v, got %v", tt.expected, role)
			}
		})
	}
}

func TestConvertMemoryToOpenRouterMessages(t *testing.T) {
	mem := []memory.MemoryEntry{
		{Role: "user", Message: "hello"},
		{Role: "assistant", Message: "hi", ToolCalls: []tools.ToolCall{
			{ID: "1", FunctionName: "func", Arguments: map[string]any{"a": 1}},
		}},
		{Role: "tool", Message: "result", ToolCallID: "1"},
	}

	messages := convertMemoryToOpenRouterMessages(mem)
	if len(messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(messages))
	}

	// Check user message
	if messages[0].Role != gopenrouter.RoleUser {
		t.Errorf("expected role user, got %v", messages[0].Role)
	}
	if messages[0].Content != "hello" {
		t.Errorf("expected content hello, got %s", messages[0].Content)
	}

	// Check assistant message with tool call
	if messages[1].Role != gopenrouter.RoleAssistant {
		t.Errorf("expected role assistant, got %v", messages[1].Role)
	}
	if messages[1].Content != "hi" {
		t.Errorf("expected content hi, got %s", messages[1].Content)
	}
	if len(messages[1].ToolCalls) != 1 {
		t.Errorf("expected 1 tool call, got %d", len(messages[1].ToolCalls))
	}
	if messages[1].ToolCalls[0].Function.Name != "func" {
		t.Errorf("expected function name func, got %s", messages[1].ToolCalls[0].Function.Name)
	}

	// Check tool message
	if messages[2].Role != gopenrouter.RoleTool {
		t.Errorf("expected role tool, got %v", messages[2].Role)
	}
	if messages[2].Content != "result" {
		t.Errorf("expected content result, got %s", messages[2].Content)
	}
	if messages[2].ToolCallID != "1" {
		t.Errorf("expected tool call ID 1, got %s", messages[2].ToolCallID)
	}
}
