# Layers Package

A Z-axis layer compositing system for terminal UIs built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Overview

The layers package provides a way to composite multiple UI elements on top of each other with proper Z-ordering and ANSI escape sequence handling. This is essential for rendering modals, popups, notifications, and other overlays in terminal applications.

## Installation

```go
import "github.com/mightymoud/arlocode/internal/tui/layers"
```

## Core Concepts

### Layer

A `Layer` represents a single piece of content with:
- **Content**: The rendered string (can include ANSI styling)
- **Z**: Z-axis position (higher values render on top)
- **X, Y**: Position offset from top-left corner
- **Visible**: Toggle visibility without removing the layer

### Canvas

A `Canvas` is the rendering surface that composites multiple layers together.

## Basic Usage

### Simple Overlay

```go
canvas := layers.NewCanvas(80, 24) // width, height

// Add background layer (Z=0)
background := lipgloss.NewStyle().
    Background(lipgloss.Color("#1e1e2e")).
    Width(80).Height(24).
    Render("Main content here...")
canvas.AddLayer(layers.NewLayer(background, 0))

// Add modal on top (Z=1)
modal := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Padding(1, 2).
    Render("Modal content")
canvas.AddLayer(layers.CenterLayer(modal, 1, 80, 24))

// Render the composite
output := canvas.Render()
```

### Positioned Layers

```go
canvas := layers.NewCanvas(80, 24)

// Add a layer at a specific position
box := "┌───┐\n│ ! │\n└───┘"
canvas.AddLayer(layers.NewLayer(box, 1).WithOffset(10, 5))
```

### Visibility Toggle

```go
layer := layers.NewLayer(content, 1).WithVisibility(false)
canvas.AddLayer(layer)

// Layer won't render until visibility is set to true
```

## API Reference

### Layer Functions

```go
// Create a new layer
NewLayer(content string, z int) Layer

// Position the layer
layer.WithOffset(x, y int) Layer

// Toggle visibility
layer.WithVisibility(visible bool) Layer

// Create a centered layer
CenterLayer(content string, z int, canvasWidth, canvasHeight int) Layer
```

### Canvas Functions

```go
// Create a new canvas
NewCanvas(width, height int) *Canvas

// Add a layer
canvas.AddLayer(layer Layer) *Canvas

// Clear all layers
canvas.ClearLayers() *Canvas

// Render the composite
canvas.Render() string
canvas.RenderWithLipgloss() string  // Alias for Render()
```

### Utility Functions

```go
// Quick overlay of two strings
OverlaySimple(background, foreground string, width, height int) string
```

## Compose API

For a more fluent interface:

```go
result := layers.NewCompose(80, 24).
    WithStyle(baseStyle).
    Layer(background, 0, 0, 0).           // content, z, x, y
    CenteredLayer(modal, 1).              // content, z (auto-centered)
    Render()
```

## Integration with Bubble Tea

```go
func (m model) View() string {
    canvas := layers.NewCanvas(m.width, m.height)
    
    // Base UI layer
    mainUI := m.renderMainContent()
    canvas.AddLayer(layers.NewLayer(mainUI, 0))
    
    // Modal overlay (conditional)
    if m.showModal {
        modal := m.renderModal()
        canvas.AddLayer(layers.CenterLayer(modal, 1, m.width, m.height))
    }
    
    // Notifications on top
    if m.notifications.HasActiveNotifications() {
        notifs := m.notifications.Render()
        canvas.AddLayer(layers.NewLayer(notifs, 2).WithOffset(m.width-42, 2))
    }
    
    return canvas.Render()
}
```

## Z-Order Priority

Layers are rendered from lowest Z to highest Z:
- Z=0: Background/base content
- Z=1: Modals, dialogs
- Z=2: Notifications, tooltips
- Z=3+: Highest priority overlays

Higher Z values always appear on top of lower Z values at overlapping positions.

## Notes

- The package handles ANSI escape sequences properly when compositing layers
- Spaces in higher layers will overwrite content in lower layers
- For transparent overlays, ensure your content only covers the area you want to change
