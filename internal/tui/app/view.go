package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/kjk/flex"
	"github.com/mightymoud/arlocode/internal/tui/layers"
	"github.com/mightymoud/arlocode/internal/tui/themes"
)

func (m AppModel) RenderModal(canvas *layers.Canvas) {
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
}

func (m AppModel) RenderWelcomeScreen(canvas *layers.Canvas) {
	t := themes.Current

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

	m.WelcomeScreen.Input.Width = lipgloss.Width(inputBoxStyle.Render(m.WelcomeScreen.Input.View())) - 4
	m.WelcomeScreen.Input.PlaceholderStyle = lipgloss.NewStyle().Foreground(t.Overlay0()).Background(t.Surface0())
	m.WelcomeScreen.Input.TextStyle = lipgloss.NewStyle().Foreground(t.Text()).Background(t.Surface0())
	m.WelcomeScreen.Input.Cursor.Style = lipgloss.NewStyle().Foreground(t.Blue()).Background(t.Surface0())
	m.WelcomeScreen.Input.PromptStyle = lipgloss.NewStyle().Foreground(t.Blue()).Background(t.Surface0())
	m.WelcomeScreen.Input.Prompt = ""

	// Render main content elements
	title := titleStyle.Render("⚡ ArloCode")
	input := inputBoxStyle.Render(m.WelcomeScreen.Input.View())
	hint := hintStyle.Render("Ctrl+O to open modal • Esc to quit")

	sections := []string{title, input, hint}
	if m.ChatScreen.Conversation.AgentThinking {
		thinkingStyle := baseLayerStyle.
			Foreground(t.Yellow()).
			PaddingTop(1)
		thinkingText := thinkingStyle.Render(m.ChatScreen.Conversation.ThinkingBuffer + "█")
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
}

func (m AppModel) RenderChatScreen(canvas *layers.Canvas) {
	t := themes.Current
	// Base layer style (faint when modal is open)
	baseLayerStyle := lipgloss.NewStyle().Faint(m.showModal)

	root := flex.NewNode()
	root.StyleSetWidth(float32(m.width))
	root.StyleSetHeight(float32(m.height))
	root.StyleSetFlexDirection(flex.FlexDirectionRow)

	sideBar := flex.NewNode()
	sideBar.StyleSetWidth(30)
	sideBar.StyleSetHeight(float32(m.height))

	contentArea := flex.NewNode()
	contentArea.StyleSetFlexDirection(flex.FlexDirectionColumn)
	contentArea.StyleSetFlexGrow(1)
	contentArea.StyleSetHeight(float32(m.height))

	// add to content area
	chatContent := flex.NewNode()
	chatContent.StyleSetFlexGrow(1)

	inputArea := flex.NewNode()
	inputArea.StyleSetHeight(5)

	statusBar := flex.NewNode()
	statusBar.StyleSetHeight(1)

	contentArea.InsertChild(chatContent, 0)
	contentArea.InsertChild(inputArea, 1)
	contentArea.InsertChild(statusBar, 2)

	root.InsertChild(sideBar, 0)
	root.InsertChild(contentArea, 1)

	flex.CalculateLayout(root, float32(m.width), float32(m.height), flex.DirectionLTR)

	hintStyle := baseLayerStyle.
		Foreground(t.Overlay1()).
		Padding(0, 2)

	// Get calculated heights for main content layout
	contentAreaHeight := int(contentArea.LayoutGetHeight())
	chatContentHeight := int(chatContent.LayoutGetHeight())
	inputHeight := int(inputArea.LayoutGetHeight())
	statusBarHeight := int(statusBar.LayoutGetHeight())

	mainAreaWidth := int(contentArea.LayoutGetWidth())

	// Sidebar width
	sidebarWidth := int(sideBar.LayoutGetWidth())

	// Render main content elements
	chatDiv := lipgloss.NewStyle().
		Width(mainAreaWidth).
		Height(chatContentHeight).
		Background(t.Base())

	inputDiv := baseLayerStyle.
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(t.Blue()).
		BorderBackground(t.Surface0()).
		Background(t.Surface0()).
		Height(inputHeight).
		Width(mainAreaWidth).
		Background(t.Blue())
	hintDiv := hintStyle.
		Height(statusBarHeight).
		Width(mainAreaWidth).
		Background(t.Base())

	m.ChatScreen.Input.Width = lipgloss.Width(inputDiv.Render(m.ChatScreen.Input.View())) - 4
	m.ChatScreen.Input.PlaceholderStyle = lipgloss.NewStyle().Foreground(t.Overlay0()).Background(t.Surface0())
	m.ChatScreen.Input.TextStyle = lipgloss.NewStyle().Foreground(t.Text()).Background(t.Surface0())
	m.ChatScreen.Input.Cursor.Style = lipgloss.NewStyle().Foreground(t.Blue()).Background(t.Surface0())
	m.ChatScreen.Input.Prompt = ""

	sideBarDiv := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Surface1())).
		Width(sidebarWidth).
		Height(contentAreaHeight).
		Margin(0, 1)

	// Render conversation history
	var messageBoxes []string

	// Shared style helper for consistent message styling
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

	// Render all completed messages from conversation
	for _, msg := range m.ChatScreen.Conversation.Conversation {
		// Skip empty messages
		if msg.Content == "" {
			continue
		}
		var style lipgloss.Style
		switch msg.Type {
		case "user":
			style = userStyle
		case "agent":
			style = agentStyle
		case "thinking", "agent_thinking":
			style = thinkingStyle
		default:
			style = defaultStyle
		}
		messageBoxes = append(messageBoxes, style.Render(msg.Content))
	}

	// Render active thinking buffer (streaming)
	if m.ChatScreen.Conversation.AgentThinking && m.ChatScreen.Conversation.ThinkingBuffer != "" {
		messageBoxes = append(messageBoxes, thinkingStyle.Faint(true).Render(m.ChatScreen.Conversation.ThinkingBuffer+"█"))
	}

	// Render active text buffer (streaming)
	if m.ChatScreen.Conversation.TextBuffer != "" {
		messageBoxes = append(messageBoxes, agentStyle.Render(m.ChatScreen.Conversation.TextBuffer+"█"))
	}

	conversationContent := lipgloss.JoinVertical(lipgloss.Left, messageBoxes...)

	// Combine all content
	mainContent := lipgloss.NewStyle().
		Width(mainAreaWidth).
		Height(contentAreaHeight).
		Margin(0, 1).
		Background(t.Surface0()).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			chatDiv.Render(conversationContent),
			inputDiv.Render(m.ChatScreen.Input.View()),
			hintDiv.Render("Ctrl+O to open modal • Esc to quit"),
		))

	fullScreen := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		mainContent,
		sideBarDiv.Render(" Sidebar\n (Placeholder)"),
	)

	// Add base layer (Z=0)
	canvas.AddLayer(layers.NewLayer(fullScreen, 0).WithOffset(0, 0))
}

func (m AppModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Create canvas for layer composition
	canvas := layers.NewCanvas(m.width, m.height)

	// Render the current screen
	switch m.currentScreen {
	case ScreenWelcome:
		m.RenderWelcomeScreen(canvas)
	case ScreenChat:
		m.RenderChatScreen(canvas)
	}

	// Render modal if open
	if m.showModal {
		m.RenderModal(canvas)
	}

	// Render notifications if any
	if m.Notifications.HasActiveNotifications() {
		notifContent, notifX, notifY := m.Notifications.RenderWithPosition()
		canvas.AddLayer(layers.NewLayer(notifContent, 2).WithOffset(notifX, notifY))
	}

	return canvas.Render()
}
