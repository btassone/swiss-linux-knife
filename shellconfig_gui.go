package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ShellConfigGUI struct {
	config        *ShellConfig
	window        fyne.Window
	aliasesTable  *widget.Table
	exportsTable  *widget.Table
	themeSelect   *widget.Select
	pluginsList   *widget.List
	pluginChecks  map[string]bool
}

func createShellConfigContent(window fyne.Window) fyne.CanvasObject {
	config := NewShellConfig()
	gui := &ShellConfigGUI{
		config:       config,
		pluginChecks: make(map[string]bool),
		window:       window,
	}

	if err := config.Load(); err != nil {
		return container.NewCenter(
			widget.NewLabel(fmt.Sprintf("Error loading config: %v", err)),
		)
	}

	for _, plugin := range config.OhMyZshPlugins {
		gui.pluginChecks[plugin] = true
	}

	tabs := container.NewAppTabs(
		container.NewTabItem("Environment", gui.createEnvironmentTab()),
		container.NewTabItem("Path", gui.createPathTab()),
		container.NewTabItem("Aliases", gui.createAliasesTab()),
		container.NewTabItem("Oh My Zsh", gui.createOhMyZshTab()),
		container.NewTabItem("Functions", gui.createFunctionsTab()),
		container.NewTabItem("History", gui.createHistoryTab()),
	)

	saveButton := widget.NewButton("Save Configuration", func() {
		gui.saveConfiguration()
	})
	saveButton.Importance = widget.HighImportance

	reloadButton := widget.NewButton("Reload", func() {
		config.Load()
		gui.refreshAll()
	})

	return container.NewBorder(
		nil,
		container.NewPadded(
			container.NewHBox(
				saveButton,
				reloadButton,
			),
		),
		nil,
		nil,
		tabs,
	)
}

func (gui *ShellConfigGUI) createEnvironmentTab() fyne.CanvasObject {
	exportData := [][]string{}
	for key, value := range gui.config.Exports {
		exportData = append(exportData, []string{key, value})
	}

	gui.exportsTable = widget.NewTable(
		func() (int, int) { return len(exportData), 2 },
		func() fyne.CanvasObject {
			return widget.NewEntry()
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			entry := cell.(*widget.Entry)
			entry.SetText(exportData[id.Row][id.Col])
			entry.OnChanged = func(text string) {
				exportData[id.Row][id.Col] = text
				gui.updateExportsFromTable(exportData)
			}
		},
	)
	gui.exportsTable.SetColumnWidth(0, 200)
	gui.exportsTable.SetColumnWidth(1, 400)

	addButton := widget.NewButton("Add Variable", func() {
		exportData = append(exportData, []string{"NEW_VAR", ""})
		gui.exportsTable.Refresh()
	})

	removeButton := widget.NewButton("Remove Selected", func() {
		dialog.ShowInformation("Info", "Select a row to remove", gui.window)
	})

	return container.NewBorder(
		widget.NewCard("Environment Variables", "", 
			container.NewHBox(addButton, removeButton),
		),
		nil,
		nil,
		nil,
		container.NewScroll(gui.exportsTable),
	)
}

