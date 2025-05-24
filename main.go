package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Tool struct {
	Name        string
	Description string
	Icon        fyne.Resource
	Content     func(window fyne.Window) fyne.CanvasObject
}

var tools = []Tool{
	{
		Name:        "Shell Config Manager",
		Description: "Manage bashrc/zshrc visually",
		Icon:        theme.SettingsIcon(),
		Content:     createShellConfigContent,
	},
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DefaultTheme())

	myWindow := myApp.NewWindow("Swiss Linux Knife")
	myWindow.Resize(fyne.NewSize(800, 600))

	content := container.NewStack()

	list := widget.NewList(
		func() int { return len(tools) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Tool Name"),
				widget.NewLabel("Description"),
			)
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			c := obj.(*fyne.Container)
			nameLabel := c.Objects[0].(*widget.Label)
			descLabel := c.Objects[1].(*widget.Label)

			nameLabel.SetText(tools[i].Name)
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			descLabel.SetText(tools[i].Description)
			descLabel.TextStyle = fyne.TextStyle{Italic: true}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		content.Objects = []fyne.CanvasObject{tools[id].Content(myWindow)}
		content.Refresh()
	}

	split := container.NewHSplit(
		container.NewBorder(
			widget.NewLabelWithStyle("Tools", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			nil, nil, nil,
			list,
		),
		container.NewBorder(
			widget.NewLabelWithStyle("Swiss Linux Knife", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			nil, nil, nil,
			container.NewPadded(content),
		),
	)
	split.SetOffset(0.3)

	myWindow.SetContent(split)

	list.Select(0)

	myWindow.ShowAndRun()
}

