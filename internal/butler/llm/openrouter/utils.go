package openrouter_llm

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/iamwavecut/gopenrouter"
	"github.com/mightymoud/arlocode/internal/butler/memory"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

func makeOpenRouterTools(agentTools []tools.Tool) []gopenrouter.Tool {
	var openRouterTools []gopenrouter.Tool
	for _, tool := range agentTools {
		openRouterTools = append(openRouterTools, gopenrouter.Tool{
			Type: "function",
			Function: gopenrouter.Function{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  generateJSONSchema(tool.ArgType),
			},
		})
	}
	return openRouterTools
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

func getRoleFromMemoryEntry(entry memory.MemoryEntry) gopenrouter.ChatCompletionMessageRole {
	switch entry.Role {
	case "user":
		return gopenrouter.RoleUser
	case "assistant", "model":
		return gopenrouter.RoleAssistant
	case "system":
		return gopenrouter.RoleSystem
	default:
		return gopenrouter.RoleUser
	}
}

func convertMemoryToOpenRouterMessages(mem []memory.MemoryEntry) []gopenrouter.ChatCompletionMessage {
	var messages []gopenrouter.ChatCompletionMessage
	for _, entry := range mem {
		if entry.Role == "tool" {
			messages = append(messages, gopenrouter.ChatCompletionMessage{
				Role:       gopenrouter.RoleTool,
				Content:    entry.Message,
				ToolCallID: entry.ToolCallID,
				Name:       entry.ToolName,
			})
		} else {
			var toolCalls []gopenrouter.ToolCall
			for _, tc := range entry.ToolCalls {
				argsBytes, _ := json.Marshal(tc.Arguments)
				toolCalls = append(toolCalls, gopenrouter.ToolCall{
					ID:   tc.ID,
					Type: "function",
					Function: gopenrouter.Function{
						Name:      tc.FunctionName,
						Arguments: string(argsBytes),
					},
				})
			}

			msg := gopenrouter.ChatCompletionMessage{
				Role: getRoleFromMemoryEntry(entry),
			}

			// Per OpenRouter docs: when assistant makes tool calls, content should be null/empty
			// Only set content if there are no tool calls or if content exists
			if len(toolCalls) > 0 {
				msg.ToolCalls = toolCalls
				// Content is intentionally not set (null) when tool calls are present
				// unless there's actual content (for interleaved thinking)
				if entry.Message != "" {
					msg.Content = entry.Message
				}
			} else {
				msg.Content = entry.Message
			}

			messages = append(messages, msg)
		}
	}
	return messages
}
