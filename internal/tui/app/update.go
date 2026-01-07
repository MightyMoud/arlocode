package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Notifications.UpdateScreenSize(msg.Width, msg.Height)

	case tickMsg:
		// Update notification animations
		if m.Notifications.Update() {
			cmds = append(cmds, tickCmd())
		}
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.showModal {
				m.showModal = false
				m.ModalInput.Blur()
				m.MainInput.Focus()
				return m, textinput.Blink
			}
			return m, tea.Quit
		case "enter":
			if m.showModal {
				// Close modal on enter
				m.showModal = false
				m.ModalInput.Blur()
				m.MainInput.Focus()
				return m, textinput.Blink
			}
		case "ctrl+o":
			// Toggle modal
			m.showModal = !m.showModal
			if m.showModal {
				m.MainInput.Blur()
				m.ModalInput.Focus()
				return m, textinput.Blink
			} else {
				m.ModalInput.Blur()
				m.MainInput.Focus()
				return m, textinput.Blink
			}
		case "w":
			// Show warning notification
			if !m.showModal {
				m.Notifications.PushWarning("Warning", "This is a warning notification!")
				return m, tickCmd()
			}
		}

		// Route input to the focused component
		if m.showModal {
			m.ModalInput, cmd = m.ModalInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.MainInput, cmd = m.MainInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}
