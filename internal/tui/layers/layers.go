package layers

import (
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Layer represents a single layer with content and Z-axis position.
// Higher Z values are rendered on top of lower Z values.
type Layer struct {
	Content string // The rendered content of this layer (can include ANSI codes)
	Z       int    // Z-axis value (higher = on top)
	X       int    // X offset from top-left corner
	Y       int    // Y offset from top-left corner
	Visible bool   // Whether this layer should be rendered
}

// NewLayer creates a new layer with the given content and Z position.
func NewLayer(content string, z int) Layer {
	return Layer{
		Content: content,
		Z:       z,
		X:       0,
		Y:       0,
		Visible: true,
	}
}

// WithOffset sets the X and Y offset for positioning the layer.
func (l Layer) WithOffset(x, y int) Layer {
	l.X = x
	l.Y = y
	return l
}

// WithVisibility sets whether the layer is visible.
func (l Layer) WithVisibility(visible bool) Layer {
	l.Visible = visible
	return l
}

// Canvas represents a rendering surface that can composite multiple layers.
type Canvas struct {
	Width  int
	Height int
	Layers []Layer
}

// NewCanvas creates a new canvas with the given dimensions.
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		Width:  width,
		Height: height,
		Layers: make([]Layer, 0),
	}
}

// AddLayer adds a layer to the canvas.
func (c *Canvas) AddLayer(layer Layer) *Canvas {
	c.Layers = append(c.Layers, layer)
	return c
}

// ClearLayers removes all layers from the canvas.
func (c *Canvas) ClearLayers() *Canvas {
	c.Layers = make([]Layer, 0)
	return c
}

// RenderWithLipgloss composites all layers.
// Layers are sorted by Z-axis (lowest first), and higher Z layers
// are placed on top of lower Z layers with proper ANSI handling.
func (c *Canvas) RenderWithLipgloss() string {
	if len(c.Layers) == 0 {
		return ""
	}

	// Sort layers by Z (lowest first, so higher Z renders on top)
	sortedLayers := make([]Layer, len(c.Layers))
	copy(sortedLayers, c.Layers)
	sort.Slice(sortedLayers, func(i, j int) bool {
		return sortedLayers[i].Z < sortedLayers[j].Z
	})

	// Filter visible layers
	visibleLayers := make([]Layer, 0, len(sortedLayers))
	for _, layer := range sortedLayers {
		if layer.Visible {
			visibleLayers = append(visibleLayers, layer)
		}
	}

	if len(visibleLayers) == 0 {
		return ""
	}

	// Build the output by placing each layer
	// Start with an empty canvas
	output := createEmptyCanvas(c.Width, c.Height)

	for _, layer := range visibleLayers {
		output = placeContentAt(output, layer.Content, layer.X, layer.Y, c.Width, c.Height)
	}

	return output
}

// Render is an alias for RenderWithLipgloss for backward compatibility.
func (c *Canvas) Render() string {
	return c.RenderWithLipgloss()
}

// createEmptyCanvas creates an empty canvas of the given dimensions.
func createEmptyCanvas(width, height int) string {
	row := strings.Repeat(" ", width)
	rows := make([]string, height)
	for i := range rows {
		rows[i] = row
	}
	return strings.Join(rows, "\n")
}

// placeContentAt places content at a specific x, y position on the background.
// It handles ANSI sequences properly by working with full styled lines.
func placeContentAt(background, content string, x, y, width, height int) string {
	bgLines := strings.Split(background, "\n")
	contentLines := strings.Split(content, "\n")

	// Ensure bgLines has correct height
	for len(bgLines) < height {
		bgLines = append(bgLines, strings.Repeat(" ", width))
	}

	// Place content lines onto background
	for i, contentLine := range contentLines {
		targetY := y + i
		if targetY < 0 || targetY >= height {
			continue
		}

		contentWidth := lipgloss.Width(contentLine)
		if contentWidth == 0 {
			continue
		}

		// Get the background line
		bgLine := bgLines[targetY]
		bgWidth := lipgloss.Width(bgLine)

		// Build the new line with content placed at x position
		newLine := placeStringAt(bgLine, contentLine, x, bgWidth, width)
		bgLines[targetY] = newLine
	}

	return strings.Join(bgLines, "\n")
}

