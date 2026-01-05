package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mightymoud/arlocode/internal/layers"
)

type model struct {
	width       int
	height      int
	showModal   bool
	activeLayer int // 0 = background focused, 1 = modal focused
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.showModal = !m.showModal
		case "tab":
			m.activeLayer = (m.activeLayer + 1) % 2
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

	// Create a canvas for layer composition
	canvas := layers.NewCanvas(m.width, m.height-1)

	// =========================================================================
	// LAYER 1 (Z=0): Background - A patterned background layer
	// =========================================================================
	bgStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 1).
		Background(lipgloss.Color("#1a1a2e")).
		Foreground(lipgloss.Color("#4a4a6a"))

	// Create a dotted pattern for the background
	pattern := ""
	for y := 0; y < m.height-1; y++ {
		row := ""
		for x := 0; x < m.width; x++ {
			if (x+y)%4 == 0 {
				row += "·"
			} else {
				row += " "
			}
		}
		pattern += row
		if y < m.height-2 {
			pattern += "\n"
		}
	}
	backgroundContent := bgStyle.Render(pattern)

	// Add background as Layer with Z=0 (bottom layer)
	canvas.AddLayer(layers.NewLayer(backgroundContent, 0))

	// =========================================================================
	// LAYER 2 (Z=1): Main content panel
	// =========================================================================
	mainPanelWidth := min(60, m.width-10)
	mainPanelHeight := min(15, m.height-8)

	mainPanelStyle := lipgloss.NewStyle().
		Width(mainPanelWidth).
		Height(mainPanelHeight).
		Background(lipgloss.Color("#16213e")).
		Foreground(lipgloss.Color("#e0e0e0")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#0f3460")).
		Padding(1, 2).
		Align(lipgloss.Center, lipgloss.Center)

	mainContent := mainPanelStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e94560")).Render("[Layer Demo]"),
			"",
			"This is the main content panel (Z=1)",
			"It sits on top of the dotted background",
			"",
			lipgloss.NewStyle().Faint(true).Render("Press SPACE to toggle modal"),
			lipgloss.NewStyle().Faint(true).Render("Press Q to quit"),
		),
	)

	// Center the main panel
	mainPanelX := (m.width - lipgloss.Width(mainContent)) / 2
	mainPanelY := (m.height - lipgloss.Height(mainContent)) / 2

	canvas.AddLayer(layers.NewLayer(mainContent, 1).WithOffset(mainPanelX, mainPanelY))

	// =========================================================================
	// LAYER 3 (Z=2): Modal overlay (conditionally shown)
	// =========================================================================
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
				"This modal is on Z=2",
				"It overlays everything below!",
				"",
				lipgloss.NewStyle().Italic(true).Render("Press SPACE to close"),
			),
		)

		// Center the modal
		modalX := (m.width - lipgloss.Width(modalContent)) / 2
		modalY := (m.height - lipgloss.Height(modalContent)) / 2

		canvas.AddLayer(layers.NewLayer(modalContent, 2).WithOffset(modalX, modalY))
	}

	// =========================================================================
	// LAYER 4 (Z=10): Status bar (always on top)
	// =========================================================================
	statusStyle := lipgloss.NewStyle().
		Width(m.width).
		Background(lipgloss.Color("#0f3460")).
		Foreground(lipgloss.Color("#e0e0e0")).
		Padding(0, 1)

	modalStatus := "Modal: OFF"
	if m.showModal {
		modalStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Render("Modal: ON")
	}

	statusContent := statusStyle.Render(
		fmt.Sprintf("Layers Demo │ Size: %dx%d │ %s │ Press SPACE toggle modal, Q quit",
			m.width, m.height, modalStatus),
	)

	// Status bar at the bottom (Y = height - 1), highest Z to always show on top
	canvas.AddLayer(layers.NewLayer(statusContent, 10).WithOffset(0, m.height-2))

	return canvas.RenderWithLipgloss()
}

func main() {
	p := tea.NewProgram(model{})
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
