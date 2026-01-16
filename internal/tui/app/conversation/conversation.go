package conversation

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mightymoud/arlocode/internal/butler/tools"
)

type MessageStatus string

const (
	StatusInProgress MessageStatus = "inProgress"
	StatusComplete   MessageStatus = "complete"
)

type ConversationMessage struct {
	Type      string // agent, user, thinking, tool_call
	Content   string
	StartTime time.Time
	EndTime   time.Time
	Status    MessageStatus
}

func (cm *ConversationMessage) IsType(t string) bool {
	return cm.Type == t
}

type ConversationManager struct {
	Conversation  []ConversationMessage
	AgentThinking bool
}

func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		Conversation:  []ConversationMessage{},
		AgentThinking: false,
	}
}

func (cm *ConversationManager) StartThinkingMessage() {
	conversationTurn := ConversationMessage{
		Type:      "thinking",
		Content:   "",
		StartTime: time.Now(),
		Status:    StatusInProgress,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
	cm.AgentThinking = true
}

func (cm *ConversationManager) UpdateThinkingMessage(content string) {
	thinkingMessage := cm.GetLastMessage()
	thinkingMessage.Content += content
}

func (cm *ConversationManager) CompleteAgentThinkingMessage() {
	thinkingMessage := cm.GetLastMessage()
	thinkingMessage.EndTime = time.Now()
	thinkingMessage.Status = StatusComplete
	cm.AgentThinking = false
}

func (cm *ConversationManager) StartAgentMessage() {
	conversationTurn := ConversationMessage{
		Type:      "agent",
		Content:   "",
		StartTime: time.Now(),
		Status:    StatusInProgress,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
}

func (cm *ConversationManager) UpdateAgentMessage(content string) {
	agentMessage := cm.GetLastMessage()
	agentMessage.Content += content
}

func (cm *ConversationManager) CompleteAgentMessage() {
	agentMessage := cm.GetLastMessage()
	agentMessage.EndTime = time.Now()
	agentMessage.Status = StatusComplete
}

func (cm *ConversationManager) AddUserMessage(content string) {
	conversationTurn := ConversationMessage{
		Type:      "user",
		Content:   content,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Status:    StatusComplete,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
}

func (cm *ConversationManager) AddToolCallMessage(tc tools.ToolCall) {
	// Convert to underlying ToolCall and create a conversation entry

	// Create a readable representation of arguments
	argsJSON, err := json.Marshal(tc.Arguments)
	argsStr := ""
	if err == nil {
		argsStr = string(argsJSON)
	} else {
		argsStr = fmt.Sprintf("(failed to marshal args: %v)", err)
	}

	content := fmt.Sprintf("%s %s", tc.FunctionName, argsStr)

	conversationTurn := ConversationMessage{
		Type:      "tool_call",
		Content:   content,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Status:    StatusComplete,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
}

func (cm *ConversationManager) IsEmpty() bool {
	return len(cm.Conversation) == 0
}

func (cm *ConversationManager) GetLastMessage() *ConversationMessage {
	if len(cm.Conversation) == 0 {
		return nil
	}
	return &cm.Conversation[len(cm.Conversation)-1]
}