// placeStringAt places a styled string at position x within a background line.
// This handles ANSI codes properly.
func placeStringAt(bgLine, content string, x, bgWidth, maxWidth int) string {
	contentWidth := lipgloss.Width(content)

	// If content starts beyond the line width, just return background
	if x >= maxWidth {
		return bgLine
	}

	// Build the result
	var result strings.Builder

	// Part 1: Background before content (0 to x)
	if x > 0 {
		leftPart := truncateWithAnsi(bgLine, 0, x)
		result.WriteString(leftPart)
	}

	// Part 2: The content itself
	result.WriteString(content)

	// Part 3: Background after content (x + contentWidth to end)
	afterX := x + contentWidth
	if afterX < bgWidth {
		rightPart := truncateWithAnsi(bgLine, afterX, bgWidth)
		result.WriteString(rightPart)
	}

	return result.String()
}

// truncateWithAnsi extracts a visual portion of a styled string.
// This is a simplified version that works for basic cases.
func truncateWithAnsi(s string, start, end int) string {
	if start >= end {
		return ""
	}

	// Use lipgloss width to properly handle styled content
	width := lipgloss.Width(s)
	if start >= width {
		return ""
	}

	// For simple cases without complex ANSI, we can use substring
	// For strings with ANSI codes, we need to be smarter
	if !strings.Contains(s, "\x1b") {
		// No ANSI codes, simple substring
		runes := []rune(s)
		if end > len(runes) {
			end = len(runes)
		}
		if start > len(runes) {
			return ""
		}
		return string(runes[start:end])
	}

	// Has ANSI codes - extract visual characters while preserving codes
	return extractVisualRange(s, start, end)
}

// extractVisualRange extracts characters from visual position start to end,
// preserving ANSI escape sequences.
func extractVisualRange(s string, start, end int) string {
	var result strings.Builder
	var currentEscape strings.Builder
	inEscape := false
	visualPos := 0
	activeEscapes := ""

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			currentEscape.Reset()
			currentEscape.WriteRune(r)
			continue
		}

		if inEscape {
			currentEscape.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
				escapeSeq := currentEscape.String()
				// Track active escape sequences for colors
				if visualPos < start {
					activeEscapes += escapeSeq
				} else if visualPos < end {
					result.WriteString(escapeSeq)
				}
			}
			continue
		}

		// Regular character
		if visualPos >= start && visualPos < end {
			// First character in range - prepend any active escape sequences
			if visualPos == start && activeEscapes != "" {
				result.WriteString(activeEscapes)
			}
			result.WriteRune(r)
		}
		visualPos++

		if visualPos >= end {
			break
		}
	}

	return result.String()
}

// CenterLayer creates a layer centered within the given dimensions.
func CenterLayer(content string, z, canvasWidth, canvasHeight int) Layer {
	contentWidth := lipgloss.Width(content)
	contentHeight := lipgloss.Height(content)

	x := (canvasWidth - contentWidth) / 2
	y := (canvasHeight - contentHeight) / 2

	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return Layer{
		Content: content,
		Z:       z,
		X:       x,
		Y:       y,
		Visible: true,
	}
}

// OverlaySimple is a simple function to overlay two strings.
// The foreground string is placed on top of the background.
func OverlaySimple(background, foreground string, width, height int) string {
	canvas := NewCanvas(width, height)
	canvas.AddLayer(NewLayer(background, 0))
	canvas.AddLayer(NewLayer(foreground, 1))
	return canvas.RenderWithLipgloss()
}

// Compose creates a styled layer composition using lipgloss styles.
type Compose struct {
	style  lipgloss.Style
	layers []Layer
	width  int
	height int
}

// NewCompose creates a new composition with base style and dimensions.
func NewCompose(width, height int) *Compose {
	return &Compose{
		style:  lipgloss.NewStyle(),
		layers: make([]Layer, 0),
		width:  width,
		height: height,
	}
}

// WithStyle sets the base style for the composition.
func (c *Compose) WithStyle(style lipgloss.Style) *Compose {
	c.style = style
	return c
}

// Layer adds a layer with content, Z position, and optional positioning.
func (c *Compose) Layer(content string, z int, x, y int) *Compose {
	c.layers = append(c.layers, Layer{
		Content: content,
		Z:       z,
		X:       x,
		Y:       y,
		Visible: true,
	})
	return c
}

// CenteredLayer adds a centered layer.
func (c *Compose) CenteredLayer(content string, z int) *Compose {
	c.layers = append(c.layers, CenterLayer(content, z, c.width, c.height))
	return c
}

// Render renders the composition.
func (c *Compose) Render() string {
	canvas := NewCanvas(c.width, c.height)
	for _, layer := range c.layers {
		canvas.AddLayer(layer)
	}
	return c.style.Render(canvas.RenderWithLipgloss())
}
