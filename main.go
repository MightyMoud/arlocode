package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kjk/flex"

	"github.com/mightymoud/arlocode/internal/layers"
	"github.com/mightymoud/arlocode/internal/notifications"
)

type model struct {
	width         int
	height        int
	showModal     bool
	activeLayer   int // 0 = background focused, 1 = modal focused
	notifications *notifications.NotificationManager
}

// tickMsg is sent on each animation frame for notifications
type tickMsg time.Time

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tick())
}

// tick returns a command that sends a tickMsg after a short delay
func tick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.notifications != nil {
			m.notifications.UpdateScreenSize(msg.Width, msg.Height)
		}
	case tickMsg:
		// Update notification animations
		if m.notifications != nil {
			m.notifications.Update()
		}
		return m, tick()
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.showModal = !m.showModal
		case "tab":
			m.activeLayer = (m.activeLayer + 1) % 2
		case "i":
			if m.notifications != nil {
				m.notifications.PushInfo("Info", "This is an informational message")
			}
		case "s":
			if m.notifications != nil {
				m.notifications.PushSuccess("Success!", "The operation completed successfully")
			}
		case "w":
			if m.notifications != nil {
				m.notifications.PushWarning("Warning", "Something might need your attention")
			}
		case "e":
			if m.notifications != nil {
				m.notifications.PushError("Error", "Something went wrong!")
			}
		case "d":
			if m.notifications != nil {
				m.notifications.DismissAll()
			}
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// =========================================================================
	// FLEXBOX LAYOUT SETUP
	// =========================================================================

	// 1. Create the Root Node
	root := flex.NewNode()
	root.StyleSetWidth(float32(m.width))
	root.StyleSetHeight(float32(m.height))
	root.StyleSetFlexDirection(flex.FlexDirectionColumn)

	// 2. Define Children Nodes
	header := flex.NewNode()
	header.StyleSetHeight(3) // 3 rows high for header
	header.StyleSetFlexShrink(0)

	mainContent := flex.NewNode()
	mainContent.StyleSetFlexGrow(1) // Fills remaining space

	footer := flex.NewNode()
	footer.StyleSetHeight(1) // 1 row for footer/status bar
	footer.StyleSetFlexShrink(0)

	// Insert children into root
	root.InsertChild(header, 0)
	root.InsertChild(mainContent, 1)
	root.InsertChild(footer, 2)

	// 3. Calculate Layout
	flex.CalculateLayout(root, float32(m.width), float32(m.height), flex.DirectionLTR)

	// =========================================================================
	// RENDER EACH PANEL USING CALCULATED FLEX DIMENSIONS
	// =========================================================================

	// Create a canvas for layer composition
	canvas := layers.NewCanvas(m.width, m.height)

	// --- HEADER PANEL ---
	headerWidth := int(header.LayoutGetWidth())
	headerHeight := int(header.LayoutGetHeight())
	headerY := int(header.LayoutGetTop())

	headerStyle := lipgloss.NewStyle().
		Width(headerWidth).
		Height(headerHeight).
		Background(lipgloss.Color("#0f3460")).
		Foreground(lipgloss.Color("#e0e0e0")).
		Padding(0, 2).
		Align(lipgloss.Left, lipgloss.Center)

	headerContent := headerStyle.Render(
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e94560")).Render("⚡ ArloCode") +
			"  │  Flexbox Layout Demo  │  SPACE toggle modal, i/s/w/e notifications, Q quit",
	)
	canvas.AddLayer(layers.NewLayer(headerContent, 0).WithOffset(0, headerY))

	// --- MAIN CONTENT PANEL ---
	mainWidth := int(mainContent.LayoutGetWidth())
	mainHeight := int(mainContent.LayoutGetHeight())
	mainY := int(mainContent.LayoutGetTop())

	// Background pattern for main area
	bgStyle := lipgloss.NewStyle().
		Width(mainWidth).
		Height(mainHeight).
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(lipgloss.Color("#4a4a6a"))

	pattern := ""
	for y := 0; y < mainHeight; y++ {
		row := ""
		for x := 0; x < mainWidth; x++ {
			if (x+y)%4 == 0 {
				row += "·"
			} else {
				row += " "
			}
		}
		pattern += row
		if y < mainHeight-1 {
			pattern += "\n"
		}
	}
	backgroundContent := bgStyle.Render(pattern)
	canvas.AddLayer(layers.NewLayer(backgroundContent, 1).WithOffset(0, mainY))

	// --- CENTERED CONTENT PANEL (inside main) ---
	contentPanelWidth := min(60, mainWidth-10)
	contentPanelHeight := min(15, mainHeight-4)

	contentPanelStyle := lipgloss.NewStyle().
		Width(contentPanelWidth).
		Height(contentPanelHeight).
		Background(lipgloss.Color("#16213e")).
		Foreground(lipgloss.Color("#e0e0e0")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#0f3460")).
		Padding(1, 2).
		Align(lipgloss.Center, lipgloss.Center)

	contentPanel := contentPanelStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e94560")).Render("[Flexbox Layout Demo]"),
			"",
			"This panel uses github.com/kjk/flex",
			"for flexbox-based terminal layout!",
			"",
			fmt.Sprintf("Header: %dw × %dh", headerWidth, headerHeight),
			fmt.Sprintf("Main:   %dw × %dh", mainWidth, mainHeight),
			fmt.Sprintf("Footer: %dw × %dh", int(footer.LayoutGetWidth()), int(footer.LayoutGetHeight())),
			"",
			lipgloss.NewStyle().Faint(true).Render("Press SPACE to toggle modal"),
		),
	)

	// Center the content panel within main area
	contentPanelX := (mainWidth - lipgloss.Width(contentPanel)) / 2
	contentPanelY := mainY + (mainHeight-lipgloss.Height(contentPanel))/2
	canvas.AddLayer(layers.NewLayer(contentPanel, 2).WithOffset(contentPanelX, contentPanelY))

	// --- MODAL OVERLAY (conditionally shown) ---
	if m.showModal {
		modalWidth := 40
		modalHeight := 10

		modalStyle := lipgloss.NewStyle().
			Width(modalWidth).
			Height(modalHeight).
			Background(lipgloss.Color("#2d132c")).
			Foreground(lipgloss.Color("#ffffff")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#ee4540")).
			Padding(1, 2).
			Align(lipgloss.Center, lipgloss.Center)

		modalContent := modalStyle.Render(
			lipgloss.JoinVertical(lipgloss.Center,
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffd460")).Render("! Modal Dialog !"),
				"",
				"Modal overlays flexbox layout",
				"at Z-index = 5",
				"",
				lipgloss.NewStyle().Italic(true).Render("Press SPACE to close"),
			),
		)

		// Center the modal on screen
		modalX := (m.width - lipgloss.Width(modalContent)) / 2
		modalY := (m.height - lipgloss.Height(modalContent)) / 2
		canvas.AddLayer(layers.NewLayer(modalContent, 5).WithOffset(modalX, modalY))
	}

	// --- FOOTER/STATUS BAR ---
	footerWidth := int(footer.LayoutGetWidth())
	footerHeight := int(footer.LayoutGetHeight())
	footerY := int(footer.LayoutGetTop())

	footerStyle := lipgloss.NewStyle().
		Width(footerWidth).
		Height(footerHeight).
		Background(lipgloss.Color("#0f3460")).
		Foreground(lipgloss.Color("#e0e0e0")).
		Padding(0, 1)

	modalStatus := "Modal: OFF"
	if m.showModal {
		modalStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render("Modal: ON")
	}

	footerContent := footerStyle.Render(
		fmt.Sprintf("Size: %dx%d │ %s │ Flexbox: Header→Main→Footer",
			m.width, m.height, modalStatus),
	)
	canvas.AddLayer(layers.NewLayer(footerContent, 10).WithOffset(0, footerY))

	// --- NOTIFICATION OVERLAY ---
	if m.notifications != nil {
		notifContent, notifX, notifY := m.notifications.RenderWithPosition()
		if notifContent != "" {
			canvas.AddLayer(layers.NewLayer(notifContent, 100).WithOffset(notifX, notifY))
		}
	}

	return canvas.RenderWithLipgloss()
}

