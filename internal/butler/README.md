# Butler Package

**Butler** is a flexible, provider-agnostic AI agent framework for Go that enables LLM-powered automation with tool calling capabilities.

> **Note:** This is an internal package used by the arlocode project. It provides the core AI agent functionality for automated code analysis and manipulation tasks.

## Overview

The butler package provides an agent-based architecture that can work with multiple LLM providers (OpenAI, Gemini, OpenRouter) and execute tools/functions based on the agent's decisions. It maintains conversation memory, supports streaming responses with real-time feedback, and uses iterative reasoning loops until a task is complete.

## Key Features

- 🔌 **Provider-Agnostic**: Swap between OpenAI, Gemini, OpenRouter without changing your code
- 🛠️ **Tool Calling**: Execute custom functions and file operations based on AI decisions
- 🧠 **Memory Management**: Track and maintain conversation history across iterations
- 📡 **Streaming Responses**: Real-time output of text, thinking process, and tool calls
- 🔄 **Iterative Execution**: Automatically loops until task completion or max iterations
- 🎯 **Event Hooks**: Customize behavior with callbacks for text, thinking, and tool calls

## Core Components

### Agent

The `Agent` is the central orchestrator that:
- Manages conversation memory
- Coordinates with LLM providers
- Handles tool execution
- Controls iteration loops to prevent infinite execution
- Supports event hooks for real-time streaming

### LLM Interface

A provider-agnostic interface that allows different LLM providers to be used interchangeably:
- `Stream(ctx, memory, tools, hooks)`: Streams responses from the LLM with tool calling support
- `Generate()`: Generates responses from the LLM

### Memory

The memory system tracks conversation history including:
- User messages
- Model responses
- Tool calls and their results
- Role-based message organization (`user`, `model`, `tool`)

### Tools

Built-in tools for common operations:
- `read_file`: Read file contents from disk
- `read_folder`: Read all files in a folder recursively
- `list_folder_contents`: List files and directories in a folder
- `search_code`: Search for text across files in a folder
- `apply_edit`: Apply precise text replacements to files
- `fetch_url_as_markdown`: Fetch web content and convert to markdown
- `make_file`: Create new files with content

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/mightymoud/arlocode/internal/butler/agent"
    "github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
    "github.com/mightymoud/arlocode/internal/butler/tools"
)

func main() {
    ctx := context.Background()
    
    // Create a provider and model
    provider := openrouter.New(ctx)
    model := provider.Model(ctx, "anthropic/claude-sonnet-4.5")
    
    // Create and configure an agent
    agent := agent.NewAgent(model).WithMaxIterations(50)
    
    // Run the agent with a prompt
    prompt := "Read the main.go file and explain what it does"
    err := agent.Run(ctx, prompt)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Advanced Example with Event Hooks

Event hooks allow you to stream responses in real-time and monitor tool execution:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/fatih/color"
    "github.com/mightymoud/arlocode/internal/butler/agent"
    "github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
    "github.com/mightymoud/arlocode/internal/butler/tools"
)

func main() {
    ctx := context.Background()
    
    provider := openrouter.New(ctx)
    model := provider.Model(ctx, "z-ai/glm-4.7")
    
    // Create an agent with event hooks
    agent := agent.NewAgent(model).
        WithMaxIterations(50).
        WithOnThinkingChunk(func(chunk string) {
            // Stream thinking/reasoning in orange
            color.RGB(255, 128, 0).Printf("%s", chunk)
        }).
        WithOnTextChunk(func(chunk string) {
            // Stream regular text output
            fmt.Printf("%s", chunk)
        }).
        WithOnToolCall(func(t tools.ToolCall) {
            // Log tool calls in blue
            color.Blue("\n[Tool Call] %s - with Arguments: %+v", 
                t.FunctionName, t.Arguments)
        })
    
    prompt := "Analyze the code structure of the butler package"
    agent.Run(ctx, prompt)
}
```

### Customizing the Agent

#### Using Custom Tools

```go
import "github.com/mightymoud/arlocode/internal/butler/tools"

// Define your tool's arguments struct
type myToolArgs struct {
    Input string `json:"input" jsonschema:"Description of the input parameter"`
    Count int    `json:"count" jsonschema:"Number of times to process"`
}

// Define your tool function
func myToolFunc(args myToolArgs) (string, error) {
    // Tool implementation
    result := fmt.Sprintf("Processed '%s' %d times", args.Input, args.Count)
    return result, nil
}

// Create the tool
customTool := tools.NewButlerTool(
    "my_tool_name",
    "Description of what the tool does",
    myToolFunc,
)

