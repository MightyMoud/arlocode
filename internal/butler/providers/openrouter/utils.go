package openrouter

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/iamwavecut/gopenrouter"
	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
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
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				name = parts[0]
				isOmitEmpty := false
				for _, part := range parts[1:] {
					if part == "omitempty" {
						isOmitEmpty = true
						break
					}
				}
				if !isOmitEmpty {
					required = append(required, name)
				}
			}

			propSchema := generateJSONSchema(field.Type)

			// Add description if available
			descTag := field.Tag.Get("description")
			if descTag != "" {
				propSchema["description"] = descTag
			}

			properties[name] = propSchema
		}

		return map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   required,
		}
	default:
		return map[string]any{"type": "string"}
	}
}

func getRoleFromMemoryEntry(entry memory.MemoryEntry) gopenrouter.ChatCompletionMessageRole {
	switch entry.Role {
	case "user":
		return gopenrouter.RoleUser
	case "assistant":
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
				Role:    getRoleFromMemoryEntry(entry),
				Content: entry.Message,
			}
			if len(toolCalls) > 0 {
				msg.ToolCalls = toolCalls
			}
			messages = append(messages, msg)
		}
	}
	return messages
}
