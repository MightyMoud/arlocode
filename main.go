package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mightymoud/arlocode/internal/coding_agent"
	state "github.com/mightymoud/arlocode/internal/tui"
	"github.com/mightymoud/arlocode/internal/tui/app"
)

func main() {
	appState := state.Get()

	// Create the app model using the new constructor
	m := app.NewAppModel()

	codingAgent := coding_agent.Agent.WithMaxIterations(10).
		WithOnThinkingChunk(func(s string) {
			appState.Program().Send(app.AgentThinkingChunkMsg(s))
		}).
		WithOnThinkingComplete(func() {
			appState.Program().Send(app.AgentThinkingCompleteMsg(""))
		}).
		WithOnTextChunk(func(s string) {
			appState.Program().Send(app.AgentTextChunkMsg(s))
		}).
		WithOnStreamComplete(func() {
			appState.Program().Send(app.AgentTextCompleteMsg(""))
		})

	appState.SetAgent(codingAgent)

	p := tea.NewProgram(m, tea.WithAltScreen())
	appState.SetProgram(p)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
