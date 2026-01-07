package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mightymoud/arlocode/internal/butler/agent"
	"github.com/mightymoud/arlocode/internal/tui/app/conversation"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
)

type SharedState struct {
	Program *tea.Program
	Agent   *agent.Agent // or whatever your agent type is
}

type AppModel struct {
	width         int
	height        int
	showModal     bool
	MainInput     textinput.Model
	ModalInput    textinput.Model
	Notifications *notifications.NotificationManager

	// Coding agent related fields
	Conversation *conversation.ConversationManager

	Shared *SharedState
}
