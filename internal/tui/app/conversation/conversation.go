package conversation

type ConversationMessage struct {
	Type    string
	Content string
}

type ConversationManager struct {
	Conversation   []ConversationMessage
	AgentThinking  bool
	ThinkingBuffer string
	TextBuffer     string
}

func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		Conversation:   []ConversationMessage{},
		AgentThinking:  false,
		ThinkingBuffer: "",
	}
}

func (cm *ConversationManager) AddThinkingMessage(content string) {
	conversationTurn := ConversationMessage{
		Type:    "thinking",
		Content: content,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
	cm.AgentThinking = false
	cm.ThinkingBuffer = ""
}

func (cm *ConversationManager) AddAgentMessage(content string) {
	conversationTurn := ConversationMessage{
		Type:    "agent",
		Content: content,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
	cm.TextBuffer = ""
}

func (cm *ConversationManager) AddUserMessage(content string) {
	conversationTurn := ConversationMessage{
		Type:    "user",
		Content: content,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
}

func (cm *ConversationManager) IsEmpty() bool {
	return len(cm.Conversation) == 0
}
