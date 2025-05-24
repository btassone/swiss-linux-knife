package shellconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/btassone/swiss-linux-knife/internal/logger"
)

type Config struct {
	FilePath        string
	Aliases         map[string]string
	Exports         map[string]string
	OhMyZshTheme    string
	OhMyZshPlugins  []string
	CustomFunctions []string
	RawSections     map[string][]string
}

func New() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		FilePath:        filepath.Join(homeDir, ".zshrc"),
		Aliases:         make(map[string]string),
		Exports:         make(map[string]string),
		OhMyZshPlugins:  []string{},
		CustomFunctions: []string{},
		RawSections:     make(map[string][]string),
	}
}

func (c *Config) Load() error {
	logger.Debug("Loading shell config from %s", c.FilePath)
	
	c.Aliases = make(map[string]string)
	c.Exports = make(map[string]string)
	c.OhMyZshTheme = ""
	c.OhMyZshPlugins = []string{}
	c.CustomFunctions = []string{}
	c.RawSections = make(map[string][]string)

	file, err := os.Open(c.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("Config file does not exist: %s", c.FilePath)
			return nil
		}
		logger.Error("Failed to open config file: %v", err)
		return fmt.Errorf("failed to open config file: %w", err)
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
				c.RawSections[currentSection] = append(c.RawSections[currentSection], line)
			}
			continue
		}

		if inFunction {
			functionLines = append(functionLines, line)
			if functionEndRegex.MatchString(line) {
				inFunction = false
				c.CustomFunctions = append(c.CustomFunctions, strings.Join(functionLines, "\n"))
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
			c.Aliases[matches[1]] = matches[2]
			currentSection = "aliases"
			logger.Debug("Found alias: %s = %s", matches[1], matches[2])
		} else if matches := exportRegex.FindStringSubmatch(line); matches != nil {
			c.Exports[matches[1]] = matches[2]
			currentSection = "exports"
			logger.Debug("Found export: %s = %s", matches[1], matches[2])
		} else if matches := themeRegex.FindStringSubmatch(line); matches != nil {
			c.OhMyZshTheme = matches[1]
			currentSection = "ohmyzsh"
			logger.Debug("Found theme: %s", matches[1])
		} else if matches := pluginsRegex.FindStringSubmatch(line); matches != nil {
			pluginsStr := matches[1]
			plugins := strings.Fields(pluginsStr)
			c.OhMyZshPlugins = plugins
			currentSection = "ohmyzsh"
			logger.Debug("Found plugins: %v", plugins)
		} else {
			c.RawSections[currentSection] = append(c.RawSections[currentSection], line)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error scanning file: %v", err)
		return fmt.Errorf("error scanning file: %w", err)
	}
	
	logger.Info("Successfully loaded config: %d aliases, %d exports, %d functions", 
		len(c.Aliases), len(c.Exports), len(c.CustomFunctions))
	return nil
}

func (c *Config) Save() error {
	logger.Debug("Saving shell config to %s", c.FilePath)
	
	tempFile := c.FilePath + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		logger.Error("Failed to create temp file: %v", err)
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if c.OhMyZshTheme != "" || len(c.OhMyZshPlugins) > 0 {
		writer.WriteString("# Oh My Zsh Configuration\n")
		if c.OhMyZshTheme != "" {
			writer.WriteString(fmt.Sprintf("ZSH_THEME=\"%s\"\n", c.OhMyZshTheme))
		}
		if len(c.OhMyZshPlugins) > 0 {
			writer.WriteString(fmt.Sprintf("plugins=(%s)\n", strings.Join(c.OhMyZshPlugins, " ")))
		}
		writer.WriteString("\n")
	}

	if len(c.Exports) > 0 {
		writer.WriteString("# Environment Variables\n")
		for key, value := range c.Exports {
			// Always quote values for consistency and safety
			writer.WriteString(fmt.Sprintf("export %s=\"%s\"\n", key, value))
		}
		writer.WriteString("\n")
	}

	if len(c.Aliases) > 0 {
		writer.WriteString("# Aliases\n")
		for name, command := range c.Aliases {
			writer.WriteString(fmt.Sprintf("alias %s='%s'\n", name, command))
		}
		writer.WriteString("\n")
	}

	if len(c.CustomFunctions) > 0 {
		writer.WriteString("# Custom Functions\n")
		for _, function := range c.CustomFunctions {
			writer.WriteString(function)
			writer.WriteString("\n\n")
		}
	}

	for section, lines := range c.RawSections {
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

	if err := writer.Flush(); err != nil {
		logger.Error("Failed to flush writer: %v", err)
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	

	// Create backup of original file
	backupPath := c.FilePath + ".bak"
	if _, err := os.Stat(c.FilePath); err == nil {
		if err := os.Rename(c.FilePath, backupPath); err != nil {
			logger.Warn("Failed to create backup: %v", err)
		}
	}

	if err := os.Rename(tempFile, c.FilePath); err != nil {
		logger.Error("Failed to rename temp file: %v", err)
		// Try to restore backup
		if _, err := os.Stat(backupPath); err == nil {
			os.Rename(backupPath, c.FilePath)
		}
		return fmt.Errorf("failed to save config: %w", err)
	}
	
	logger.Info("Successfully saved config to %s", c.FilePath)
	return nil
}

func (c *Config) GetAvailableThemes() ([]string, error) {
	homeDir, _ := os.UserHomeDir()
	themesDir := filepath.Join(homeDir, ".oh-my-zsh", "themes")

	entries, err := os.ReadDir(themesDir)
	if err != nil {
		logger.Warn("Failed to read themes directory: %v", err)
		return []string{}, nil
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

func (c *Config) GetAvailablePlugins() ([]string, error) {
	homeDir, _ := os.UserHomeDir()
	pluginsDir := filepath.Join(homeDir, ".oh-my-zsh", "plugins")

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		logger.Warn("Failed to read plugins directory: %v", err)
		return []string{}, nil
	}

	plugins := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}
	return plugins, nil
}