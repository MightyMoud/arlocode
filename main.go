package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mightymoud/arlocode/internal/coding_agent"
	"github.com/mightymoud/arlocode/internal/tui/app"
	"github.com/mightymoud/arlocode/internal/tui/app/conversation"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
	"github.com/mightymoud/arlocode/internal/tui/themes"
)

func main() {
	// Get theme for input styling
	t := themes.Current

	// Create main input
	mainInput := textinput.New()
	mainInput.Placeholder = "What would you like to do?"
	mainInput.Width = 56
	mainInput.CharLimit = 200
	mainInput.TextStyle = lipgloss.NewStyle().Foreground(t.Text())
	mainInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(t.Overlay0()).Background(t.Surface0())
	mainInput.Cursor.Style = lipgloss.NewStyle().Foreground(t.Rosewater())
	mainInput.Focus()

	// Create modal input
	modalInput := textinput.New()
	modalInput.Placeholder = "Enter value..."
	modalInput.Width = 40
	modalInput.CharLimit = 100
	modalInput.TextStyle = lipgloss.NewStyle().Foreground(t.Text())
	modalInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(t.Overlay0())
	modalInput.Cursor.Style = lipgloss.NewStyle().Foreground(t.Rosewater())
	shared := &app.SharedState{}

	m := &app.AppModel{
		MainInput:     mainInput,
		ModalInput:    modalInput,
		Notifications: notifications.NewNotificationManager(80, 24),
		Conversation:  &conversation.ConversationManager{},
		Shared:        shared,
	}

	codingAgent := coding_agent.Agent.WithMaxIterations(10).
		WithOnThinkingChunk(func(s string) {
			shared.Program.Send(app.AgentThinkingChunkMsg(s))
		}).
		WithOnThinkingComplete(func() {
			shared.Program.Send(app.AgentThinkingCompleteMsg(""))
		})

	shared.Agent = codingAgent
	p := tea.NewProgram(m, tea.WithAltScreen())
	shared.Program = p
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

}
