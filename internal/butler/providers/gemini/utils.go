package gemini

import (
	"reflect"
	"strings"

	"github.com/mightymoud/sidekick-agent/internal/butler/memory"
	"github.com/mightymoud/sidekick-agent/internal/butler/tools"
	"google.golang.org/genai"
)

func makeGeminiTools(tools []tools.Tool) []*genai.Tool {
	var geminiFunDecls []*genai.FunctionDeclaration
	for _, tool := range tools {
		geminiFunDecls = append(geminiFunDecls, &genai.FunctionDeclaration{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  generateGenAISchema(tool.ArgType),
		})
	}
	return []*genai.Tool{{FunctionDeclarations: geminiFunDecls}}
}

func generateGenAISchema(t reflect.Type) *genai.Schema {
	// Handle pointer types by dereferencing them
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return &genai.Schema{
			Type: genai.TypeString,
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &genai.Schema{
			Type: genai.TypeInteger,
		}
	case reflect.Float32, reflect.Float64:
		return &genai.Schema{
			Type: genai.TypeNumber,
		}
	case reflect.Bool:
		return &genai.Schema{
			Type: genai.TypeBoolean,
		}
	case reflect.Slice, reflect.Array:
		return &genai.Schema{
			Type:  genai.TypeArray,
			Items: generateGenAISchema(t.Elem()),
		}
	case reflect.Map:
		// Maps are treated as objects with dynamic keys if possible,
		// but GenAI schema usually expects specific properties for objects.
		// For a generic map[string]T, we might not be able to fully represent it
		// without knowing the keys ahead of time, or we treat it as an object
		// where additional properties are allowed.
		// However, standard JSON schema for maps usually implies an object.
		// Let's assume string keys.
		return &genai.Schema{
			Type: genai.TypeObject,
			// In OpenAPI/JSON Schema, maps are objects with "additionalProperties".
			// The genai.Schema struct might not strictly support additionalProperties in the same way
			// as full JSON schema, but we can try to approximate or just return TypeObject.
		}
	case reflect.Struct:
		properties := make(map[string]*genai.Schema)
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
				// Check if "omitempty" is NOT present to mark as required?
				// Usually in JSON schema generation from Go structs:
				// - If no tag, it's required (or optional depending on convention).
				// - If `omitempty` is present, it's definitely optional.
				// Let's assume fields are optional unless we decide otherwise,
				// but for function calling, usually we want to be explicit.
				// A common convention: if it doesn't have omitempty, it's required.
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
				// No JSON tag, default to field name, assume optional or required?
				// Let's assume optional to be safe, or required.
				// For strict schemas, maybe required.
				// Let's stick to: if no tag, use field name, treat as optional.
			}

			properties[name] = generateGenAISchema(field.Type)

			// Add description if available (e.g. from a "description" tag)
			descTag := field.Tag.Get("description")
			if descTag != "" {
				properties[name].Description = descTag
			}

			// Handle enums if there's a tag or specific type logic (omitted for brevity)
		}

		return &genai.Schema{
			Type:       genai.TypeObject,
			Properties: properties,
			Required:   required,
		}
	default:
		// Fallback for unknown types (interface{}, func, chan, etc.)
		return &genai.Schema{
			Type: genai.TypeString, // Safe fallback or maybe TypeObject
		}
	}
}

func convertMemoryToGeminiHistory(memory []memory.MemoryEntry) []*genai.Content {
	history := []*genai.Content{}
	for _, entry := range memory {
		var genAIEntry genai.Content
		if entry.Role == "tool" {
			genAIEntry = genai.Content{
				Role: "tool",
				Parts: []*genai.Part{{
					FunctionResponse: &genai.FunctionResponse{
						Name:     entry.ToolName,
						Response: map[string]any{"content": entry.Message},
						ID:       entry.ToolCallID,
					},
				}},
			}
		} else {
			parts := []*genai.Part{}
			if entry.Message != "" {
				parts = append(parts, &genai.Part{Text: entry.Message})
			}
			for _, tc := range entry.ToolCalls {
				parts = append(parts, &genai.Part{
					ThoughtSignature: tc.ThoughtSignature,
					FunctionCall: &genai.FunctionCall{
						Name: tc.FunctionName,
						Args: tc.Arguments,
						ID:   tc.ID,
					},
				})
			}
			genAIEntry = genai.Content{
				Role:  entry.Role,
				Parts: parts,
			}
		}
		history = append(history, &genAIEntry)
	}
	return history
}
