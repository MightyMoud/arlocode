package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m AppModel) Init() tea.Cmd {
	return textinput.Blink
}
