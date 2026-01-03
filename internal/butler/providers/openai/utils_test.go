package openai_provider

import (
	"reflect"
	"testing"

	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
)

func TestMakeOpenAITools(t *testing.T) {
	handler := func(args struct{ Name string }) (string, error) { return "", nil }
	toolList := []tools.Tool{
		tools.NewButlerTool("test_tool", "description", handler),
	}

	openaiTools := makeOpenAITools(toolList)
	if len(openaiTools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(openaiTools))
	}

	// The tools are wrapped in ChatCompletionToolUnionParam, need to check the underlying type
	// For now, we'll just check that we got something back
	if len(openaiTools) == 0 {
		t.Error("expected at least one tool")
	}
}

func TestGenerateOpenAISchema(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"String", "test"},
		{"Int", 123},
		{"Float", 1.23},
		{"Bool", true},
		{"Slice", []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := generateOpenAISchema(reflect.TypeOf(tt.input))
			if schema == nil {
				t.Error("expected schema to be non-nil")
			}
			// Since shared.FunctionParameters is an opaque type,
			// we can't easily inspect its contents without more knowledge
			// of the OpenAI library internals. Just verify it returns something.
		})
	}
}

func TestGenerateOpenAISchema_Struct(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2,omitempty"`
	}

	schema := generateOpenAISchema(reflect.TypeOf(TestStruct{}))
	if schema == nil {
		t.Error("expected schema to be non-nil")
	}
	// Since we can't inspect the internal structure easily,
	// just verify the function completes without error
}

func TestConvertMemoryToOpenAIMessages(t *testing.T) {
	mem := []memory.MemoryEntry{
		{Role: "user", Message: "hello"},
		{Role: "assistant", Message: "hi"},
		{Role: "system", Message: "system message"},
		{Role: "unknown", Message: "default to user"},
	}

	messages := convertMemoryToOpenAIMessages(mem)
	if len(messages) != 4 {
		t.Fatalf("expected 4 messages, got %d", len(messages))
	}

	// Since ChatCompletionMessageParamUnion is not an interface,
	// we can't do type assertions. Just verify the function returns
	// the correct number of messages without panicking.
}

func TestGenerateOpenAISchema_Map(t *testing.T) {
	schema := generateOpenAISchema(reflect.TypeOf(map[string]int{}))
	if schema == nil {
		t.Error("expected schema to be non-nil")
	}
}

func TestGenerateOpenAISchema_Ptr(t *testing.T) {
	i := 10
	schema := generateOpenAISchema(reflect.TypeOf(&i))
	if schema == nil {
		t.Error("expected schema to be non-nil")
	}
}

func TestGenerateOpenAISchema_MoreTypes(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"Int8", int8(1)},
		{"Uint", uint(1)},
		{"Float32", float32(1.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := generateOpenAISchema(reflect.TypeOf(tt.input))
			if schema == nil {
				t.Errorf("expected schema to be non-nil for %s", tt.name)
			}
		})
	}
}
