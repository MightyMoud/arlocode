package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	state "github.com/mightymoud/arlocode/internal/tui"
	"github.com/mightymoud/arlocode/internal/tui/themes"
)

var appState = state.Get()

// Update handles all messages and routes them to the appropriate screen
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.ChatScreen.Viewport, cmd = m.ChatScreen.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Notifications.UpdateScreenSize(msg.Width, msg.Height)
	}

	// if m.currentScreen == ScreenChat && !m.showModal && m.ChatScreen.Conversation.TextBuffer != "" {
	// 	// Check if it's a scroll event using IsWheel()
	// 	m.ChatScreen.Viewport, cmd = m.ChatScreen.Viewport.Update(msg)
	// 	cmds = append(cmds, cmd)
	// 	return m, tea.Batch(cmds...)
	// }

	// Handle mouse events first before textinput can process them
	if mouseMsg, ok := msg.(tea.MouseMsg); ok {
		// Forward mouse wheel events to viewport for scrolling
		if m.currentScreen == ScreenChat && !m.showModal {
			// Check if it's a scroll event using IsWheel()
			mouseEvent := tea.MouseEvent(mouseMsg)
			if mouseEvent.IsWheel() {
				m.ChatScreen.Viewport, cmd = m.ChatScreen.Viewport.Update(msg)
				// fmt.Print(m.ChatScreen.Viewport.Height, "vueport height\n")
				// fmt.Print(m.ChatScreen.Viewport.Width)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Notifications.UpdateScreenSize(msg.Width, msg.Height)

		// Update viewport size when window changes
		// Calculate chat content height (total height - input area - status bar - margins)
		chatContentHeight := msg.Height - 5 - 1 - 2
		sidebarWidth := 30
		mainAreaWidth := msg.Width - sidebarWidth - 2

		if chatContentHeight > 0 && mainAreaWidth > 0 {
			m.ChatScreen.Viewport.Width = mainAreaWidth
			m.ChatScreen.Viewport.Height = chatContentHeight
			// Update content after resizing
			m.updateViewportContent()
		}

	case tickMsg:
		// Update notification animations
		if m.Notifications.Update() {
			cmds = append(cmds, tickCmd())
		}
		// Update viewport content on tick to keep it current
		m.updateViewportContent()
		return m, tea.Batch(cmds...)

	case AgentTextChunkMsg:
		m.ChatScreen.Conversation.TextBuffer += string(msg)
		m.updateViewportContent()
		// Auto-scroll when streaming
		m.ChatScreen.Viewport.GotoBottom()
		return m, tea.Batch(cmds...)

	case AgentTextCompleteMsg:
		m.ChatScreen.Conversation.AddAgentMessage(m.ChatScreen.Conversation.TextBuffer)
		m.ChatScreen.Conversation.TextBuffer = ""
		m.updateViewportContent()
		return m, tea.Batch(cmds...)

	case AgentThinkingChunkMsg:
		m.ChatScreen.Conversation.AgentThinking = true
		m.ChatScreen.Conversation.ThinkingBuffer += string(msg)
		m.updateViewportContent()
		// Auto-scroll when thinking
		m.ChatScreen.Viewport.GotoBottom()
		return m, tea.Batch(cmds...)

	case AgentThinkingCompleteMsg:
		m.ChatScreen.Conversation.AddThinkingMessage(m.ChatScreen.Conversation.ThinkingBuffer)
		m.updateViewportContent()
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		// Handle global key bindings first
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.showModal {
				m.showModal = false
				m.ModalInput.Blur()
				m.focusCurrentScreenInput()
				// Get blink command for the focused input
				cmds = append(cmds, m.getCurrentScreenBlinkCmd())
				return m, tea.Batch(cmds...)
			}
			return m, nil
		case "ctrl+o":
			// Toggle modal
			m.showModal = !m.showModal
			if m.showModal {
				m.blurCurrentScreenInput()
				m.ModalInput.Focus()
				// Get blink command for modal input
				cmds = append(cmds, m.ModalInput.Cursor.BlinkCmd())
			} else {
				m.ModalInput.Blur()
				m.focusCurrentScreenInput()
				// Get blink command for the focused input
				cmds = append(cmds, m.getCurrentScreenBlinkCmd())
			}
			return m, tea.Batch(cmds...)
			// case "w":
			// 	// Show warning notification
			// 	if !m.showModal {
			// 		m.Notifications.PushWarning("Warning", "This is a warning notification!")
			// 		cmds = append(cmds, tickCmd())
			// 		return m, tea.Batch(cmds...)
			// 	}
		}

		// Handle modal input
		if m.showModal {
			if msg.String() == "enter" {
				// Close modal on enter
				m.showModal = false
				m.ModalInput.Blur()
				m.focusCurrentScreenInput()
				// Get blink command for the focused input
				cmds = append(cmds, m.getCurrentScreenBlinkCmd())
			}
			return m, tea.Batch(cmds...)
		}

		// Route key input to the current screen
		switch m.currentScreen {
		case ScreenWelcome:
			m, cmd = m.handleWelcomeScreenKeys(msg)
			cmds = append(cmds, cmd)
		case ScreenChat:
			m, cmd = m.handleChatScreenKeys(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Always update the focused textinput with all messages (for cursor blinking)
	if m.showModal {
		m.ModalInput, cmd = m.ModalInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		switch m.currentScreen {
		case ScreenWelcome:
			m.WelcomeScreen.Input, cmd = m.WelcomeScreen.Input.Update(msg)
			cmds = append(cmds, cmd)
		case ScreenChat:
			m.ChatScreen.Input, cmd = m.ChatScreen.Input.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleWelcomeScreenKeys handles key events for the welcome screen
func (m AppModel) handleWelcomeScreenKeys(msg tea.KeyMsg) (AppModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		value := m.WelcomeScreen.Input.Value()
		if value != "" {
			// Clear input and transition to chat screen
			m.WelcomeScreen.Input.SetValue("")
			m.currentScreen = ScreenChat
			m.WelcomeScreen.Input.Blur()
			m.ChatScreen.Input.Focus()
			// Add user message to conversation
			m.ChatScreen.Conversation.AddUserMessage(value)
			// Start the agent
			go appState.Agent().Run(context.Background(), value)
			return m, tickCmd()
		}
	}
	return m, nil
}

// handleChatScreenKeys handles key events for the chat screen
func (m AppModel) handleChatScreenKeys(msg tea.KeyMsg) (AppModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		value := m.ChatScreen.Input.Value()
		if value != "" {
			// Clear input after submission
			m.ChatScreen.Input.SetValue("")
			// Add user message to conversation
			m.ChatScreen.Conversation.AddUserMessage(value)
			// Start the agent
			go appState.Agent().Run(context.Background(), value)
			return m, tickCmd()
		}
	}
	return m, nil
}

// focusCurrentScreenInput focuses the input of the current screen
func (m *AppModel) focusCurrentScreenInput() {
	switch m.currentScreen {
	case ScreenWelcome:
		m.WelcomeScreen.Input.Focus()
	case ScreenChat:
		m.ChatScreen.Input.Focus()
	}
}

// blurCurrentScreenInput blurs the input of the current screen
func (m *AppModel) blurCurrentScreenInput() {
	switch m.currentScreen {
	case ScreenWelcome:
		m.WelcomeScreen.Input.Blur()
	case ScreenChat:
		m.ChatScreen.Input.Blur()
	}
}

// getCurrentScreenBlinkCmd returns the blink command for the current screen's input
func (m *AppModel) getCurrentScreenBlinkCmd() tea.Cmd {
	switch m.currentScreen {
	case ScreenWelcome:
		return m.WelcomeScreen.Input.Cursor.BlinkCmd()
	case ScreenChat:
		return m.ChatScreen.Input.Cursor.BlinkCmd()
	}
	return nil
}

// updateViewportContent builds and sets the conversation content on the viewport
func (m *AppModel) updateViewportContent() {
	if m.ChatScreen.Viewport.Width == 0 {
		return // Not initialized yet
	}

	t := themes.Current
	mainAreaWidth := m.ChatScreen.Viewport.Width

	// Create glamour renderer for agent messages
	glamourRenderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(mainAreaWidth-10),
	)

	// Shared style helper for consistent message styling
	baseLayerStyle := lipgloss.NewStyle().Faint(m.showModal)

	agentStyle := baseLayerStyle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Green()).
		Foreground(t.Text()).
		Padding(1, 1).
		MarginBottom(1).
		Width(mainAreaWidth - 4)

	thinkingStyle := baseLayerStyle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Yellow()).
		Foreground(t.Overlay1()).
		Background(t.Surface1()).
		Padding(1, 1).
		MarginBottom(1).
		Width(mainAreaWidth - 4)

	userStyle := baseLayerStyle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Blue()).
		Foreground(t.Text()).
		Padding(1, 1).
		MarginBottom(1).
		Width(mainAreaWidth - 4)

	defaultStyle := baseLayerStyle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Overlay0()).
		Foreground(t.Text()).
		Padding(1, 1).
		MarginBottom(1).
		Width(mainAreaWidth - 4)

	var messageBoxes []string

	// Render all completed messages from conversation
	for _, msg := range m.ChatScreen.Conversation.Conversation {
		if msg.Content == "" {
			continue
		}
		var style lipgloss.Style
		var content string
		switch msg.Type {
		case "user":
			style = userStyle
			content = msg.Content
		case "agent":
			style = agentStyle
			if glamourRenderer != nil {
				rendered, err := glamourRenderer.Render(msg.Content)
				if err == nil {
					content = rendered
				} else {
					content = msg.Content
				}
			} else {
				content = msg.Content
			}
		case "thinking", "agent_thinking":
			style = thinkingStyle
			content = msg.Content
		default:
			style = defaultStyle
			content = msg.Content
		}
		messageBoxes = append(messageBoxes, style.Render(content))
	}

	// Render active thinking buffer (streaming)
	if m.ChatScreen.Conversation.AgentThinking && m.ChatScreen.Conversation.ThinkingBuffer != "" {
		messageBoxes = append(messageBoxes, thinkingStyle.Faint(true).Render(m.ChatScreen.Conversation.ThinkingBuffer+"█"))
	}

	// Render active text buffer (streaming)
	if m.ChatScreen.Conversation.TextBuffer != "" {
		streamContent := m.ChatScreen.Conversation.TextBuffer
		if glamourRenderer != nil {
			rendered, err := glamourRenderer.Render(streamContent)
			if err == nil {
				streamContent = rendered + "█"
			} else {
				streamContent = streamContent + "█"
			}
		} else {
			streamContent = streamContent + "█"
		}
		messageBoxes = append(messageBoxes, agentStyle.Render(streamContent))
	}

	conversationContent := lipgloss.JoinVertical(lipgloss.Left, messageBoxes...)
	m.ChatScreen.Viewport.SetContent(conversationContent)
}
