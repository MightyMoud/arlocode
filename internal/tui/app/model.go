package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/mightymoud/arlocode/internal/tui/app/conversation"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
)

// Screen represents the different screens in the application
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenChat
)

// WelcomeScreenModel represents the welcome screen state
type WelcomeScreenModel struct {
	Input textinput.Model
}

// ChatScreenModel represents the chat screen state
type ChatScreenModel struct {
	Input        textinput.Model
	Conversation *conversation.ConversationManager
	Viewport     viewport.Model
}

// AppModel is the main application model that manages screen routing
type AppModel struct {
	width         int
	height        int
	showModal     bool
	currentScreen Screen
	ModalInput    textinput.Model
	Notifications *notifications.NotificationManager

	// Screen models
	WelcomeScreen WelcomeScreenModel
	ChatScreen    ChatScreenModel
}
