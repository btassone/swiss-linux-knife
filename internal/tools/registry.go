package tools

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/btassone/swiss-linux-knife/internal/gui"
)

type Tool struct {
	Name        string
	Description string
	Icon        fyne.Resource
	Content     func(window fyne.Window) fyne.CanvasObject
}

func GetAvailableTools() []Tool {
	return []Tool{
		{
			Name:        "Shell Config Manager",
			Description: "Manage bashrc/zshrc visually",
			Icon:        theme.SettingsIcon(),
			Content: func(window fyne.Window) fyne.CanvasObject {
				shellConfigGUI := gui.NewShellConfigGUI(window)
				return shellConfigGUI.CreateContent()
			},
		},
	}
}