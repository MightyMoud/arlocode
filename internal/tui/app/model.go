package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mightymoud/arlocode/internal/tui/notifications"
)

type AppModel struct {
	width         int
	height        int
	showModal     bool
	MainInput     textinput.Model
	ModalInput    textinput.Model
	Notifications *notifications.NotificationManager
}
