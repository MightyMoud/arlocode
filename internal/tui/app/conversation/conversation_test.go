package conversation

import (
	"testing"
	"time"
)

func TestConversationMessageWithStatus(t *testing.T) {
	cm := NewConversationManager()

	// Test that we can create and manage thinking messages with status
	cm.StartThinkingMessage()

	if len(cm.Conversation) != 1 {
		t.Errorf("Expected 1 message, got %d", len(cm.Conversation))
	}

	msg := cm.Conversation[0]
	if msg.Type != "thinking" {
		t.Errorf("Expected type 'thinking', got '%s'", msg.Type)
	}

	if msg.Status != StatusInProgress {
		t.Errorf("Expected status '%s', got '%s'", StatusInProgress, msg.Status)
	}

	if msg.StartTime.IsZero() {
		t.Error("Expected StartTime to be set")
	}

	if !msg.EndTime.IsZero() {
		t.Error("Expected EndTime to be zero for in-progress message")
	}

	// Update thinking message
	cm.UpdateThinkingMessage("thinking content")
	msg = cm.Conversation[0]
	if msg.Content != "thinking content" {
		t.Errorf("Expected content 'thinking content', got '%s'", msg.Content)
	}

	// Sleep briefly to ensure duration is measurable
	time.Sleep(10 * time.Millisecond)

	// Complete thinking message
	cm.CompleteAgentThinkingMessage()

	msg = cm.Conversation[0]
	if msg.Status != StatusComplete {
		t.Errorf("Expected status '%s', got '%s'", StatusComplete, msg.Status)
	}

	if msg.EndTime.IsZero() {
		t.Error("Expected EndTime to be set for complete message")
	}

	// Content is preserved by the current implementation
	if msg.Content != "thinking content" {
		t.Errorf("Expected content 'thinking content', got '%s'", msg.Content)
	}

	duration := msg.EndTime.Sub(msg.StartTime)
	if duration < 10*time.Millisecond {
		t.Errorf("Expected duration at least 10ms, got %v", duration)
	}
}

func TestAgentTextMessage(t *testing.T) {
	cm := NewConversationManager()
	cm.StartAgentMessage()
	cm.UpdateAgentMessage("Test response")
	cm.CompleteAgentMessage()

	if len(cm.Conversation) != 1 {
		t.Errorf("Expected 1 message, got %d", len(cm.Conversation))
	}

	msg := cm.Conversation[0]
	if msg.Type != "agent" {
		t.Errorf("Expected type 'agent', got '%s'", msg.Type)
	}

	if msg.Status != StatusComplete {
		t.Errorf("Expected status '%s', got '%s'", StatusComplete, msg.Status)
	}

	if msg.Content != "Test response" {
		t.Errorf("Expected content 'Test response', got '%s'", msg.Content)
	}
}

func TestUserMessage(t *testing.T) {
	cm := NewConversationManager()
	cm.AddUserMessage("Test question")

	if len(cm.Conversation) != 1 {
		t.Errorf("Expected 1 message, got %d", len(cm.Conversation))
	}

	msg := cm.Conversation[0]
	if msg.Type != "user" {
		t.Errorf("Expected type 'user', got '%s'", msg.Type)
	}

	if msg.Status != StatusComplete {
		t.Errorf("Expected status '%s', got '%s'", StatusComplete, msg.Status)
	}

	if msg.Content != "Test question" {
		t.Errorf("Expected content 'Test question', got '%s'", msg.Content)
	}
}

func TestMultipleThinkingMessages(t *testing.T) {
	cm := NewConversationManager()

	// First thinking round
	cm.StartThinkingMessage()
	cm.UpdateThinkingMessage("first thinking")
	cm.CompleteAgentThinkingMessage()

	// Agent response
	cm.StartAgentMessage()
	cm.UpdateAgentMessage("First answer")
	cm.CompleteAgentMessage()

	// Second thinking round
	cm.StartThinkingMessage()
	cm.UpdateThinkingMessage("second thinking")
	cm.CompleteAgentThinkingMessage()

	// Second agent response
	cm.StartAgentMessage()
	cm.UpdateAgentMessage("Second answer")
	cm.CompleteAgentMessage()

	if len(cm.Conversation) != 4 {
		t.Errorf("Expected 4 messages, got %d", len(cm.Conversation))
	}

	// Verify message order
	expectedTypes := []string{"thinking", "agent", "thinking", "agent"}
	for i, expectedType := range expectedTypes {
		if cm.Conversation[i].Type != expectedType {
			t.Errorf("Message %d: expected type '%s', got '%s'", i, expectedType, cm.Conversation[i].Type)
		}
	}
}

func TestThinkingMessageDuration(t *testing.T) {
	cm := NewConversationManager()

	cm.StartThinkingMessage()
	time.Sleep(50 * time.Millisecond)
	cm.CompleteAgentThinkingMessage()

	msg := cm.Conversation[0]
	duration := msg.EndTime.Sub(msg.StartTime)

	if duration < 50*time.Millisecond {
		t.Errorf("Expected duration at least 50ms, got %v", duration)
	}
}
