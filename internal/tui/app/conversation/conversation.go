package conversation

type ConversationMessage struct {
	Type    string
	Content string
}

type ConversationManager struct {
	Conversation   []ConversationMessage
	AgentThinking  bool
	ThinkingBuffer string
}

func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		Conversation: []ConversationMessage{
			{
				Type:    "user",
				Content: "Why are you here",
			},
			{
				Type:    "thinking",
				Content: "I think therefore I am.",
			},
			{
				Type:    "agent",
				Content: "this is my first message!",
			},
		},
		AgentThinking:  false,
		ThinkingBuffer: "",
	}
}

func (cm *ConversationManager) AddThinkingMessage(content string) {
	conversationTurn := ConversationMessage{
		Type:    "agent_thinking",
		Content: cm.ThinkingBuffer,
	}
	cm.Conversation = append(cm.Conversation, conversationTurn)
	cm.AgentThinking = false
	cm.ThinkingBuffer = ""
}

func (cm *ConversationManager) IsEmpty() bool {
	return len(cm.Conversation) == 0
}
