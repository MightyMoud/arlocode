package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mightymoud/arlocode/internal/tui/app/conversation"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
)

type MainModel struct {
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
}
