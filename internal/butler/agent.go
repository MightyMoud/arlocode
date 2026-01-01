package butler

import (
	"context"

	"github.com/mightymoud/sidekick-agent/internal/butler/common"
)

type LLM interface {
	AddHistory(entry common.History)
	History() []common.History
	Stream(ctx context.Context) error
	Generate(ctx context.Context) error
}

type Agent struct {
	llm LLM
}

func NewAgent(l LLM) *Agent {
	return &Agent{llm: l}
}

func (a *Agent) Run(ctx context.Context, prompt string) error {
	initMessage := common.History{Message: prompt, Role: "user"}
	a.llm.AddHistory(initMessage)
	return a.llm.Stream(ctx)
}
