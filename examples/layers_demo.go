// This file demonstrates various examples of using the layers package
// for Z-axis compositing in terminal UIs with bubbletea and lipgloss.
package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mightymoud/arlocode/internal/layers"
)

func main() {
	fmt.Println("=== Layer Rendering Examples ===")

	// Example 1: Basic Overlay
	example1BasicOverlay()

	// Example 2: Z-Order Priority
	example2ZOrderPriority()

	// Example 3: Positioned Layers
	example3PositionedLayers()

	// Example 4: Modal Dialog Over Content
	example4ModalDialog()

	// Example 5: Using Compose API
	example5ComposeAPI()

	// Example 6: Transparency vs Opaque
	example6TransparencyVsOpaque()
}

// Example 1: Basic overlay of foreground on background
func example1BasicOverlay() {
	fmt.Println("--- Example 1: Basic Overlay ---")

	width, height := 20, 5

	// Background: filled with dots
	background := strings.Repeat(strings.Repeat(".", width)+"\n", height-1) +
		strings.Repeat(".", width)

	// Foreground: A simple box
	foreground := "┌────────┐\n│ HELLO! │\n└────────┘"

	result := layers.OverlaySimple(background, foreground, width, height)
	fmt.Println(result)
	fmt.Println()
}

// Example 2: Demonstrating Z-order priority
func example2ZOrderPriority() {
	fmt.Println("--- Example 2: Z-Order Priority ---")
	fmt.Println("Three layers at same position with Z=0, Z=1, Z=2")
	fmt.Println("Highest Z wins:")

	canvas := layers.NewCanvas(15, 3)

	// Layer at Z=0
	canvas.AddLayer(layers.NewLayer("AAAAAAAAAA", 0))

	// Layer at Z=1 (partially overlaps)
	canvas.AddLayer(layers.NewLayer("  BBBBBB", 1))

	// Layer at Z=2 (overlaps B)
	canvas.AddLayer(layers.NewLayer("    CCCC", 2))

	fmt.Println(canvas.Render())
	fmt.Println("Result: A shows where nothing covers it, B shows where C doesn't cover, C on top")
	fmt.Println()
}

// Example 3: Layers at different positions
func example3PositionedLayers() {
	fmt.Println("--- Example 3: Positioned Layers ---")

	canvas := layers.NewCanvas(30, 10)

	// Fill background
	bg := ""
	for i := 0; i < 10; i++ {
		bg += strings.Repeat("░", 30)
		if i < 9 {
			bg += "\n"
		}
	}
	canvas.AddLayer(layers.NewLayer(bg, 0))

	// Top-left box at Z=1
	box1 := "┌───┐\n│ 1 │\n└───┘"
	canvas.AddLayer(layers.NewLayer(box1, 1).WithOffset(2, 1))

	// Overlapping box at Z=2
	box2 := "┌───┐\n│ 2 │\n└───┘"
	canvas.AddLayer(layers.NewLayer(box2, 2).WithOffset(4, 2))

	// Bottom-right box at Z=1
	box3 := "┌───┐\n│ 3 │\n└───┘"
	canvas.AddLayer(layers.NewLayer(box3, 1).WithOffset(22, 6))

	fmt.Println(canvas.Render())
	fmt.Println()
}

// Example 4: Modal dialog overlaying content
func example4ModalDialog() {
	fmt.Println("--- Example 4: Modal Dialog Over Content ---")

	width, height := 40, 12

	canvas := layers.NewCanvas(width, height)

	// Main content (Z=0)
	contentStyle := lipgloss.NewStyle().
		Width(width - 4).
		Border(lipgloss.RoundedBorder())

	mainContent := contentStyle.Render(
		"This is the main application content.\n" +
			"It contains important information.\n" +
			"A modal will appear on top of this.")

	canvas.AddLayer(layers.NewLayer(mainContent, 0).WithOffset(2, 1))

	// Modal dialog (Z=1)
	modalStyle := lipgloss.NewStyle().
		Width(24).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("9")).
		Padding(0, 1)

	modal := modalStyle.Render(
		"⚠ Confirm Action?\n" +
			"[Yes]  [No]")

	// Center the modal
	modalLayer := layers.CenterLayer(modal, 1, width, height)
	canvas.AddLayer(modalLayer)

	fmt.Println(canvas.Render())
	fmt.Println()
}

// Example 5: Using the Compose API for fluent layer building
func example5ComposeAPI() {
	fmt.Println("--- Example 5: Using Compose API ---")

	result := layers.NewCompose(25, 7).
		Layer(strings.Repeat(strings.Repeat("▒", 25)+"\n", 6)+strings.Repeat("▒", 25), 0, 0, 0).
		CenteredLayer("═══════════\n  CENTER  \n═══════════", 1).
		Render()

	fmt.Println(result)
	fmt.Println()
}

// Example 6: Difference between Render (transparent) and RenderOpaque
func example6TransparencyVsOpaque() {
	fmt.Println("--- Example 6: Transparency vs Opaque ---")

	width, height := 20, 5

	// Background
	bg := ""
	for i := 0; i < height; i++ {
		bg += strings.Repeat("X", width)
		if i < height-1 {
			bg += "\n"
		}
	}

	// Foreground with spaces (which act as transparent or opaque)
	fg := "  HELLO  "

	// Transparent render - spaces show background through
	canvas1 := layers.NewCanvas(width, height)
	canvas1.AddLayer(layers.NewLayer(bg, 0))
	canvas1.AddLayer(layers.NewLayer(fg, 1).WithOffset(5, 2))

	fmt.Println("Transparent render (spaces show background):")
	fmt.Println(canvas1.Render())

	// Opaque render - spaces overwrite background
	canvas2 := layers.NewCanvas(width, height)
	canvas2.AddLayer(layers.NewLayer(bg, 0))
	canvas2.AddLayer(layers.NewLayer(fg, 1).WithOffset(5, 2))

	fmt.Println("\nOpaque render (spaces cover background):")
	fmt.Println()
}
