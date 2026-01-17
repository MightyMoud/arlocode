package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mightymoud/arlocode/internal/butler/tools"
)

// tickMsg is sent on each animation frame
type tickMsg time.Time

// tickCmd returns a command that ticks at 60fps for animations
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickLoadingMsg time.Time

func TickLoading() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return tickLoadingMsg(t)
	})
}

type AgentTextChunkMsg string
type AgentTextCompleteMsg string

type AgentThinkingChunkMsg string
type AgentThinkingCompleteMsg string

type ToolCallMsg tools.ToolCall