func main() {
	// Initialize the model with notification manager
	m := model{
		notifications: notifications.NewNotificationManager(80, 24),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// geminiProvider := gemini.New(ctx)
	// model2 := geminiProvider.Model(ctx, "gemini-3-flash-preview")
	// geminiBasedAgent := butler.NewAgent(model2)
	// geminiBasedAgent.Run(ctx, prompt)
}

// prompt := "Some tool tests are missing in the tools package. Please add them."
// 	// prompt := "Add a README.md file to the internal package named butler that will explain what the butler package is doing and how to use it."
// 	// prompt += "In the file Agent.go, I added a new confirmation step before executing any tool call. This step prompts the user to confirm whether they want to proceed with the tool call or not. If the user declines, the tool call is skipped, and a message is displayed indicating that the tool call was cancelled by the user. The logic is not working. Fix it"
// 	// prompt := "Add a github action workflow file that will run tests for all the packages on every push to the repository."
// 	// prompt := "Read all the code in the butler package and tell me what it does:"
// 	// prompt := "I have added a new function called runCommand - Implement the function properly and make sure to pipe the output back to the LLM from both std and error streams"
// 	// prompt := "run the git status command and tell me what you see in the output"

// 	openrouterProvider := openrouter.New(ctx)
// 	model := openrouterProvider.Model(ctx, "z-ai/glm-4.7")

// 	openrouterBasedAgent := agent.NewAgent(model).
// 		WithMaxIterations(50).
// 		WithOnThinkingChunk(func(chunk string) {
// 			color.RGB(255, 128, 0).Printf("%s", chunk)

// 		}).
// 		WithOnTextChunk(func(chunk string) {
// 			fmt.Printf("%s", chunk)

// 		}).
// 		WithOnToolCall(func(t tools.ToolCall) {
// 			color.Blue("\n[Tool Call] %s - with Arguments: %+v", t.FunctionName, t.Arguments)
// 		})
// 	openrouterBasedAgent.Run(ctx, prompt)
