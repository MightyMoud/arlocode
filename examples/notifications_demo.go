package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mightymoud/arlocode/internal/tui/layers"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
)

// tickMsg is sent on each animation frame
type tickMsg time.Time

// Model holds the application state
type model struct {
	width         int
	height        int
	notifications *notifications.NotificationManager
	canvas        *layers.Canvas
}

func initialModel() model {
	return model{
		width:         80,
		height:        24,
		notifications: notifications.NewNotificationManager(80, 24),
	}
}

func (m model) Init() tea.Cmd {
	// Start with a welcome notification
	m.notifications.PushInfo("Welcome!", "Press keys to trigger different notifications")
	return tick()
}

// tick returns a command that sends a tickMsg after a short delay
func tick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "i":
			m.notifications.PushInfo("Info", "This is an informational message")
			return m, nil

		case "s":
			m.notifications.PushSuccess("Success!", "The operation completed successfully")
			return m, nil

		case "w":
			m.notifications.PushWarning("Warning", "Something might need your attention")
			return m, nil

		case "e":
			m.notifications.PushError("Error", "Something went wrong! Please check the logs.")
			return m, nil

		case "d":
			m.notifications.DismissAll()
			return m, nil

		case "m":
			// Push multiple notifications rapidly
			m.notifications.PushInfo("First", "This is the first notification")
			m.notifications.PushSuccess("Second", "This is the second one")
			m.notifications.PushWarning("Third", "And here's a third!")
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.notifications.UpdateScreenSize(msg.Width, msg.Height)
		return m, nil

	case tickMsg:
		// Update notification animations
		m.notifications.Update()
		// Continue ticking if there are active notifications
		if m.notifications.HasActiveNotifications() {
			return m, tick()
		}
		// Even with no notifications, keep ticking to catch new ones
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	// Create the main content (background)
	mainContent := m.renderMainContent()

	// Create canvas for layering
	canvas := layers.NewCanvas(m.width, m.height)

	// Add main content as base layer
	canvas.AddLayer(layers.NewLayer(mainContent, 0))

	// Add notifications as overlay layer
	notifContent, x, y := m.notifications.RenderWithPosition()
	if notifContent != "" {
		canvas.AddLayer(layers.NewLayer(notifContent, 100).WithOffset(x, y))
	}

	return canvas.Render()
}

func (m model) renderMainContent() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true)

	title := titleStyle.Render("🔔 Notification System Demo")

	help := lipgloss.JoinVertical(lipgloss.Left,
		"",
		helpStyle.Render("Press keys to trigger notifications:"),
		"",
		fmt.Sprintf("  %s  Info notification", keyStyle.Render("i")),
		fmt.Sprintf("  %s  Success notification", keyStyle.Render("s")),
		fmt.Sprintf("  %s  Warning notification", keyStyle.Render("w")),
		fmt.Sprintf("  %s  Error notification", keyStyle.Render("e")),
		fmt.Sprintf("  %s  Multiple notifications", keyStyle.Render("m")),
		fmt.Sprintf("  %s  Dismiss all", keyStyle.Render("d")),
		"",
		fmt.Sprintf("  %s  Quit", keyStyle.Render("q")),
		"",
		helpStyle.Render(fmt.Sprintf("Active notifications: %d", m.notifications.Count())),
	)

	content := lipgloss.JoinVertical(lipgloss.Left, title, help)

	// Center the content
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