// Create an agent with custom tools
customTools := []tools.Tool{customTool}
agent := agent.NewAgent(model).WitTools(customTools)
```

#### Using No Tools (Chat-Only Mode)

```go
// Create an agent without any tools - just for conversation
agent := agent.NewAgent(model).WithNoTools()
```

#### Setting Max Iterations

```go
// Limit the agent to 20 iterations (default is 10)
agent := agent.NewAgent(model).WithMaxIterations(20)

// For complex tasks, increase iterations
agent := agent.NewAgent(model).WithMaxIterations(100)
```

#### Using Existing Memory

Continue from a previous conversation:

```go
import "github.com/mightymoud/arlocode/internal/butler/memory"

existingMemory := []memory.MemoryEntry{
    {Role: "user", Message: "Hello"},
    {Role: "model", Message: "Hi! How can I help?"},
    {Role: "user", Message: "What can you do?"},
    {Role: "model", Message: "I can read files, analyze code, and more!"},
}

agent := agent.NewAgent(model).WithMemory(existingMemory)
```

## Supported Providers

### OpenRouter

OpenRouter provides access to many different LLM models:

```go
import "github.com/mightymoud/arlocode/internal/butler/providers/openrouter"

provider := openrouter.New(ctx)
model := provider.Model(ctx, "anthropic/claude-sonnet-4.5")
// or
model = provider.Model(ctx, "openai/gpt-4-turbo")
// or
model = provider.Model(ctx, "z-ai/glm-4.7")
```

### Gemini

Google's Gemini models:

```go
import "github.com/mightymoud/arlocode/internal/butler/providers/gemini"

provider := gemini.New(ctx)
model := provider.Model(ctx, "gemini-3-flash-preview")
```

### OpenAI

Official OpenAI models:

```go
import "github.com/mightymoud/arlocode/internal/butler/providers/openai"

provider := openai.New(ctx)
model := provider.Model(ctx, "gpt-4")
// or
model = provider.Model(ctx, "gpt-3.5-turbo")
```

## Creating Custom Tools

Tools are created using the `NewButlerTool` function which uses reflection to handle any function with a struct argument:

### Tool Function Requirements

1. **Arguments**: Must be a single struct with JSON tags
2. **Return Values**: Must return `(string, error)`
3. **JSON Schema**: Use `jsonschema` tags to describe parameters for the LLM

### Example Tool

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

// 1. Define the arguments struct
type createDirectoryArgs struct {
    Path        string `json:"path" jsonschema:"The directory path to create"`
    MakeParents bool   `json:"make_parents" jsonschema:"Create parent directories if they don't exist"`
}

// 2. Define the tool function
func createDirectory(args createDirectoryArgs) (string, error) {
    if args.MakeParents {
        err := os.MkdirAll(args.Path, 0755)
        if err != nil {
            return "", fmt.Errorf("failed to create directory: %w", err)
        }
        return fmt.Sprintf("Created directory: %s (with parents)", args.Path), nil
    }
    
    err := os.Mkdir(args.Path, 0755)
    if err != nil {
        return "", fmt.Errorf("failed to create directory: %w", err)
    }
    return fmt.Sprintf("Created directory: %s", args.Path), nil
}

// 3. Create the tool
func main() {
    dirTool := tools.NewButlerTool(
        "create_directory",
        "Creates a new directory at the specified path",
        createDirectory,
    )
    
    // Use the tool with an agent
    agent := agent.NewAgent(model).WitTools([]tools.Tool{dirTool})
}
```

## How It Works

1. **Initialization**: Create an agent with an LLM provider and optional configuration
2. **Prompt**: Submit a user prompt via `agent.Run(ctx, prompt)`
3. **Iteration Loop**: The agent enters a loop where it:
   - Sends the current memory (conversation history) to the LLM
   - Receives a streaming response (text, thinking, and/or tool calls)
   - Executes any requested tool calls using reflection
   - Adds tool results to memory
   - Repeats until no more tool calls are needed or max iterations is reached
4. **Completion**: The loop exits when the LLM decides the task is complete

### Event Hooks Flow

```
User Prompt
    ↓
Agent.Run()
    ↓
LLM.Stream() → Event Hooks Triggered:
    ├─ OnThinkingChunk() → Stream reasoning
    ├─ OnTextChunk()     → Stream response text
    └─ OnToolCall()      → Notify of tool execution
    ↓
Tool Execution (via reflection)
    ↓
Add results to memory
    ↓
Repeat if more tool calls needed
```

## Memory Management

### Retrieve Conversation History

```go
history := agent.GetMemory()
for _, entry := range history {
    fmt.Printf("%s: %s\n", entry.Role, entry.Message)
}
```

