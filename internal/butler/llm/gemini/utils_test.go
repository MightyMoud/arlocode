package gemini_llm

import (
	"reflect"
	"testing"

	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/tools"
	"google.golang.org/genai"
)

func TestGenerateGenAISchema(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected genai.Type
	}{
		{"String", "test", genai.TypeString},
		{"Int", 123, genai.TypeInteger},
		{"Float", 1.23, genai.TypeNumber},
		{"Bool", true, genai.TypeBoolean},
		{"Slice", []string{"a", "b"}, genai.TypeArray},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := generateGenAISchema(reflect.TypeOf(tt.input))
			if schema.Type != tt.expected {
				t.Errorf("expected type %v, got %v", tt.expected, schema.Type)
			}
		})
	}
}

func TestGenerateGenAISchema_Struct(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2,omitempty"`
	}

	schema := generateGenAISchema(reflect.TypeOf(TestStruct{}))
	if schema.Type != genai.TypeObject {
		t.Errorf("expected type Object, got %v", schema.Type)
	}
	if len(schema.Properties) != 2 {
		t.Errorf("expected 2 properties, got %d", len(schema.Properties))
	}
	if schema.Properties["field1"].Type != genai.TypeString {
		t.Errorf("expected field1 to be String, got %v", schema.Properties["field1"].Type)
	}

	// Check required fields
	found := false
	for _, req := range schema.Required {
		if req == "field1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected field1 to be required")
	}
}

func TestMakeGeminiTools(t *testing.T) {
	handler := func(args struct{ Name string }) (string, error) { return "", nil }
	toolList := []tools.Tool{
		tools.NewButlerTool("test_tool", "description", handler),
	}

	geminiTools := makeGeminiTools(toolList)
	if len(geminiTools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(geminiTools))
	}

	if len(geminiTools[0].FunctionDeclarations) != 1 {
		t.Fatalf("expected 1 function declaration, got %d", len(geminiTools[0].FunctionDeclarations))
	}

	decl := geminiTools[0].FunctionDeclarations[0]
	if decl.Name != "test_tool" {
		t.Errorf("expected name test_tool, got %s", decl.Name)
	}
	if decl.Description != "description" {
		t.Errorf("expected description description, got %s", decl.Description)
	}
}

func TestConvertMemoryToGeminiHistory(t *testing.T) {
	mem := []memory.MemoryEntry{
		{Role: "user", Message: "hello"},
		{Role: "model", Message: "hi", ToolCalls: []tools.ToolCall{
			{ID: "1", FunctionName: "func", Arguments: map[string]any{"a": 1}},
		}},
		{Role: "tool", Message: "result", ToolName: "func", ToolCallID: "1"},
	}

	history := convertMemoryToGeminiHistory(mem)
	if len(history) != 3 {
		t.Fatalf("expected 3 history items, got %d", len(history))
	}

	// Check user message
	if history[0].Role != "user" {
		t.Errorf("expected role user, got %s", history[0].Role)
	}
	if history[0].Parts[0].Text != "hello" {
		t.Errorf("expected text hello, got %s", history[0].Parts[0].Text)
	}

	// Check model message with tool call
	if history[1].Role != "model" {
		t.Errorf("expected role model, got %s", history[1].Role)
	}
	if history[1].Parts[0].Text != "hi" {
		t.Errorf("expected text hi, got %s", history[1].Parts[0].Text)
	}
	if history[1].Parts[1].FunctionCall == nil {
		t.Error("expected function call")
	}

	// Check tool response
	if history[2].Role != "tool" {
		t.Errorf("expected role tool, got %s", history[2].Role)
	}
	if history[2].Parts[0].FunctionResponse == nil {
		t.Error("expected function response")
	}
	if history[2].Parts[0].FunctionResponse.Name != "func" {
		t.Errorf("expected function name func, got %s", history[2].Parts[0].FunctionResponse.Name)
	}
}

func TestGenerateGenAISchema_Map(t *testing.T) {
	schema := generateGenAISchema(reflect.TypeOf(map[string]int{}))
	if schema.Type != genai.TypeObject {
		t.Errorf("expected type Object for map, got %v", schema.Type)
	}
}

func TestGenerateGenAISchema_StructTags(t *testing.T) {
	type TagStruct struct {
		NoTag     string
		DescField string `description:"A description"`
		Ignored   string `json:"-"`
	}

	schema := generateGenAISchema(reflect.TypeOf(TagStruct{}))

	if _, ok := schema.Properties["NoTag"]; !ok {
		t.Error("expected NoTag field")
	}

	if prop, ok := schema.Properties["DescField"]; ok {
		if prop.Description != "A description" {
			t.Errorf("expected description 'A description', got '%s'", prop.Description)
		}
	} else {
		t.Error("expected DescField field")
	}

	if _, ok := schema.Properties["Ignored"]; ok {
		t.Error("expected Ignored field to be ignored")
	}
}

func TestGenerateGenAISchema_Ptr(t *testing.T) {
	i := 10
	schema := generateGenAISchema(reflect.TypeOf(&i))
	if schema.Type != genai.TypeInteger {
		t.Errorf("expected type Integer for pointer to int, got %v", schema.Type)
	}
}

func TestGenerateGenAISchema_MoreTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected genai.Type
	}{
		{"Int8", int8(1), genai.TypeInteger},
		{"Uint", uint(1), genai.TypeInteger},
		{"Float32", float32(1.0), genai.TypeNumber},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := generateGenAISchema(reflect.TypeOf(tt.input))
			if schema.Type != tt.expected {
				t.Errorf("expected type %v, got %v", tt.expected, schema.Type)
			}
		})
	}
}
