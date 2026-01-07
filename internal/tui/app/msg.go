package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg is sent on each animation frame
type tickMsg time.Time

// tickCmd returns a command that ticks at 60fps for animations
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
