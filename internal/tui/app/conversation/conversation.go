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
		Conversation: []ConversationMessage{
			{
				Type:    "user",
				Content: "Hi there!",
			},
			{
				Type:    "thinking",
				Content: "Let me think about that...",
			},
			{
				Type:    "agent",
				Content: "Hello! How can I assist you today?",
			},
		},
		AgentThinking:  false,
		ThinkingBuffer: "",
	}
}

func (cm *ConversationManager) AddThinkingMessage(content string) {
	conversationTurn := ConversationMessage{
		Type:    "thinking",
		Content: cm.ThinkingBuffer,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
	cm.AgentThinking = false
	cm.ThinkingBuffer = ""
}

func (cm *ConversationManager) AddAgentMessage(content string) {
	conversationTurn := ConversationMessage{
		Type:    "agent",
		Content: cm.TextBuffer,
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
