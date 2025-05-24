package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/btassone/swiss-linux-knife/internal/tools"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DefaultTheme())

	myWindow := myApp.NewWindow("Swiss Linux Knife")
	myWindow.Resize(fyne.NewSize(1200, 700))

	content := container.NewStack()
	availableTools := tools.GetAvailableTools()

	list := widget.NewList(
		func() int { return len(availableTools) },
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

			nameLabel.SetText(availableTools[i].Name)
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			descLabel.SetText(availableTools[i].Description)
			descLabel.TextStyle = fyne.TextStyle{Italic: true}
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		content.Objects = []fyne.CanvasObject{availableTools[id].Content(myWindow)}
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
	split.SetOffset(0.22)

	myWindow.SetContent(split)

	list.Select(0)

	myWindow.ShowAndRun()
}