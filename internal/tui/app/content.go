package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mightymoud/arlocode/internal/tui/themes"
)

// buildConversationContent builds the styled conversation content string for the viewport.
// This is called from View() to keep styling/presentation logic separate from Update().
func (m AppModel) buildConversationContent(mainAreaWidth int, baseLayerStyle lipgloss.Style) string {
	t := themes.Current

	// Create glamour renderer for agent messages with themed styles
	glamourRenderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(themes.GlamourStyle()),
		glamour.WithWordWrap(mainAreaWidth-10),
		glamour.WithPreservedNewLines(),
	)

	agentStyle := baseLayerStyle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Green()).
		Foreground(t.Text()).
		PaddingBottom(1).
		PaddingLeft(1).
		MarginBottom(1)

	activeThinkingStyle := baseLayerStyle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Yellow()).
		Foreground(t.Text()).
		Padding(1, 1).
		MarginBottom(1).
		Faint(true).
		Width(mainAreaWidth - 2)

	finishedThinkingStyle := activeThinkingStyle.
		UnsetPadding().
		PaddingLeft(1)

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
					content = strings.TrimRight(rendered, "\n")
					// content = rendered
				} else {
					content = msg.Content
				}
			} else {
				content = msg.Content
			}
		case "thinking":
			if msg.Status == "complete" {
				style = finishedThinkingStyle
				// Show duration instead of content
				duration := msg.EndTime.Sub(msg.StartTime)
				seconds := int(duration.Seconds())
				if seconds < 1 {
					content = "thought for < 1s"
				} else {
					content = fmt.Sprintf("thought for %ds", seconds)
				}
			} else {
				style = activeThinkingStyle
				// In progress - show the actual content (streaming tokens)
				content = msg.Content + "█"
			}
		default:
			style = defaultStyle
			content = msg.Content
		}
		if content != "" {
			messageBoxes = append(messageBoxes, style.Render(content))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, messageBoxes...)
}
