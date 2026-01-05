package layers

import (
	"strings"
	"testing"
)

func TestNewLayer(t *testing.T) {
	content := "Hello"
	layer := NewLayer(content, 5)

	if layer.Content != content {
		t.Errorf("Expected content %q, got %q", content, layer.Content)
	}
	if layer.Z != 5 {
		t.Errorf("Expected Z=5, got Z=%d", layer.Z)
	}
	if layer.X != 0 || layer.Y != 0 {
		t.Errorf("Expected X=0, Y=0, got X=%d, Y=%d", layer.X, layer.Y)
	}
	if !layer.Visible {
		t.Error("Expected layer to be visible by default")
	}
}

func TestLayerWithOffset(t *testing.T) {
	layer := NewLayer("Test", 1).WithOffset(10, 20)

	if layer.X != 10 || layer.Y != 20 {
		t.Errorf("Expected X=10, Y=20, got X=%d, Y=%d", layer.X, layer.Y)
	}
}

func TestLayerWithVisibility(t *testing.T) {
	layer := NewLayer("Test", 1).WithVisibility(false)

	if layer.Visible {
		t.Error("Expected layer to be invisible")
	}
}

func TestCanvasRender(t *testing.T) {
	canvas := NewCanvas(10, 3)

	// Background layer (Z=0)
	background := strings.Repeat(".", 10) + "\n" +
		strings.Repeat(".", 10) + "\n" +
		strings.Repeat(".", 10)

	// Foreground layer (Z=1) with some content
	foreground := "Hi"

	canvas.AddLayer(NewLayer(background, 0))
	canvas.AddLayer(NewLayer(foreground, 1))

	result := canvas.Render()
	lines := strings.Split(result, "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// First line should have "Hi" at the start
	if !strings.HasPrefix(lines[0], "Hi") {
		t.Errorf("Expected first line to start with 'Hi', got %q", lines[0])
	}
}

func TestCanvasRenderWithOffset(t *testing.T) {
	canvas := NewCanvas(10, 5)

	// Background layer
	bg := strings.Repeat("-", 10) + "\n" +
		strings.Repeat("-", 10) + "\n" +
		strings.Repeat("-", 10) + "\n" +
		strings.Repeat("-", 10) + "\n" +
		strings.Repeat("-", 10)

	// Foreground layer with offset
	fg := "XX\nXX"

	canvas.AddLayer(NewLayer(bg, 0))
	canvas.AddLayer(NewLayer(fg, 1).WithOffset(4, 2))

	result := canvas.Render()
	lines := strings.Split(result, "\n")

	// Check that XX appears at position (4, 2) and (4, 3)
	if len(lines) < 4 {
		t.Fatalf("Expected at least 4 lines, got %d", len(lines))
	}

	expectedLine2 := "----XX----"
	if lines[2] != expectedLine2 {
		t.Errorf("Line 2: expected %q, got %q", expectedLine2, lines[2])
	}

	expectedLine3 := "----XX----"
	if lines[3] != expectedLine3 {
		t.Errorf("Line 3: expected %q, got %q", expectedLine3, lines[3])
	}
}

func TestZOrderRendering(t *testing.T) {
	canvas := NewCanvas(5, 1)

	// Layer at Z=0 with "A"
	canvas.AddLayer(NewLayer("A", 0))
	// Layer at Z=2 with "B" at same position
	canvas.AddLayer(NewLayer("B", 2))
	// Layer at Z=1 with "C" at same position
	canvas.AddLayer(NewLayer("C", 1))

	result := canvas.Render()

	// B should be on top since it has highest Z
	if !strings.HasPrefix(result, "B") {
		t.Errorf("Expected 'B' on top, got %q", result)
	}
}

func TestInvisibleLayer(t *testing.T) {
	canvas := NewCanvas(5, 1)

	canvas.AddLayer(NewLayer("A", 0))
	canvas.AddLayer(NewLayer("B", 1).WithVisibility(false))

	result := canvas.Render()

	// B should not appear since it's invisible
	if strings.Contains(result, "B") {
		t.Errorf("Expected 'B' to be hidden, got %q", result)
	}
	if !strings.HasPrefix(result, "A") {
		t.Errorf("Expected 'A' to show, got %q", result)
	}
}

func TestCenterLayer(t *testing.T) {
	content := "Hi"
	layer := CenterLayer(content, 1, 10, 5)

	// "Hi" is 2 chars wide, so X should be (10-2)/2 = 4
	// 1 line high, so Y should be (5-1)/2 = 2
	if layer.X != 4 {
		t.Errorf("Expected X=4, got X=%d", layer.X)
	}
	if layer.Y != 2 {
		t.Errorf("Expected Y=2, got Y=%d", layer.Y)
	}
}

func TestOverlaySimple(t *testing.T) {
	bg := "-----\n-----\n-----"
	fg := "XX"

	result := OverlaySimple(bg, fg, 5, 3)
	lines := strings.Split(result, "\n")

	if !strings.HasPrefix(lines[0], "XX") {
		t.Errorf("Expected foreground 'XX' on first line, got %q", lines[0])
	}
}

func TestComposeRender(t *testing.T) {
	compose := NewCompose(10, 3).
		Layer(strings.Repeat(".", 10)+"\n"+strings.Repeat(".", 10)+"\n"+strings.Repeat(".", 10), 0, 0, 0).
		CenteredLayer("Hi", 1)

	result := compose.Render()

	if !strings.Contains(result, "Hi") {
		t.Errorf("Expected 'Hi' in result, got %q", result)
	}
}

func TestEmptyCanvas(t *testing.T) {
	canvas := NewCanvas(10, 5)
	result := canvas.Render()

	if result != "" {
		t.Errorf("Expected empty string for empty canvas, got %q", result)
	}
}

func TestClearLayers(t *testing.T) {
	canvas := NewCanvas(10, 5)
	canvas.AddLayer(NewLayer("Test", 0))
	canvas.ClearLayers()

	if len(canvas.Layers) != 0 {
		t.Errorf("Expected 0 layers after clear, got %d", len(canvas.Layers))
	}
}
