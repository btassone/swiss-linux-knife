package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/btassone/swiss-linux-knife/internal/tools"
)

func createMainMenu(app fyne.App, window fyne.Window) *fyne.MainMenu {
	// File menu items with shortcuts
	newItem := fyne.NewMenuItem("New", func() {
		dialog.ShowInformation("New", "Create new configuration", window)
	})
	newItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}

	openItem := fyne.NewMenuItem("Open...", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				defer reader.Close()
				dialog.ShowInformation("Open", "Opened: "+reader.URI().Path(), window)
			}
		}, window)
	})
	openItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}

	saveItem := fyne.NewMenuItem("Save", func() {
		dialog.ShowInformation("Save", "Configuration saved", window)
	})
	saveItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}

	saveAsItem := fyne.NewMenuItem("Save As...", func() {
		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err == nil && writer != nil {
				defer writer.Close()
				dialog.ShowInformation("Save As", "Saved to: "+writer.URI().Path(), window)
			}
		}, window)
	})
	saveAsItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}

	exitItem := fyne.NewMenuItem("Exit", func() {
		app.Quit()
	})
	exitItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyQ, Modifier: fyne.KeyModifierControl}

	fileMenu := fyne.NewMenu("File",
		newItem,
		openItem,
		saveItem,
		saveAsItem,
		fyne.NewMenuItemSeparator(),
		exitItem,
	)

	// Edit menu items with shortcuts
	cutItem := fyne.NewMenuItem("Cut", func() {
		dialog.ShowInformation("Cut", "Cut functionality not yet implemented", window)
	})
	cutItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyX, Modifier: fyne.KeyModifierControl}

	copyItem := fyne.NewMenuItem("Copy", func() {
		dialog.ShowInformation("Copy", "Copy functionality not yet implemented", window)
	})
	copyItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyC, Modifier: fyne.KeyModifierControl}

	pasteItem := fyne.NewMenuItem("Paste", func() {
		dialog.ShowInformation("Paste", "Paste functionality not yet implemented", window)
	})
	pasteItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyV, Modifier: fyne.KeyModifierControl}

	selectAllItem := fyne.NewMenuItem("Select All", func() {
		dialog.ShowInformation("Select All", "Select All functionality not yet implemented", window)
	})
	selectAllItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyA, Modifier: fyne.KeyModifierControl}

	editMenu := fyne.NewMenu("Edit",
		cutItem,
		copyItem,
		pasteItem,
		fyne.NewMenuItemSeparator(),
		selectAllItem,
	)

	// Tools menu
	toolsMenu := fyne.NewMenu("Tools",
		fyne.NewMenuItem("Preferences...", func() {
			showPreferences(window)
		}),
		fyne.NewMenuItem("Options...", func() {
			dialog.ShowInformation("Options", "Options dialog not yet implemented", window)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Reload Configuration", func() {
			dialog.ShowInformation("Reload", "Configuration reloaded", window)
		}),
	)

	// Help menu
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			dialog.ShowInformation("Documentation", "Visit: https://github.com/btassone/swiss-linux-knife", window)
		}),
		fyne.NewMenuItem("Report Issue", func() {
			dialog.ShowInformation("Report Issue", "Visit: https://github.com/btassone/swiss-linux-knife/issues", window)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("About", func() {
			showAbout(window)
		}),
	)

	return fyne.NewMainMenu(fileMenu, editMenu, toolsMenu, helpMenu)
}

func showPreferences(window fyne.Window) {
	// Create preferences dialog
	themeSelect := widget.NewSelect([]string{"Light", "Dark"}, func(selected string) {
		// Theme selection logic
	})
	themeSelect.SetSelected("Light")

	content := container.NewVBox(
		widget.NewLabel("Application Preferences"),
		widget.NewSeparator(),
		container.NewHBox(widget.NewLabel("Theme:"), themeSelect),
	)

	dialog.ShowCustom("Preferences", "Close", content, window)
}

func showAbout(window fyne.Window) {
	content := container.NewVBox(
		widget.NewLabelWithStyle("Swiss Linux Knife", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Version 1.0.0"),
		widget.NewSeparator(),
		widget.NewLabel("A collection of GUI tools for Linux system management"),
		widget.NewLabel(""),
		widget.NewLabel("© 2024 Swiss Linux Knife Contributors"),
		widget.NewLabel("Licensed under MIT License"),
	)

	dialog.ShowCustom("About Swiss Linux Knife", "OK", content, window)
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DefaultTheme())

	myWindow := myApp.NewWindow("Swiss Linux Knife")
	myWindow.Resize(fyne.NewSize(1200, 700))

	// Set up the menu
	mainMenu := createMainMenu(myApp, myWindow)
	myWindow.SetMainMenu(mainMenu)

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