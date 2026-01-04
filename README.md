# Arlo code Project - WIP

This is an ambitious project aiming to make a coding agent that is:
- Long running capable in one command -> full with infra integration
- Cheaper to run same models vs other agents -> I think there is waste in how models are used
- Performs better on large codebases -> better code semantic understanding
- Support the vision of an agent companion and live team coding sessions
- Support hats for different types of work including teach, plan, long running, test, design and write

## Butler Package

**Butler** is a flexible, provider-agnostic AI agent framework that enables LLM-powered automation with tool calling capabilities. It supports multiple providers (OpenAI, Gemini, OpenRouter), maintains conversation memory, and executes iterative tasks until completion.

**Key Features:**
- 🔌 Swap between LLM providers without changing code
- 🛠️ Execute custom tools and file operations
- 🧠 Track and maintain conversation history
- 📡 Stream responses in real-time

**Simple Example:**
```go
package main

import (
    "context"
    "github.com/mightymoud/arlocode/internal/butler/agent"
    "github.com/mightymoud/arlocode/internal/butler/providers/openrouter"
)

func main() {
    ctx := context.Background()
    provider := openrouter.New(ctx)
    model := provider.Model(ctx, "anthropic/claude-sonnet-4.5")
    
    agent := agent.NewAgent(model)
    agent.Run(ctx, "Read the main.go file and explain what it does")
}
```

# But why?
> When the stars are within reach it's foolish to not aim for the moon.