package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mightymoud/arlocode/internal/tui/app/conversation"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
)

func NewAppModel() AppModel {
	// Initialize welcome screen input
	welcomeInput := textinput.New()
	welcomeInput.Placeholder = "Enter your prompt"
	welcomeInput.Focus()

	// Initialize chat screen input
	chatInput := textinput.New()
	chatInput.Placeholder = "Type a message..."
	chatInput.Focus()

	// Initialize modal input
	modalInput := textinput.New()
	modalInput.Placeholder = "Enter input..."

	// Initialize viewport for chat content
	vp := viewport.New(80, 20)
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 3
	vp.Height = 5

	return AppModel{
		currentScreen: ScreenChat,
		WelcomeScreen: WelcomeScreenModel{
			Input: welcomeInput,
		},
		ChatScreen: ChatScreenModel{
			Input:        chatInput,
			Conversation: conversation.NewConversationManager(),
			Viewport:     vp,
		},
		ModalInput:    modalInput,
		Notifications: notifications.NewNotificationManager(80, 24),
	}
}

func (m AppModel) Init() tea.Cmd {
	return textinput.Blink
}
