package state

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mightymoud/arlocode/internal/butler/agent"
)

var (
	instance *AppState
	once     sync.Once
)

type AppState struct {
	mu      sync.RWMutex
	program *tea.Program
	agent   *agent.Agent
}

func Get() *AppState {
	once.Do(func() {
		instance = &AppState{}
	})
	return instance
}

func (s *AppState) SetProgram(p *tea.Program) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.program = p
}

func (s *AppState) Program() *tea.Program {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.program
}

func (s *AppState) SetAgent(a *agent.Agent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.agent = a
}

func (s *AppState) Agent() *agent.Agent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.agent
}

// Helper to send messages safely
func (s *AppState) Send(msg tea.Msg) {
	if p := s.Program(); p != nil {
		p.Send(msg)
	}
}
