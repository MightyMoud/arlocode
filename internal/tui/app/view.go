package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kjk/flex"
	"github.com/mightymoud/arlocode/internal/tui/layers"
	"github.com/mightymoud/arlocode/internal/tui/themes"
)

// =========================================================================
// MODAL LAYER - Flexbox layout for modal overlay
// =========================================================================
func (m AppModel) RenderModal() (canvas layers.Canvas) {
	// Get theme styles
	t := themes.Current

	// Modal flexbox container
	modalRoot := flex.NewNode()
	modalRoot.StyleSetWidth(50)
	modalRoot.StyleSetFlexDirection(flex.FlexDirectionColumn)
	modalRoot.StyleSetPadding(flex.EdgeAll, 1)

	// Modal content nodes
	modalTitleNode := flex.NewNode()
	modalTitleNode.StyleSetHeight(2)

	modalInputNode := flex.NewNode()
	modalInputNode.StyleSetHeight(1)
	modalInputNode.StyleSetMargin(flex.EdgeVertical, 1)

	modalHintNode := flex.NewNode()
	modalHintNode.StyleSetHeight(1)

	modalRoot.InsertChild(modalTitleNode, 0)
	modalRoot.InsertChild(modalInputNode, 1)
	modalRoot.InsertChild(modalHintNode, 2)

	flex.CalculateLayout(modalRoot, 50, flex.Undefined, flex.DirectionLTR)

	// Modal styles with consistent background
	modalBg := t.Surface1()
	modalWidth := int(modalRoot.LayoutGetWidth())

	modalStyle := lipgloss.NewStyle().
		Background(modalBg).
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Pink()).
		BorderBackground(modalBg).
		Padding(1, 2).
		Width(modalWidth)

	modalTitleStyle := lipgloss.NewStyle().
		Background(modalBg).
		Bold(true).
		Foreground(t.Peach()).
		Width(modalWidth - 6).
		PaddingBottom(1)

	modalInputBoxStyle := lipgloss.NewStyle().
		Background(modalBg).
		Width(modalWidth - 6)

	modalHintStyle := lipgloss.NewStyle().
		Background(modalBg).
		Foreground(t.Overlay1()).
		Width(modalWidth - 6).
		PaddingTop(1)

	// Update modal input styles to match modal background
	m.ModalInput.TextStyle = lipgloss.NewStyle().Foreground(t.Text()).Background(modalBg)
	m.ModalInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(t.Overlay0()).Background(modalBg)
	m.ModalInput.Cursor.Style = lipgloss.NewStyle().Foreground(t.Rosewater()).Background(modalBg)

	modalContent := modalStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			modalTitleStyle.Render("Modal"),
			modalInputBoxStyle.Render(m.ModalInput.View()),
			modalHintStyle.Render("Enter to close • Esc to cancel"),
		),
	)

	// Center modal on screen
	modalX := (m.width - lipgloss.Width(modalContent)) / 2
	modalY := (m.height - lipgloss.Height(modalContent)) / 2

	// Add modal layer (Z=1, renders on top)
	canvas.AddLayer(layers.NewLayer(modalContent, 1).WithOffset(modalX, modalY))
	return
}

func (m AppModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Get theme styles
	t := themes.Current

	// Create canvas for layer composition
	canvas := layers.NewCanvas(m.width, m.height)

	// =========================================================================
	// BASE LAYER - Flexbox layout for main content
	// =========================================================================

	// Root flexbox container
	root := flex.NewNode()
	root.StyleSetWidth(float32(m.width))
	root.StyleSetHeight(float32(m.height))
	root.StyleSetFlexDirection(flex.FlexDirectionColumn)
	root.StyleSetJustifyContent(flex.JustifyCenter)
	root.StyleSetAlignItems(flex.AlignCenter)

	// Content container (centered vertically and horizontally)
	contentNode := flex.NewNode()
	contentNode.StyleSetFlexDirection(flex.FlexDirectionColumn)
	contentNode.StyleSetAlignItems(flex.AlignCenter)

	// Calculate layout
	root.InsertChild(contentNode, 0)

	agentThinkingNode := flex.NewNode()
	agentThinkingNode.StyleSetHeight(5)
	root.InsertChild(agentThinkingNode, 1)
	flex.CalculateLayout(root, float32(m.width), float32(m.height), flex.DirectionLTR)

	// Base layer style (faint when modal is open)
	baseLayerStyle := lipgloss.NewStyle().Faint(m.showModal)

	// Styles using theme colors
	titleStyle := baseLayerStyle.
		Bold(true).
		Foreground(t.Mauve()).
		PaddingBottom(2)

	inputBoxStyle := baseLayerStyle.
		Border(lipgloss.ThickBorder()).
		BorderTop(false).
		BorderBottom(false).
		BorderRight(false).
		Height(5).
		Background(t.Surface0()).
		BorderForeground(t.Blue()).
		Padding(0, 1).
		Width(60)

	hintStyle := baseLayerStyle.
		Foreground(t.Overlay1()).
		PaddingTop(2)

	// Render main content elements
	title := titleStyle.Render("⚡ ArloCode")
	input := inputBoxStyle.Render(m.MainInput.View())
	hint := hintStyle.Render("Ctrl+O to open modal • Esc to quit")

	sections := []string{title, input, hint}
	if m.Conversation.AgentThinking {
		thinkingStyle := baseLayerStyle.
			Foreground(t.Yellow()).
			PaddingTop(1)
		thinkingText := thinkingStyle.Render(m.Conversation.ThinkingBuffer + "█")
		sections = append(sections, thinkingText)
	}

	mainContent := lipgloss.JoinVertical(lipgloss.Center,
		sections...,
	)

	// Center in the full screen area using flex-calculated position
	contentX := (m.width - lipgloss.Width(mainContent)) / 2
	contentY := (m.height - lipgloss.Height(mainContent)) / 2

	// Add base layer (Z=0)
	canvas.AddLayer(layers.NewLayer(mainContent, 0).WithOffset(contentX, contentY))

	if m.showModal {
		m.RenderModal()
	}

	// =========================================================================
	// NOTIFICATIONS LAYER - Rendered on top of everything
	// =========================================================================

	if m.Notifications.HasActiveNotifications() {
		notifContent, notifX, notifY := m.Notifications.RenderWithPosition()
		canvas.AddLayer(layers.NewLayer(notifContent, 2).WithOffset(notifX, notifY))
	}

	return canvas.Render()
}
