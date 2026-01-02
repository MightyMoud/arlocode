package butler

import (
	"context"

	"github.com/mightymoud/sidekick-agent/internal/butler/common"
)

type Agent struct {
	llm    LLM
	memory []common.MemoryEntry
}

func NewAgent(l LLM) *Agent {
	return &Agent{llm: l, memory: []common.MemoryEntry{}}
}

// Mock for Memory stuff later this is where Agent will use it
func (a *Agent) AddMemoryEntry(entry common.MemoryEntry) {
	a.memory = append(a.memory, entry)
}

// Memory stuff later
func (a *Agent) GetMemory() []common.MemoryEntry {
	return a.memory
}

// Mods
func (a *Agent) WithMemory(memory []common.MemoryEntry) *Agent {
	a.memory = memory
	return a
}

func (a *Agent) Run(ctx context.Context, prompt string) error {
	initMessage := common.MemoryEntry{Message: prompt, Role: "user"}
	a.AddMemoryEntry(initMessage)
	return a.llm.Generate(ctx, a.memory)
}
