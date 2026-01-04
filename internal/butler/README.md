# Butler Package

Butler is a flexible, provider-agnostic AI agent framework for Go that enables LLM-powered automation with tool calling capabilities.

## Overview

The butler package provides an agent-based architecture that can work with multiple LLM providers (OpenAI, Gemini, OpenRouter) and execute tools/functions based on the agent's decisions. It maintains conversation memory and supports iterative reasoning loops until a task is complete.

## Core Components

### Agent

The `Agent` is the central orchestrator that:
- Manages conversation memory
- Coordinates with LLM providers
- Handles tool execution
- Controls iteration loops to prevent infinite execution

### LLM Interface

A provider-agnostic interface that allows different LLM providers to be used interchangeably:
- `Stream()`: Streams responses from the LLM with tool calling support
- `Generate()`: Generates responses from the LLM

### Memory

The memory system tracks conversation history including:
- User messages
- Model responses
- Tool calls and their results
- Role-based message organization

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
    
    "github.com/mightymoud/arlocode/internal/butler"
    "github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
)

func main() {
    ctx := context.Background()
    
    // Create a provider and model
    provider := openrouter.New(ctx)
    model := provider.Model(ctx, "anthropic/claude-sonnet-4.5")
    
    // Create and configure an agent
    agent := butler.NewAgent(model).WithMaxIterations(50)
    
    // Run the agent with a prompt
    prompt := "Read the main.go file and explain what it does"
    agent.Run(ctx, prompt)
}
```

### Customizing the Agent

#### Using Custom Tools

```go
// Create an agent with custom tools
customTools := []tools.Tool{
    tools.NewButlerTool("my_tool", "Description", myToolFunc),
}

agent := butler.NewAgent(model).WitTools(customTools)
```

#### Using No Tools

```go
// Create an agent without any tools (chat-only mode)
agent := butler.NewAgent(model).WithNoTools()
```

#### Setting Max Iterations

```go
// Limit the agent to 20 iterations
agent := butler.NewAgent(model).WithMaxIterations(20)
```

#### Using Existing Memory

```go
// Continue from previous conversation
existingMemory := []memory.MemoryEntry{
    {Role: "user", Message: "Hello"},
    {Role: "model", Message: "Hi! How can I help?"},
}

agent := butler.NewAgent(model).WithMemory(existingMemory)
```

## Supported Providers

### OpenRouter

```go
provider := openrouter.New(ctx)
model := provider.Model(ctx, "anthropic/claude-sonnet-4.5")
```

### Gemini

```go
provider := gemini.New(ctx)
model := provider.Model(ctx, "gemini-3-flash-preview")
```

### OpenAI

```go
provider := openai.New(ctx)
model := provider.Model(ctx, "gpt-4")
```

## Creating Custom Tools

Tools are created using the `NewButlerTool` function which uses reflection to handle any function with a struct argument:

```go
// Define your tool's arguments
type myToolArgs struct {
    Input string `json:"input" jsonschema:"Description of the input parameter"`
}

// Define your tool function
func myToolFunc(args myToolArgs) (string, error) {
    // Tool implementation
    return "result", nil
}

// Create the tool
myTool := tools.NewButlerTool(
    "my_tool_name",
    "Description of what the tool does",
    myToolFunc,
)
```

## How It Works

1. **Initialization**: Create an agent with an LLM provider and optional configuration
2. **Prompt**: Submit a user prompt via `agent.Run()`
3. **Iteration Loop**: The agent enters a loop where it:
   - Sends the current memory (conversation history) to the LLM
   - Receives a response (text and/or tool calls)
   - Executes any requested tool calls
   - Adds results to memory
   - Repeats until no more tool calls are needed or max iterations is reached
4. **Completion**: The loop exits when the LLM decides the task is complete

## Memory Management

Retrieve conversation history at any time:

```go
history := agent.GetMemory()
for _, entry := range history {
    fmt.Printf("%s: %s\n", entry.Role, entry.Message)
}
```

Add custom memory entries:

```go
agent.AddMemoryEntry(memory.MemoryEntry{
    Role:    "user",
    Message: "Custom message",
})
```

## Best Practices

1. **Set Appropriate Max Iterations**: Default is 10. Increase for complex tasks, decrease for simple ones
2. **Choose the Right Provider**: Different providers have different strengths and pricing
3. **Use Specific Tools**: Only include tools relevant to your task for better performance
4. **Monitor Tool Calls**: The agent logs tool calls in blue for easy debugging
5. **Handle Errors**: The `Run()` method returns an error that should be checked

## Architecture

```
Agent
├── LLM Interface (provider-agnostic)
│   ├── OpenRouter Provider
│   ├── Gemini Provider
│   └── OpenAI Provider
├── Memory System
│   └── MemoryEntry (messages, roles, tool calls)
└── Tools System
    ├── Tool Definition (name, description, handler)
    └── Standard Toolset (file operations, web fetch, etc.)
```

## Future Enhancements

- Persistent memory storage
- More providers (Anthropic, Cohere, etc.)
- Additional built-in tools
- Streaming responses to user
- Tool result validation
- Multi-agent collaboration