func (gui *ShellConfigGUI) createPathTab() fyne.CanvasObject {
	currentPath := gui.config.Exports["PATH"]
	if currentPath == "" {
		// Get system PATH if not set
		currentPath = os.Getenv("PATH")
	}
	paths := strings.Split(currentPath, ":")
	
	// Remove empty paths
	filteredPaths := []string{}
	for _, p := range paths {
		if p != "" {
			filteredPaths = append(filteredPaths, p)
		}
	}
	paths = filteredPaths

	pathData := make([][]string, len(paths))
	for i, p := range paths {
		pathData[i] = []string{p}
	}

	var pathTable *widget.Table
	pathTable = widget.NewTable(
		func() (int, int) { return len(pathData), 1 },
		func() fyne.CanvasObject {
			orderLabel := widget.NewLabel("1.")
			orderLabel.TextStyle = fyne.TextStyle{Bold: true}
			entry := widget.NewEntry()
			btn := widget.NewButton("X", func() {})
			btn.Resize(fyne.NewSize(40, 30))
			return container.NewBorder(nil, nil, orderLabel, btn, entry)
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			box := cell.(*fyne.Container)
			// In BorderContainer: Objects[0] is center (entry), Objects[1] is left (orderLabel), Objects[2] is right (btn)
			entry := box.Objects[0].(*widget.Entry)
			orderLabel := box.Objects[1].(*widget.Label)
			btn := box.Objects[2].(*widget.Button)
			
			if id.Row < len(pathData) {
				// Set order number
				orderLabel.SetText(fmt.Sprintf("%d.", id.Row+1))
				
				entry.SetText(pathData[id.Row][0])
				entry.OnChanged = func(text string) {
					if id.Row < len(pathData) {
						pathData[id.Row][0] = text
					}
				}
				
				row := id.Row
				btn.OnTapped = func() {
					// Remove the path
					if row < len(pathData) {
						pathData = append(pathData[:row], pathData[row+1:]...)
						pathTable.Refresh()
					}
				}
			}
		},
	)
	pathTable.SetColumnWidth(0, 800)

	// Add path with file browser
	addWithBrowserBtn := widget.NewButton("Add Directory...", func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err == nil && dir != nil {
				pathData = append(pathData, []string{dir.Path()})
				pathTable.Refresh()
			}
		}, gui.window)
	})

	// Add custom path
	customPathEntry := widget.NewEntry()
	customPathEntry.SetPlaceHolder("/path/to/directory")
	
	addCustomBtn := widget.NewButton("Add", func() {
		if customPathEntry.Text != "" {
			pathData = append(pathData, []string{customPathEntry.Text})
			pathTable.Refresh()
			customPathEntry.SetText("")
		}
	})

	// Common paths to add
	commonPaths := []struct{ name, path string }{
		{"Home bin", "$HOME/bin"},
		{"Local bin", "/usr/local/bin"},
		{"Snap bin", "/snap/bin"},
		{"Go bin", "$HOME/go/bin"},
		{"Cargo bin", "$HOME/.cargo/bin"},
		{"NPM bin", "$HOME/.npm/bin"},
		{"Python bin", "$HOME/.local/bin"},
	}

	quickAddSelect := widget.NewSelect([]string{}, func(selected string) {
		for _, cp := range commonPaths {
			if cp.name == selected {
				expandedPath := os.ExpandEnv(cp.path)
				pathData = append(pathData, []string{expandedPath})
				pathTable.Refresh()
				break
			}
		}
	})
	
	options := []string{}
	for _, cp := range commonPaths {
		options = append(options, cp.name)
	}
	quickAddSelect.Options = options
	quickAddSelect.PlaceHolder = "Quick add common paths..."

	// Track selected row
	var selectedRow int = -1
	
	// Add selection handling to the table
	pathTable.OnSelected = func(id widget.TableCellID) {
		selectedRow = id.Row
	}
	
	// Move up/down buttons
	moveUpBtn := widget.NewButton("Move Up", func() {
		if selectedRow > 0 && selectedRow < len(pathData) {
			// Swap with previous item
			pathData[selectedRow-1], pathData[selectedRow] = pathData[selectedRow], pathData[selectedRow-1]
			selectedRow--
			pathTable.Refresh()
			pathTable.Select(widget.TableCellID{Row: selectedRow, Col: 0})
		}
	})
	
	moveDownBtn := widget.NewButton("Move Down", func() {
		if selectedRow >= 0 && selectedRow < len(pathData)-1 {
			// Swap with next item
			pathData[selectedRow], pathData[selectedRow+1] = pathData[selectedRow+1], pathData[selectedRow]
			selectedRow++
			pathTable.Refresh()
			pathTable.Select(widget.TableCellID{Row: selectedRow, Col: 0})
		}
	})

	// Save the PATH
	savePathBtn := widget.NewButton("Apply PATH Changes", func() {
		newPaths := []string{}
		for _, row := range pathData {
			if row[0] != "" {
				newPaths = append(newPaths, row[0])
			}
		}
		gui.config.Exports["PATH"] = strings.Join(newPaths, ":")
		if gui.exportsTable != nil {
			gui.exportsTable.Refresh()
		}
		dialog.ShowInformation("Success", "PATH updated. Save configuration to make permanent.", gui.window)
	})
	savePathBtn.Importance = widget.HighImportance

	// Info card
	infoCard := widget.NewCard("About PATH", "", widget.NewLabel(
		"The PATH variable tells the system where to find executable programs.\n"+
		"Directories are searched in order from top to bottom.\n"+
		"Changes will be applied to your shell configuration when saved.",
	))

	topControls := container.NewVBox(
		widget.NewCard("Add Path Entry", "", container.NewVBox(
			addWithBrowserBtn,
			widget.NewSeparator(),
			container.NewBorder(nil, nil, nil, addCustomBtn, customPathEntry),
			quickAddSelect,
		)),
		widget.NewCard("Reorder", "", container.NewHBox(moveUpBtn, moveDownBtn)),
	)

	return container.NewBorder(
		container.NewHBox(topControls, infoCard),
		container.NewPadded(savePathBtn),
		nil,
		nil,
		container.NewScroll(pathTable),
	)
}

