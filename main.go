package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kjk/flex"
	"github.com/mightymoud/arlocode/internal/themes"
	"github.com/mightymoud/arlocode/internal/tui/layers"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
)

type model struct {
	width         int
	height        int
	showModal     bool
	mainInput     textinput.Model
	modalInput    textinput.Model
	notifications *notifications.NotificationManager
}

// tickMsg is sent on each animation frame
type tickMsg time.Time

// tickCmd returns a command that ticks at 60fps for animations
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.notifications.UpdateScreenSize(msg.Width, msg.Height)

	case tickMsg:
		// Update notification animations
		if m.notifications.Update() {
			cmds = append(cmds, tickCmd())
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.showModal {
				m.showModal = false
				m.modalInput.Blur()
				m.mainInput.Focus()
				return m, textinput.Blink
			}
			return m, tea.Quit
		case "enter":
			if m.showModal {
				// Close modal on enter
				m.showModal = false
				m.modalInput.Blur()
				m.mainInput.Focus()
				return m, textinput.Blink
			}
		case "ctrl+o":
			// Toggle modal
			m.showModal = !m.showModal
			if m.showModal {
				m.mainInput.Blur()
				m.modalInput.Focus()
				return m, textinput.Blink
			} else {
				m.modalInput.Blur()
				m.mainInput.Focus()
				return m, textinput.Blink
			}
		case "w":
			// Show warning notification
			if !m.showModal {
				m.notifications.PushWarning("Warning", "This is a warning notification!")
				return m, tickCmd()
			}
		}

		// Route input to the focused component
		if m.showModal {
			m.modalInput, cmd = m.modalInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.mainInput, cmd = m.mainInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
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
	input := inputBoxStyle.Render(m.mainInput.View())
	hint := hintStyle.Render("Ctrl+O to open modal • Esc to quit")

	mainContent := lipgloss.JoinVertical(lipgloss.Center,
		title,
		input,
		hint,
	)

	// Center in the full screen area using flex-calculated position
	contentX := (m.width - lipgloss.Width(mainContent)) / 2
	contentY := (m.height - lipgloss.Height(mainContent)) / 2

	// Add base layer (Z=0)
	canvas.AddLayer(layers.NewLayer(mainContent, 0).WithOffset(contentX, contentY))

	// =========================================================================
	// MODAL LAYER - Flexbox layout for modal overlay
	// =========================================================================

	if m.showModal {
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
		m.modalInput.TextStyle = lipgloss.NewStyle().Foreground(t.Text()).Background(modalBg)
		m.modalInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(t.Overlay0()).Background(modalBg)
		m.modalInput.Cursor.Style = lipgloss.NewStyle().Foreground(t.Rosewater()).Background(modalBg)

		modalContent := modalStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				modalTitleStyle.Render("Modal"),
				modalInputBoxStyle.Render(m.modalInput.View()),
				modalHintStyle.Render("Enter to close • Esc to cancel"),
			),
		)

		// Center modal on screen
		modalX := (m.width - lipgloss.Width(modalContent)) / 2
		modalY := (m.height - lipgloss.Height(modalContent)) / 2

		// Add modal layer (Z=1, renders on top)
		canvas.AddLayer(layers.NewLayer(modalContent, 1).WithOffset(modalX, modalY))
	}

	// =========================================================================
	// NOTIFICATIONS LAYER - Rendered on top of everything
	// =========================================================================

	if m.notifications.HasActiveNotifications() {
		notifContent, notifX, notifY := m.notifications.RenderWithPosition()
		canvas.AddLayer(layers.NewLayer(notifContent, 2).WithOffset(notifX, notifY))
	}

	return canvas.Render()
}

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

	m := model{
		mainInput:     mainInput,
		modalInput:    modalInput,
		notifications: notifications.NewNotificationManager(80, 24),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
