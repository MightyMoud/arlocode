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
			{Type: "user", Content: "You are a helpful coding assistant1"},
			{Type: "user", Content: "You are a helpful coding assistant2."},
			{Type: "user", Content: "You are a helpful coding assistant3."},
			{Type: "user", Content: "You are a helpful coding assistant4."},
			{Type: "user", Content: "You are a helpful coding assistant5."},
			{Type: "user", Content: "You are a helpful coding assistant6."},
			{Type: "agent", Content: "You are a helpful coding assistant7."},
			{Type: "agent", Content: "You are a helpful coding assistant8."},
			{Type: "agent", Content: "You are a helpful coding assistant9."},
			{Type: "agent", Content: "You are a helpful coding assistant10."},
			{Type: "agent", Content: "You are a helpful coding assistant11."},
			{Type: "agent", Content: "You are a helpful coding assistant12."},
			{Type: "agent", Content: "You are a helpful coding assistant13."},
			{Type: "agent", Content: "You are a helpful coding assistant14."},
			{Type: "agent", Content: "You are a helpful coding assistant15."},
		},
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