func (gui *ShellConfigGUI) createAliasesTab() fyne.CanvasObject {
	aliasData := [][]string{}
	for name, command := range gui.config.Aliases {
		aliasData = append(aliasData, []string{name, command})
	}

	gui.aliasesTable = widget.NewTable(
		func() (int, int) { return len(aliasData), 2 },
		func() fyne.CanvasObject {
			return widget.NewEntry()
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			entry := cell.(*widget.Entry)
			entry.SetText(aliasData[id.Row][id.Col])
			entry.OnChanged = func(text string) {
				aliasData[id.Row][id.Col] = text
				gui.updateAliasesFromTable(aliasData)
			}
		},
	)
	gui.aliasesTable.SetColumnWidth(0, 150)
	gui.aliasesTable.SetColumnWidth(1, 450)

	addButton := widget.NewButton("Add Alias", func() {
		aliasData = append(aliasData, []string{"newalias", "command"})
		gui.aliasesTable.Refresh()
	})

	commonAliases := []struct{ name, cmd string }{
		{"ll", "ls -alF"},
		{"la", "ls -A"},
		{"l", "ls -CF"},
		{"..", "cd .."},
		{"...", "cd ../.."},
		{"gs", "git status"},
		{"gd", "git diff"},
		{"gc", "git commit"},
		{"gp", "git push"},
	}

	quickAdd := widget.NewSelect([]string{}, func(selected string) {
		for _, alias := range commonAliases {
			if fmt.Sprintf("%s: %s", alias.name, alias.cmd) == selected {
				aliasData = append(aliasData, []string{alias.name, alias.cmd})
				gui.aliasesTable.Refresh()
				gui.updateAliasesFromTable(aliasData)
				break
			}
		}
	})

	options := []string{}
	for _, alias := range commonAliases {
		options = append(options, fmt.Sprintf("%s: %s", alias.name, alias.cmd))
	}
	quickAdd.Options = options
	quickAdd.PlaceHolder = "Quick add common aliases..."

	return container.NewBorder(
		widget.NewCard("Shell Aliases", "",
			container.NewVBox(
				container.NewHBox(addButton, quickAdd),
			),
		),
		nil,
		nil,
		nil,
		container.NewScroll(gui.aliasesTable),
	)
}

