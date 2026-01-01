package toolset

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ToolMetadata is what we send to the LLM (OpenAI/Gemini format)
type ToolMetadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

type Registry struct {
	tools map[string]reflect.Value
}

func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]reflect.Value)}
}

// Register takes any Go function and adds it to our framework
func (r *Registry) Register(name, desc string, fn interface{}) ToolMetadata {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic("Register requires a function")
	}
	r.tools[name] = v

	// Basic Reflection to build a JSON Schema
	properties := make(map[string]map[string]string)
	t := v.Type()

	// We assume function arguments are named arg0, arg1...
	// In a real framework, you'd use a struct as a single argument to get field names.
	for i := 0; i < t.NumIn(); i++ {
		argName := fmt.Sprintf("arg%d", i)
		properties[argName] = map[string]string{
			"type": goTypeToJSONType(t.In(i).Kind()),
		}
	}

	return ToolMetadata{
		Name:        name,
		Description: desc,
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": properties,
		},
	}
}

func goTypeToJSONType(k reflect.Kind) string {
	switch k {
	case reflect.Int, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	default:
		return "string"
	}
}

// Call executes the stored function using JSON arguments from the LLM
func (r *Registry) Call(name string, jsonArgs string) (string, error) {
	fn, ok := r.tools[name]
	if !ok {
		return "", fmt.Errorf("tool not found")
	}

	// 1. Unmarshal JSON into a map
	var argsMap map[string]interface{}
	json.Unmarshal([]byte(jsonArgs), &argsMap)

	// 2. Prepare reflect.Values for calling
	in := make([]reflect.Value, fn.Type().NumIn())
	for i := 0; i < fn.Type().NumIn(); i++ {
		val := argsMap[fmt.Sprintf("arg%d", i)]
		in[i] = reflect.ValueOf(val).Convert(fn.Type().In(i))
	}

	// 3. Execute the Go function
	results := fn.Call(in)
	return fmt.Sprintf("%v", results[0].Interface()), nil
}
