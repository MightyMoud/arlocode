// Package notifications provides an animated notification system for terminal UIs.
// Notifications slide in from the top-right corner, display for a configurable duration,
// then slide out. Uses charmbracelet/harmonica for smooth spring-based animations.
package notifications

import (
	"strings"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// NotificationType defines the visual style of a notification
type NotificationType int

const (
	NotificationInfo NotificationType = iota
	NotificationSuccess
	NotificationWarning
	NotificationError
)

// AnimationState tracks the current phase of a notification's lifecycle
type AnimationState int

const (
	AnimationSlideIn AnimationState = iota
	AnimationVisible
	AnimationSlideOut
	AnimationDone
)

// Notification represents a single notification with its content and animation state
type Notification struct {
	ID        string
	Title     string
	Message   string
	Type      NotificationType
	State     AnimationState
	CreatedAt time.Time

	// Animation spring for smooth movement
	spring   harmonica.Spring
	position float64 // Current X position (offset from right edge)
	velocity float64

	// Timing
	displayDuration time.Duration
	visibleSince    time.Time

	// Dimensions
	width int
}

// NotificationManager handles multiple notifications and their animations
type NotificationManager struct {
	notifications []*Notification
	maxVisible    int
	screenWidth   int
	screenHeight  int

	// Default styling
	defaultWidth    int
	defaultDuration time.Duration

	// Spring configuration
	springFrequency float64
	springDamping   float64

	// Notification ID counter
	idCounter int
}

// NewNotificationManager creates a new manager for handling notifications
func NewNotificationManager(screenWidth, screenHeight int) *NotificationManager {
	return &NotificationManager{
		notifications:   make([]*Notification, 0),
		maxVisible:      5,
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		defaultWidth:    40,
		defaultDuration: 4 * time.Second,
		springFrequency: 6.0,
		springDamping:   0.8,
		idCounter:       0,
	}
}

// WithMaxVisible sets the maximum number of visible notifications
func (m *NotificationManager) WithMaxVisible(max int) *NotificationManager {
	m.maxVisible = max
	return m
}

// WithDefaultWidth sets the default notification width
func (m *NotificationManager) WithDefaultWidth(width int) *NotificationManager {
	m.defaultWidth = width
	return m
}

// WithDefaultDuration sets the default display duration
func (m *NotificationManager) WithDefaultDuration(d time.Duration) *NotificationManager {
	m.defaultDuration = d
	return m
}

// WithSpringConfig sets the spring animation parameters
func (m *NotificationManager) WithSpringConfig(frequency, damping float64) *NotificationManager {
	m.springFrequency = frequency
	m.springDamping = damping
	return m
}

// UpdateScreenSize updates the screen dimensions
func (m *NotificationManager) UpdateScreenSize(width, height int) {
	m.screenWidth = width
	m.screenHeight = height
}

// Push adds a new notification to the queue
func (m *NotificationManager) Push(title, message string, notifType NotificationType) string {
	m.idCounter++
	id := time.Now().Format("20060102150405") + string(rune('A'+m.idCounter%26))

	// Create spring for animation - starts off-screen to the right (positive = off right edge)
	spring := harmonica.NewSpring(harmonica.FPS(60), m.springFrequency, m.springDamping)

	notif := &Notification{
		ID:              id,
		Title:           title,
		Message:         message,
		Type:            notifType,
		State:           AnimationSlideIn,
		CreatedAt:       time.Now(),
		spring:          spring,
		position:        float64(m.defaultWidth + 10), // Start off-screen to the right (positive value)
		velocity:        0,
		displayDuration: m.defaultDuration,
		width:           m.defaultWidth,
	}

	m.notifications = append(m.notifications, notif)
	return id
}

// PushInfo adds an info notification
func (m *NotificationManager) PushInfo(title, message string) string {
	return m.Push(title, message, NotificationInfo)
}

// PushSuccess adds a success notification
func (m *NotificationManager) PushSuccess(title, message string) string {
	return m.Push(title, message, NotificationSuccess)
}

// PushWarning adds a warning notification
func (m *NotificationManager) PushWarning(title, message string) string {
	return m.Push(title, message, NotificationWarning)
}

// PushError adds an error notification
func (m *NotificationManager) PushError(title, message string) string {
	return m.Push(title, message, NotificationError)
}

// Dismiss removes a notification by ID immediately (no animation)
func (m *NotificationManager) Dismiss(id string) {
	for _, n := range m.notifications {
		if n.ID == id && n.State != AnimationDone {
			n.State = AnimationDone
		}
	}
}

// DismissAll removes all notifications immediately (no animation)
func (m *NotificationManager) DismissAll() {
	for _, n := range m.notifications {
		if n.State != AnimationDone {
			n.State = AnimationDone
		}
	}
}

// Update advances all notification animations. Call this on each frame/tick.
// Returns true if any notifications are still animating (need more updates).
func (m *NotificationManager) Update() bool {
	hasActive := false
	toRemove := make([]int, 0)

	for i, n := range m.notifications {
		switch n.State {
		case AnimationSlideIn:
			// Animate towards visible position (0 = flush with right edge)
			targetPos := float64(0)
			n.position, n.velocity = n.spring.Update(n.position, n.velocity, targetPos)

			// Check if animation is complete (close enough to target)
			if abs(n.position-targetPos) < 0.5 && abs(n.velocity) < 0.5 {
				n.position = targetPos
				n.velocity = 0
				n.State = AnimationVisible
				n.visibleSince = time.Now()
			}
			hasActive = true

		case AnimationVisible:
			// Check if display duration has elapsed - then just remove (no exit animation)
			if time.Since(n.visibleSince) >= n.displayDuration {
				n.State = AnimationDone
				toRemove = append(toRemove, i)
			}
			hasActive = true

		case AnimationSlideOut:
			// No exit animation - immediately mark as done
			n.State = AnimationDone
			toRemove = append(toRemove, i)

		case AnimationDone:
			toRemove = append(toRemove, i)
		}
	}

	// Remove completed notifications (iterate in reverse to preserve indices)
	for i := len(toRemove) - 1; i >= 0; i-- {
		idx := toRemove[i]
		m.notifications = append(m.notifications[:idx], m.notifications[idx+1:]...)
	}

	return hasActive
}

// HasActiveNotifications returns true if there are any notifications being displayed
func (m *NotificationManager) HasActiveNotifications() bool {
	return len(m.notifications) > 0
}

// Count returns the number of active notifications
func (m *NotificationManager) Count() int {
	return len(m.notifications)
}

// Render returns the rendered notification stack positioned for the top-right corner.
// The returned string should be overlaid on your main content using the layers package.
func (m *NotificationManager) Render() string {
	if len(m.notifications) == 0 {
		return ""
	}

	var rendered []string
	visibleCount := 0

	for _, n := range m.notifications {
		if n.State == AnimationDone {
			continue
		}
		if visibleCount >= m.maxVisible {
			break
		}

		rendered = append(rendered, m.renderNotification(n))
		visibleCount++
	}

	if len(rendered) == 0 {
		return ""
	}

	// Stack notifications vertically
	return lipgloss.JoinVertical(lipgloss.Right, rendered...)
}

// RenderWithPosition returns the rendered notifications along with their X,Y position
// for use with the layers package positioning system
func (m *NotificationManager) RenderWithPosition() (content string, x, y int) {
	content = m.Render()
	if content == "" {
		return "", 0, 0
	}

	// Calculate X position based on the first notification's animation
	// position is how far off-screen (positive = right), so we add it to push right
	xOffset := 0
	if len(m.notifications) > 0 {
		xOffset = int(m.notifications[0].position)
	}

	// Get actual rendered width (includes border and padding)
	actualWidth := m.GetNotificationWidth()

	// Position at top-right with margin from edge
	x = m.screenWidth - actualWidth - 1 + xOffset
	y = 4 // Position below header (header is 3 rows)

	return content, x, y
}

// renderNotification renders a single notification with appropriate styling
func (m *NotificationManager) renderNotification(n *Notification) string {
	// Content width (excluding border and padding)
	contentWidth := n.width

	// Base style - fixed width box
	baseStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())

	// Apply type-specific colors
	switch n.Type {
	case NotificationInfo:
		baseStyle = baseStyle.
			BorderForeground(lipgloss.Color("39")). // Blue
			Foreground(lipgloss.Color("252"))
	case NotificationSuccess:
		baseStyle = baseStyle.
			BorderForeground(lipgloss.Color("42")). // Green
			Foreground(lipgloss.Color("252"))
	case NotificationWarning:
		baseStyle = baseStyle.
			BorderForeground(lipgloss.Color("214")). // Orange
			Foreground(lipgloss.Color("252"))
	case NotificationError:
		baseStyle = baseStyle.
			BorderForeground(lipgloss.Color("196")). // Red
			Foreground(lipgloss.Color("252"))
	}

	// Title style
	titleStyle := lipgloss.NewStyle().Bold(true)
	switch n.Type {
	case NotificationInfo:
		titleStyle = titleStyle.Foreground(lipgloss.Color("39"))
	case NotificationSuccess:
		titleStyle = titleStyle.Foreground(lipgloss.Color("42"))
	case NotificationWarning:
		titleStyle = titleStyle.Foreground(lipgloss.Color("214"))
	case NotificationError:
		titleStyle = titleStyle.Foreground(lipgloss.Color("196"))
	}

	// Build content
	var content strings.Builder
	if n.Title != "" {
		content.WriteString(titleStyle.Render(n.Title))
		if n.Message != "" {
			content.WriteString("\n")
		}
	}
	if n.Message != "" {
		// Word wrap the message to fit the content width
		wrapped := wordWrap(n.Message, contentWidth)
		content.WriteString(wrapped)
	}

	return baseStyle.Render(content.String())
}

// GetNotificationWidth returns the actual rendered width of notifications
func (m *NotificationManager) GetNotificationWidth() int {
	// Content width + left padding (2) + right padding (2) + left border (1) + right border (1)
	return m.defaultWidth + 6
}

// wordWrap wraps text to fit within the specified width
func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	lineLen := 0

	for i, word := range words {
		wordLen := len([]rune(word))

		if lineLen+wordLen > width && lineLen > 0 {
			result.WriteString("\n")
			lineLen = 0
		} else if i > 0 && lineLen > 0 {
			result.WriteString(" ")
			lineLen++
		}

		result.WriteString(word)
		lineLen += wordLen
	}

	return result.String()
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