func (gui *ShellConfigGUI) createOhMyZshTab() fyne.CanvasObject {
	themes, _ := gui.config.GetAvailableThemes()
	gui.themeSelect = widget.NewSelect(themes, func(selected string) {
		gui.config.OhMyZshTheme = selected
	})
	gui.themeSelect.SetSelected(gui.config.OhMyZshTheme)

	plugins, _ := gui.config.GetAvailablePlugins()
	
	gui.pluginsList = widget.NewList(
		func() int { return len(plugins) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", func(bool) {}),
				widget.NewLabel("Plugin Name"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			pluginName := plugins[id]
			box := item.(*fyne.Container)
			check := box.Objects[0].(*widget.Check)
			label := box.Objects[1].(*widget.Label)
			
			label.SetText(pluginName)
			check.SetChecked(gui.pluginChecks[pluginName])
			check.OnChanged = func(checked bool) {
				gui.pluginChecks[pluginName] = checked
				gui.updatePluginsFromChecks()
			}
		},
	)

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search plugins...")
	searchEntry.OnChanged = func(text string) {
		// TODO: Implement plugin search/filter
	}

	popularPlugins := widget.NewCard("Popular Plugins", "", container.NewVBox(
		widget.NewButton("Enable git", func() {
			gui.pluginChecks["git"] = true
			gui.updatePluginsFromChecks()
			gui.pluginsList.Refresh()
		}),
		widget.NewButton("Enable docker", func() {
			gui.pluginChecks["docker"] = true
			gui.updatePluginsFromChecks()
			gui.pluginsList.Refresh()
		}),
		widget.NewButton("Enable kubectl", func() {
			gui.pluginChecks["kubectl"] = true
			gui.updatePluginsFromChecks()
			gui.pluginsList.Refresh()
		}),
	))

	themeCard := widget.NewCard("Theme Selection", "", container.NewVBox(
		gui.themeSelect,
		widget.NewButton("Preview Theme", func() {
			dialog.ShowInformation("Theme Preview", 
				"Theme preview coming soon!\nCurrent theme: "+gui.config.OhMyZshTheme, 
				gui.window)
		}),
	))

	return container.NewBorder(
		container.NewVBox(
			themeCard,
			widget.NewCard("Plugin Selection", "", searchEntry),
		),
		popularPlugins,
		nil,
		nil,
		container.NewScroll(gui.pluginsList),
	)
}

func (gui *ShellConfigGUI) createFunctionsTab() fyne.CanvasObject {
	functionsList := widget.NewList(
		func() int { return len(gui.config.CustomFunctions) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Function")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			label := item.(*widget.Label)
			funcText := gui.config.CustomFunctions[id]
			lines := strings.Split(funcText, "\n")
			if len(lines) > 0 {
				label.SetText(lines[0])
			}
		},
	)

	functionEditor := widget.NewMultiLineEntry()
	functionEditor.SetPlaceHolder("Select a function to edit...")
	functionEditor.Resize(fyne.NewSize(600, 400))

	functionsList.OnSelected = func(id widget.ListItemID) {
		if id < len(gui.config.CustomFunctions) {
			functionEditor.SetText(gui.config.CustomFunctions[id])
		}
	}

	addButton := widget.NewButton("Add Function", func() {
		newFunc := "newfunction() {\n    # Add your code here\n    echo \"Hello from new function\"\n}"
		gui.config.CustomFunctions = append(gui.config.CustomFunctions, newFunc)
		functionsList.Refresh()
	})

	saveButton := widget.NewButton("Save Changes", func() {
		if functionsList.Length() > 0 {
			// TODO: Update the selected function with editor content
		}
	})

	templateSelect := widget.NewSelect(
		[]string{"Directory Navigation", "Git Helper", "Docker Shortcut"},
		func(selected string) {
			var template string
			switch selected {
			case "Directory Navigation":
				template = "mkcd() {\n    mkdir -p \"$1\" && cd \"$1\"\n}"
			case "Git Helper":
				template = "gcommit() {\n    git add . && git commit -m \"$1\"\n}"
			case "Docker Shortcut":
				template = "dexec() {\n    docker exec -it \"$1\" /bin/bash\n}"
			}
			functionEditor.SetText(template)
		},
	)
	templateSelect.PlaceHolder = "Function templates..."

	leftPanel := container.NewBorder(
		container.NewVBox(
			addButton,
			templateSelect,
		),
		nil,
		nil,
		nil,
		functionsList,
	)

	rightPanel := container.NewBorder(
		nil,
		saveButton,
		nil,
		nil,
		container.NewScroll(functionEditor),
	)

	return container.NewHSplit(leftPanel, rightPanel)
}

func (gui *ShellConfigGUI) createHistoryTab() fyne.CanvasObject {
	historyData := []string{}
	currentFilter := ""
	
	historyList := widget.NewList(
		func() int { 
			if currentFilter == "" {
				return len(historyData)
			}
			count := 0
			for _, entry := range historyData {
				if strings.Contains(strings.ToLower(entry), strings.ToLower(currentFilter)) {
					count++
				}
			}
			return count
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("History entry")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			label := item.(*widget.Label)
			if currentFilter == "" {
				if id < len(historyData) {
					label.SetText(historyData[id])
				}
			} else {
				// Show filtered results
				count := 0
				for _, entry := range historyData {
					if strings.Contains(strings.ToLower(entry), strings.ToLower(currentFilter)) {
						if count == id {
							label.SetText(entry)
							break
						}
						count++
					}
				}
			}
		},
	)
	
	loadHistory := func() []string {
		// Try to read history using zsh with proper initialization
		cmd := exec.Command("zsh", "-i", "-c", "history -1000")
		output, err := cmd.Output()
		
		if err != nil || len(output) == 0 {
			// Try reading the history file directly
			homeDir, _ := os.UserHomeDir()
			histFile := filepath.Join(homeDir, ".zsh_history")
			
			// Check if .zsh_history exists
			if _, err := os.Stat(histFile); err != nil {
				// Try .bash_history as fallback
				histFile = filepath.Join(homeDir, ".bash_history")
			}
			
			data, err := os.ReadFile(histFile)
			if err != nil {
				return []string{"Error reading history: " + err.Error()}
			}
			
			lines := strings.Split(string(data), "\n")
			result := make([]string, 0, len(lines))
			
			// Process lines in reverse order (newest first)
			for i := len(lines) - 1; i >= 0; i-- {
				line := strings.TrimSpace(lines[i])
				// Skip metadata lines in zsh history (they start with :)
				if line != "" && !strings.HasPrefix(line, ":") {
					result = append(result, line)
				}
			}
			
			return result
		}
		
		// Parse the history output (format: "number  command")
		lines := strings.Split(string(output), "\n")
		result := make([]string, 0, len(lines))
		
		// Process lines in reverse order (newest first)
		for i := len(lines) - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}
			
			// Remove the line number prefix (e.g., "  123  command" -> "command")
			parts := strings.Fields(line)
			if len(parts) > 1 {
				// Skip the first field (line number) and join the rest
				command := strings.Join(parts[1:], " ")
				result = append(result, command)
			}
		}
		
		return result
	}
	
	// Initial load
	historyData = loadHistory()
	if len(historyData) == 0 {
		historyData = []string{"No history found. Click Refresh to try again."}
	}
	historyList.Refresh()
	
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search history...")
	searchEntry.OnChanged = func(text string) {
		currentFilter = text
		historyList.Refresh()
	}
	
	refreshButton := widget.NewButton("Refresh", func() {
		historyData = loadHistory()
		historyList.Refresh()
	})
	
	clearButton := widget.NewButton("Clear History", func() {
		dialog.ShowConfirm("Clear History", "Are you sure you want to clear your shell history?", func(clear bool) {
			if clear {
				exec.Command("zsh", "-c", "history -c").Run()
				historyData = loadHistory()
				historyList.Refresh()
			}
		}, gui.window)
	})
	
	copyButton := widget.NewButton("Copy Selected", func() {
		// TODO: Implement copy to clipboard when item is selected
		dialog.ShowInformation("Info", "Select a history item to copy", gui.window)
	})
	
	topBar := container.NewBorder(nil, nil, nil,
		container.NewHBox(refreshButton, clearButton, copyButton),
		searchEntry,
	)
	
	// Clean up timer when tab is closed
	historyTab := container.NewBorder(
		topBar,
		widget.NewLabel("Recent commands (click Refresh to update)"),
		nil,
		nil,
		container.NewScroll(historyList),
	)
	
	return historyTab
}

