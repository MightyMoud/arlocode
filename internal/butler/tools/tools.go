package tools

import "reflect"

type Tool struct {
	Name        string
	Description string
	Handler     reflect.Value
	ArgType     reflect.Type
}

func NewButlerTool(name, desc string, fn interface{}) *Tool {
	// Reflect to get the function and its signature
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// Simplify -> tools take one arg as struct
	return &Tool{
		Name:        name,
		Description: desc,
		Handler:     fnValue,
		ArgType:     fnType.In(0), // Assumes the tool takes 1 argument (the struct)
	}
}

// ToolCall is a generic representation of an LLM's request to run a tool.
// It is provider-agnostic.
type ToolCall struct {
	ID           string         // Unique ID from the LLM
	FunctionName string         // e.g., "read_file"
	Arguments    map[string]any // The raw arguments (unmarshaled from JSON)
}