### Add Custom Memory Entries

```go
import "github.com/mightymoud/arlocode/internal/butler/memory"

agent.AddMemoryEntry(memory.MemoryEntry{
    Role:    "user",
    Message: "Custom message",
})

agent.AddMemoryEntry(memory.MemoryEntry{
    Role:       "tool",
    Message:    "Tool output",
    ToolName:   "my_tool",
    ToolCallID: "call_123",
})
```

### Memory Entry Structure

```go
type MemoryEntry struct {
    Role       string            // "user", "model", "tool"
    Message    string            // The content
    ToolCalls  []ToolCall        // For model entries with tool calls
    ToolName   string            // For tool entries
    ToolCallID string            // For tool entries
}
```

## Event Hooks Reference

### OnTextChunk

Called for each chunk of regular text output from the LLM:

```go
.WithOnTextChunk(func(chunk string) {
    fmt.Print(chunk) // Stream directly to console
})
```

### OnThinkingChunk

Called for each chunk of thinking/reasoning output (provider-dependent):

```go
.WithOnThinkingChunk(func(chunk string) {
    fmt.Printf("[Thinking] %s", chunk)
})
```

### OnToolCall

Called when the LLM requests to execute a tool:

```go
.WithOnToolCall(func(t tools.ToolCall) {
    fmt.Printf("Executing: %s with args: %+v\n", t.FunctionName, t.Arguments)
})
```

## Best Practices

1. **Set Appropriate Max Iterations**: 
   - Default is 10 (recommended by OpenRouter)
   - Increase to 50-100 for complex multi-step tasks
   - Decrease to 5-10 for simple queries

2. **Choose the Right Provider**: 
   - Different providers have different strengths, speeds, and pricing
   - Consider your specific use case when selecting

3. **Use Specific Tools**: 
   - Only include tools relevant to your task
   - Too many tools can confuse the LLM and degrade performance

4. **Monitor Tool Calls**: 
   - Use event hooks to track tool execution
   - Helps debug and understand the agent's reasoning

5. **Handle Errors**: 
   - Always check the error returned by `Run()`
   - The agent will continue on tool errors but may fail completely on LLM errors

6. **Design Good Tool Interfaces**: 
   - Use clear, descriptive JSON schema tags
   - Keep tool functions focused and single-purpose
   - Return meaningful error messages

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Agent                               │
├─────────────────────────────────────────────────────────┤
│  Memory (Conversation History)                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│  │  User    │  │  Model   │  │  Tool    │              │
│  │ Messages │  │ Responses│  │  Results │              │
│  └──────────┘  └──────────┘  └──────────┘              │
├─────────────────────────────────────────────────────────┤
│  Event Hooks                                            │
│  ├─ OnTextChunk     → Stream text output                │
│  ├─ OnThinkingChunk → Stream reasoning                  │
│  └─ OnToolCall      → Monitor tool execution            │
├─────────────────────────────────────────────────────────┤
│  Tool Execution (Reflection-based)                     │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Standard Toolset (File I/O, Web Fetch, etc.)    │   │
│  │ Custom Tools (User-defined functions)          │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              LLM Interface (Provider-Agnostic)           │
├─────────────────────────────────────────────────────────┤
│  ┌───────────┐  ┌─────────┐  ┌─────────────┐           │
│  │OpenRouter │  │ Gemini  │  │   OpenAI    │           │
│  │ Provider  │  │Provider │  │  Provider   │           │
│  └───────────┘  └─────────┘  └─────────────┘           │
└─────────────────────────────────────────────────────────┘
```

## Troubleshooting

### Max Iterations Reached

If you see "Warning: Maximum iterations reached", the agent couldn't complete the task in the allotted iterations. Solutions:
- Increase `WithMaxIterations()`
- Simplify the task prompt
- Check if tools are working correctly

### Tool Execution Errors

Tools that return errors are logged but don't stop execution. To debug:
- Use `WithOnToolCall()` to see which tools are called
- Check tool error messages in the conversation history
- Ensure tool arguments have proper JSON tags

### Provider-Specific Issues

Different providers have different capabilities:
- Some providers don't support thinking/reasoning chunks
- Tool calling format may vary slightly
- Streaming behavior differs between providers

## Future Enhancements

Potential improvements for the butler package:

- Persistent memory storage (database, files)
- More providers (Anthropic, Cohere, etc.)
- Additional built-in tools (git operations, API calls)
- Tool result validation and error recovery
- Multi-agent collaboration patterns
- Structured output support
- Rate limiting and cost management
- Better tool discovery and auto-documentation

## License

This package is part of the arlocode project and follows the same license terms.
