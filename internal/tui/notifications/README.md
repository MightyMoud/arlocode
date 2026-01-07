# Notifications Package

An animated notification system for terminal UIs using [Harmonica](https://github.com/charmbracelet/harmonica) spring physics for smooth slide-in animations.

## Overview

The notifications package provides toast-style notifications that animate in from the right side of the screen, display for a configurable duration, then disappear. It integrates seamlessly with Bubble Tea applications and the layers package for proper Z-ordering.

## Installation

```go
import "github.com/mightymoud/arlocode/internal/tui/notifications"
```

## Quick Start

```go
// Create a notification manager
notifs := notifications.NewNotificationManager(screenWidth, screenHeight)

// Push notifications
notifs.PushInfo("Info", "This is an informational message")
notifs.PushSuccess("Success", "Operation completed!")
notifs.PushWarning("Warning", "Please review your input")
notifs.PushError("Error", "Something went wrong")
```

## Notification Types

| Type | Method | Default Color |
|------|--------|---------------|
| Info | `PushInfo(title, message)` | Blue |
| Success | `PushSuccess(title, message)` | Green |
| Warning | `PushWarning(title, message)` | Yellow |
| Error | `PushError(title, message)` | Red |

## Configuration

```go
notifs := notifications.NewNotificationManager(width, height).
    WithMaxVisible(5).                          // Max notifications shown at once
    WithDefaultWidth(40).                       // Width in characters
    WithDefaultDuration(4 * time.Second).       // How long to display
    WithSpringConfig(12.0, 0.6)                 // Animation frequency & damping
```

## Integration with Bubble Tea

### Model Setup

```go
type model struct {
    width         int
    height        int
    notifications *notifications.NotificationManager
}

func initialModel() model {
    return model{
        notifications: notifications.NewNotificationManager(80, 24),
    }
}
```

### Update Function

The notification manager needs to be updated on each tick to animate:

```go
type tickMsg time.Time

func tickCmd() tea.Cmd {
    return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.notifications.UpdateScreenSize(msg.Width, msg.Height)

    case tickMsg:
        // Update notification animations
        if m.notifications.Update() {
            return m, tickCmd() // Continue ticking if animations active
        }

    case tea.KeyMsg:
        switch msg.String() {
        case "i":
            m.notifications.PushInfo("Info", "Hello world!")
            return m, tickCmd() // Start animation tick
        }
    }
    return m, nil
}
```

### View Function with Layers

```go
func (m model) View() string {
    canvas := layers.NewCanvas(m.width, m.height)
    
    // Base content (Z=0)
    mainContent := m.renderMain()
    canvas.AddLayer(layers.NewLayer(mainContent, 0))
    
    // Notifications (Z=2, above modals)
    if m.notifications.HasActiveNotifications() {
        notifContent := m.notifications.Render()
        // Position at top-right, offset from edge
        x := m.width - m.notifications.GetWidth() - 2
        y := 2
        canvas.AddLayer(layers.NewLayer(notifContent, 2).WithOffset(x, y))
    }
    
    return canvas.Render()
}
```

## API Reference

### NotificationManager

```go
// Create manager
NewNotificationManager(screenWidth, screenHeight int) *NotificationManager

// Configuration (chainable)
WithMaxVisible(max int) *NotificationManager
WithDefaultWidth(width int) *NotificationManager
WithDefaultDuration(d time.Duration) *NotificationManager
WithSpringConfig(frequency, damping float64) *NotificationManager

// Update screen size (call on WindowSizeMsg)
UpdateScreenSize(width, height int)

// Push notifications (returns notification ID)
Push(title, message string, notifType NotificationType) string
PushInfo(title, message string) string
PushSuccess(title, message string) string
PushWarning(title, message string) string
PushError(title, message string) string

// Dismiss notifications
Dismiss(id string)      // Dismiss specific notification
DismissAll()            // Dismiss all notifications

// Animation update (call each frame, returns true if still animating)
Update() bool

// State queries
HasActiveNotifications() bool
Count() int

// Rendering
Render() string                                    // Standalone render
RenderWithPosition(width, height int) string       // Positioned for layers integration
```

## Animation Lifecycle

Each notification goes through these states:

1. **SlideIn**: Animates from off-screen right to visible position using spring physics
2. **Visible**: Displays for the configured duration (default 4 seconds)
3. **Done**: Removed from the notification stack

## Spring Animation

The package uses Harmonica spring physics for natural-feeling animations:

- **Frequency** (default 12.0): Higher = faster animation
- **Damping** (default 0.6): Higher = less bounce/overshoot

```go
// Snappy, minimal bounce
notifs.WithSpringConfig(15.0, 0.8)

// Bouncy, playful
notifs.WithSpringConfig(8.0, 0.3)
```

## Example: Full Integration

```go
package main

import (
    "time"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/mightymoud/arlocode/internal/tui/layers"
    "github.com/mightymoud/arlocode/internal/tui/notifications"
)

type model struct {
    width, height int
    notifs        *notifications.NotificationManager
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
    return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.notifs.UpdateScreenSize(m.width, m.height)
        
    case tickMsg:
        if m.notifs.Update() {
            return m, tickCmd()
        }
        
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "1":
            m.notifs.PushInfo("Info", "Press keys 1-4 for notifications")
            return m, tickCmd()
        case "2":
            m.notifs.PushSuccess("Success", "It worked!")
            return m, tickCmd()
        case "3":
            m.notifs.PushWarning("Warning", "Be careful...")
            return m, tickCmd()
        case "4":
            m.notifs.PushError("Error", "Something broke!")
            return m, tickCmd()
        }
    }
    return m, nil
}

func (m model) View() string {
    canvas := layers.NewCanvas(m.width, m.height)
    
    // Main content
    canvas.AddLayer(layers.NewLayer("Press 1-4 for notifications", 0))
    
    // Notifications overlay
    if m.notifs.HasActiveNotifications() {
        canvas.AddLayer(
            layers.NewLayer(m.notifs.Render(), 2).
                WithOffset(m.width-42, 2),
        )
    }
    
    return canvas.Render()
}

func main() {
    m := model{
        notifs: notifications.NewNotificationManager(80, 24),
    }
    tea.NewProgram(m, tea.WithAltScreen()).Run()
}
```