func (gui *ShellConfigGUI) updateAliasesFromTable(data [][]string) {
	gui.config.Aliases = make(map[string]string)
	for _, row := range data {
		if len(row) == 2 && row[0] != "" {
			gui.config.Aliases[row[0]] = row[1]
		}
	}
}

func (gui *ShellConfigGUI) updateExportsFromTable(data [][]string) {
	gui.config.Exports = make(map[string]string)
	for _, row := range data {
		if len(row) == 2 && row[0] != "" {
			gui.config.Exports[row[0]] = row[1]
		}
	}
}

func (gui *ShellConfigGUI) updatePluginsFromChecks() {
	gui.config.OhMyZshPlugins = []string{}
	for plugin, checked := range gui.pluginChecks {
		if checked {
			gui.config.OhMyZshPlugins = append(gui.config.OhMyZshPlugins, plugin)
		}
	}
}


func (gui *ShellConfigGUI) saveConfiguration() {
	if err := gui.config.Save(); err != nil {
		dialog.ShowError(err, gui.window)
		return
	}
	dialog.ShowInformation("Success", "Configuration saved successfully!", gui.window)
}

func (gui *ShellConfigGUI) refreshAll() {
	// Reset plugin checks based on loaded data
	gui.pluginChecks = make(map[string]bool)
	for _, plugin := range gui.config.OhMyZshPlugins {
		gui.pluginChecks[plugin] = true
	}
	
	// Refresh all UI components
	if gui.aliasesTable != nil {
		gui.aliasesTable.Refresh()
	}
	if gui.exportsTable != nil {
		gui.exportsTable.Refresh()
	}
	if gui.themeSelect != nil {
		gui.themeSelect.SetSelected(gui.config.OhMyZshTheme)
	}
	if gui.pluginsList != nil {
		gui.pluginsList.Refresh()
	}
}