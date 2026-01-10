package app

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mightymoud/arlocode/internal/tui/themes"
)

// buildConversationContent builds the styled conversation content string for the viewport.
// This is called from View() to keep styling/presentation logic separate from Update().
func (m AppModel) buildConversationContent(mainAreaWidth int, baseLayerStyle lipgloss.Style) string {
	t := themes.Current

	// Create glamour renderer for agent messages
	glamourRenderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(mainAreaWidth-10),
	)

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

	return lipgloss.JoinVertical(lipgloss.Left, messageBoxes...)
}
