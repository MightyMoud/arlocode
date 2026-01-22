package ollama_llm

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/tools"
	"github.com/ollama/ollama/api"
)

func makeOllamaTools(agentTools []tools.Tool) []api.Tool {
	var ollamaTools []api.Tool
	for _, tool := range agentTools {
		ollamaTools = append(ollamaTools, api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  mapToToolFunctionParameters(generateJSONSchema(tool.ArgType)),
			},
		})
	}
	return ollamaTools
}

func mapToToolFunctionParameters(m map[string]any) api.ToolFunctionParameters {
	var p api.ToolFunctionParameters
	bts, _ := json.Marshal(m)
	_ = json.Unmarshal(bts, &p)
	return p
}

func mapToToolCallFunctionArguments(m map[string]any) api.ToolCallFunctionArguments {
	var args api.ToolCallFunctionArguments
	bts, _ := json.Marshal(m)
	_ = json.Unmarshal(bts, &args)
	return args
}

func generateJSONSchema(t reflect.Type) map[string]any {
	// Handle pointer types by dereferencing them
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return map[string]any{"type": "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]any{"type": "integer"}
	case reflect.Float32, reflect.Float64:
		return map[string]any{"type": "number"}
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Slice, reflect.Array:
		return map[string]any{
			"type":  "array",
			"items": generateJSONSchema(t.Elem()),
		}
	case reflect.Map:
		return map[string]any{"type": "object"}
	case reflect.Struct:
		properties := make(map[string]any)
		var required []string

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}

			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			name := field.Name
			isOmitEmpty := false
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				name = parts[0]
				for _, part := range parts[1:] {
					if part == "omitempty" {
						isOmitEmpty = true
						break
					}
				}
			}

			// Add to required if not omitempty
			if !isOmitEmpty {
				required = append(required, name)
			}

			propSchema := generateJSONSchema(field.Type)

			// Add description if available - important for LLM understanding
			descTag := field.Tag.Get("description")
			if descTag != "" {
				propSchema["description"] = descTag
			}

			// Add enum constraints if available
			enumTag := field.Tag.Get("enum")
			if enumTag != "" {
				enumValues := strings.Split(enumTag, ",")
				propSchema["enum"] = enumValues
			}

			// Add default value if available
			defaultTag := field.Tag.Get("default")
			if defaultTag != "" {
				propSchema["default"] = defaultTag
			}

			properties[name] = propSchema
		}

		schema := map[string]any{
			"type":       "object",
			"properties": properties,
		}

		// Only add required array if there are required fields
		if len(required) > 0 {
			schema["required"] = required
		}

		return schema
	default:
		return map[string]any{"type": "string"}
	}
}

func convertMemoryToOllamaMessages(mem []memory.MemoryEntry) []api.Message {
	var messages []api.Message
	for _, entry := range mem {
		role := entry.Source
		switch role {
		case "agent", "assistant", "model":
			role = "assistant"
		case "user", "system":
			// keep
		default:
			role = "user"
		}

		msg := api.Message{
			Role:    role,
			Content: entry.Message,
		}

		if role == "assistant" {
			if len(entry.ToolCalls) > 0 {
				var toolCalls []api.ToolCall
				for _, tc := range entry.ToolCalls {
					toolCalls = append(toolCalls, api.ToolCall{
						Function: api.ToolCallFunction{
							Name:      tc.FunctionName,
							Arguments: mapToToolCallFunctionArguments(tc.Arguments),
						},
					})
				}
				msg.ToolCalls = toolCalls
			}
		}
		messages = append(messages, msg)

		if len(entry.Observation.Results) > 0 {
			toolCallNameByID := map[string]string{}
			for _, tc := range entry.ToolCalls {
				toolCallNameByID[tc.ToolCallID] = tc.FunctionName
			}
			for _, result := range entry.Observation.Results {
				name := toolCallNameByID[result.SourceCallID]
				if name == "" && len(entry.ToolCalls) == 1 {
					name = entry.ToolCalls[0].FunctionName
				}
				messages = append(messages, api.Message{
					Role:       "tool",
					Content:    result.Content,
					ToolName:   name,
					ToolCallID: result.SourceCallID,
				})
			}
		}
	}
	return messages
}
