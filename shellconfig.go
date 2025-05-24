package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ShellConfig struct {
	FilePath        string
	Aliases         map[string]string
	Exports         map[string]string
	OhMyZshTheme    string
	OhMyZshPlugins  []string
	CustomFunctions []string
	RawSections     map[string][]string
}

func NewShellConfig() *ShellConfig {
	homeDir, _ := os.UserHomeDir()
	return &ShellConfig{
		FilePath:        filepath.Join(homeDir, ".zshrc"),
		Aliases:         make(map[string]string),
		Exports:         make(map[string]string),
		OhMyZshPlugins:  []string{},
		CustomFunctions: []string{},
		RawSections:     make(map[string][]string),
	}
}

func (sc *ShellConfig) Load() error {
	// Clear existing data before loading
	sc.Aliases = make(map[string]string)
	sc.Exports = make(map[string]string)
	sc.OhMyZshTheme = ""
	sc.OhMyZshPlugins = []string{}
	sc.CustomFunctions = []string{}
	sc.RawSections = make(map[string][]string)

	file, err := os.Open(sc.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentSection := "other"
	inFunction := false
	functionLines := []string{}

	aliasRegex := regexp.MustCompile(`^\s*alias\s+(\w+)=['"](.+)['"]`)
	exportRegex := regexp.MustCompile(`^\s*export\s+(\w+)=['"]?(.+?)['"]?\s*$`)
	themeRegex := regexp.MustCompile(`^\s*ZSH_THEME=['"](.+)['"]`)
	pluginsRegex := regexp.MustCompile(`^\s*plugins=\((.*)\)`)
	functionStartRegex := regexp.MustCompile(`^\s*(\w+)\s*\(\)\s*{`)
	functionEndRegex := regexp.MustCompile(`^\s*}`)

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			if !inFunction {
				sc.RawSections[currentSection] = append(sc.RawSections[currentSection], line)
			}
			continue
		}

		if inFunction {
			functionLines = append(functionLines, line)
			if functionEndRegex.MatchString(line) {
				inFunction = false
				sc.CustomFunctions = append(sc.CustomFunctions, strings.Join(functionLines, "\n"))
				functionLines = []string{}
			}
			continue
		}

		if matches := functionStartRegex.FindStringSubmatch(line); matches != nil {
			inFunction = true
			functionLines = []string{line}
			continue
		}

		if matches := aliasRegex.FindStringSubmatch(line); matches != nil {
			sc.Aliases[matches[1]] = matches[2]
			currentSection = "aliases"
		} else if matches := exportRegex.FindStringSubmatch(line); matches != nil {
			sc.Exports[matches[1]] = matches[2]
			currentSection = "exports"
		} else if matches := themeRegex.FindStringSubmatch(line); matches != nil {
			sc.OhMyZshTheme = matches[1]
			currentSection = "ohmyzsh"
		} else if matches := pluginsRegex.FindStringSubmatch(line); matches != nil {
			pluginsStr := matches[1]
			plugins := strings.Fields(pluginsStr)
			sc.OhMyZshPlugins = plugins
			currentSection = "ohmyzsh"
		} else {
			sc.RawSections[currentSection] = append(sc.RawSections[currentSection], line)
		}
	}

	return scanner.Err()
}

func (sc *ShellConfig) Save() error {
	file, err := os.Create(sc.FilePath + ".tmp")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if sc.OhMyZshTheme != "" || len(sc.OhMyZshPlugins) > 0 {
		writer.WriteString("# Oh My Zsh Configuration\n")
		if sc.OhMyZshTheme != "" {
			writer.WriteString(fmt.Sprintf("ZSH_THEME=\"%s\"\n", sc.OhMyZshTheme))
		}
		if len(sc.OhMyZshPlugins) > 0 {
			writer.WriteString(fmt.Sprintf("plugins=(%s)\n", strings.Join(sc.OhMyZshPlugins, " ")))
		}
		writer.WriteString("\n")
	}

	if len(sc.Exports) > 0 {
		writer.WriteString("# Environment Variables\n")
		for key, value := range sc.Exports {
			if strings.Contains(value, " ") || strings.Contains(value, "$") {
				writer.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, value))
			} else {
				writer.WriteString(fmt.Sprintf("export %s=%s\n", key, value))
			}
		}
		writer.WriteString("\n")
	}

	if len(sc.Aliases) > 0 {
		writer.WriteString("# Aliases\n")
		for name, command := range sc.Aliases {
			writer.WriteString(fmt.Sprintf("alias %s='%s'\n", name, command))
		}
		writer.WriteString("\n")
	}

	if len(sc.CustomFunctions) > 0 {
		writer.WriteString("# Custom Functions\n")
		for _, function := range sc.CustomFunctions {
			writer.WriteString(function)
			writer.WriteString("\n\n")
		}
	}

	for section, lines := range sc.RawSections {
		if section != "other" && len(lines) > 0 {
			writer.WriteString(fmt.Sprintf("# %s\n", section))
		}
		for _, line := range lines {
			writer.WriteString(line)
			writer.WriteString("\n")
		}
		if len(lines) > 0 {
			writer.WriteString("\n")
		}
	}

	writer.Flush()
	file.Close()

	return os.Rename(sc.FilePath+".tmp", sc.FilePath)
}

func (sc *ShellConfig) GetAvailableThemes() ([]string, error) {
	homeDir, _ := os.UserHomeDir()
	themesDir := filepath.Join(homeDir, ".oh-my-zsh", "themes")

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		return []string{}, err
	}

	themes := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".zsh-theme") {
			themeName := strings.TrimSuffix(entry.Name(), ".zsh-theme")
			themes = append(themes, themeName)
		}
	}
	return themes, nil
}

func (sc *ShellConfig) GetAvailablePlugins() ([]string, error) {
	homeDir, _ := os.UserHomeDir()
	pluginsDir := filepath.Join(homeDir, ".oh-my-zsh", "plugins")

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return []string{}, err
	}

	plugins := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}
	return plugins, nil
}
