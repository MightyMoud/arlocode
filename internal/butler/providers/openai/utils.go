package openai_provider

import (
	"reflect"
	"strings"

	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

func makeOpenAITools(agentTools []tools.Tool) []openai.ChatCompletionToolUnionParam {
	var openaiTools []openai.ChatCompletionToolUnionParam
	for _, tool := range agentTools {
		openaiTools = append(openaiTools, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
			Name:        tool.Name,
			Description: openai.String(tool.Description),
			Parameters:  generateOpenAISchema(tool.ArgType),
		}))
	}
	return openaiTools
}

func generateOpenAISchema(t reflect.Type) shared.FunctionParameters {
	// Handle pointer types by dereferencing them
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// OpenAI expects a JSON Schema object.
	// We will construct a map[string]interface{} that represents the schema.

	schema := map[string]interface{}{}

	switch t.Kind() {
	case reflect.String:
		schema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Slice, reflect.Array:
		schema["type"] = "array"
		schema["items"] = generateOpenAISchema(t.Elem())
	case reflect.Map:
		schema["type"] = "object"
	case reflect.Struct:
		schema["type"] = "object"
		properties := make(map[string]interface{})
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
			} else {
				// Default to required if no tag? Or optional?
				// Let's assume required for function args usually.
				required = append(required, name)
			}

			properties[name] = generateOpenAISchema(field.Type)
		}
		schema["properties"] = properties
		if len(required) > 0 {
			schema["required"] = required
		}
	default:
		schema["type"] = "string" // Fallback
	}

	return shared.FunctionParameters(schema)
}

func convertMemoryToOpenAIMessages(mem []memory.MemoryEntry) []openai.ChatCompletionMessageParamUnion {
	var messages []openai.ChatCompletionMessageParamUnion
	for _, entry := range mem {
		switch entry.Role {
		case "user":
			messages = append(messages, openai.UserMessage(entry.Message))
		case "assistant", "model":
			messages = append(messages, openai.AssistantMessage(entry.Message))
		case "system":
			messages = append(messages, openai.SystemMessage(entry.Message))
		default:
			// Default to user or handle error?
			messages = append(messages, openai.UserMessage(entry.Message))
		}
	}
	return messages
}
