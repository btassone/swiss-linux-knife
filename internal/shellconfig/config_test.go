package shellconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	config := New()
	
	if config == nil {
		t.Fatal("Expected non-nil config")
	}
	
	if config.Aliases == nil {
		t.Error("Expected Aliases map to be initialized")
	}
	
	if config.Exports == nil {
		t.Error("Expected Exports map to be initialized")
	}
	
	if config.OhMyZshPlugins == nil {
		t.Error("Expected OhMyZshPlugins slice to be initialized")
	}
	
	if config.CustomFunctions == nil {
		t.Error("Expected CustomFunctions slice to be initialized")
	}
	
	if config.RawSections == nil {
		t.Error("Expected RawSections map to be initialized")
	}
	
	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".zshrc")
	if config.FilePath != expectedPath {
		t.Errorf("Expected FilePath to be %s, got %s", expectedPath, config.FilePath)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	config := New()
	config.FilePath = "/tmp/nonexistent_zshrc_test_file"
	
	err := config.Load()
	if err != nil {
		t.Errorf("Expected no error for non-existent file, got %v", err)
	}
}

func TestLoadParseContent(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, ".zshrc")
	
	testContent := `# Test zshrc file
export PATH=/usr/local/bin:$PATH
export EDITOR="vim"

alias ll='ls -la'
alias gs='git status'

ZSH_THEME="robbyrussell"
plugins=(git docker kubectl)

# Custom function
hello() {
    echo "Hello, World!"
}

# Some other content
source ~/.zsh/completions`
	
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	config := New()
	config.FilePath = testFile
	
	if err := config.Load(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Test exports
	if config.Exports["PATH"] != "/usr/local/bin:$PATH" {
		t.Errorf("Expected PATH export, got %s", config.Exports["PATH"])
	}
	if config.Exports["EDITOR"] != "vim" {
		t.Errorf("Expected EDITOR export, got %s", config.Exports["EDITOR"])
	}
	
	// Test aliases
	if config.Aliases["ll"] != "ls -la" {
		t.Errorf("Expected ll alias, got %s", config.Aliases["ll"])
	}
	if config.Aliases["gs"] != "git status" {
		t.Errorf("Expected gs alias, got %s", config.Aliases["gs"])
	}
	
	// Test Oh My Zsh
	if config.OhMyZshTheme != "robbyrussell" {
		t.Errorf("Expected robbyrussell theme, got %s", config.OhMyZshTheme)
	}
	
	expectedPlugins := []string{"git", "docker", "kubectl"}
	if len(config.OhMyZshPlugins) != len(expectedPlugins) {
		t.Errorf("Expected %d plugins, got %d", len(expectedPlugins), len(config.OhMyZshPlugins))
	}
	for i, plugin := range expectedPlugins {
		if i < len(config.OhMyZshPlugins) && config.OhMyZshPlugins[i] != plugin {
			t.Errorf("Expected plugin %s at index %d, got %s", plugin, i, config.OhMyZshPlugins[i])
		}
	}
	
	// Test custom functions
	if len(config.CustomFunctions) != 1 {
		t.Errorf("Expected 1 custom function, got %d", len(config.CustomFunctions))
	}
	if len(config.CustomFunctions) > 0 && !strings.Contains(config.CustomFunctions[0], "hello()") {
		t.Errorf("Expected hello function, got %s", config.CustomFunctions[0])
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, ".zshrc")
	
	config := New()
	config.FilePath = testFile
	
	// Set up test data
	config.Exports["PATH"] = "/usr/local/bin:$PATH"
	config.Exports["EDITOR"] = "vim"
	config.Aliases["ll"] = "ls -la"
	config.Aliases["gs"] = "git status"
	config.OhMyZshTheme = "robbyrussell"
	config.OhMyZshPlugins = []string{"git", "docker"}
	config.CustomFunctions = []string{"hello() {\n    echo \"Hello\"\n}"}
	
	// Save the config
	if err := config.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Read the saved file
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}
	
	contentStr := string(content)
	
	// Verify content
	if !strings.Contains(contentStr, "export PATH=\"/usr/local/bin:$PATH\"") {
		t.Error("Expected PATH export in saved file")
	}
	if !strings.Contains(contentStr, "export EDITOR=\"vim\"") {
		t.Error("Expected EDITOR export in saved file")
	}
	if !strings.Contains(contentStr, "alias ll='ls -la'") {
		t.Error("Expected ll alias in saved file")
	}
	if !strings.Contains(contentStr, "alias gs='git status'") {
		t.Error("Expected gs alias in saved file")
	}
	if !strings.Contains(contentStr, "ZSH_THEME=\"robbyrussell\"") {
		t.Error("Expected theme in saved file")
	}
	if !strings.Contains(contentStr, "plugins=(git docker)") {
		t.Error("Expected plugins in saved file")
	}
	if !strings.Contains(contentStr, "hello()") {
		t.Error("Expected hello function in saved file")
	}
}

func TestBackupCreation(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, ".zshrc")
	backupFile := testFile + ".bak"
	
	// Create original file
	originalContent := "# Original content\nexport FOO=bar"
	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatalf("Failed to create original file: %v", err)
	}
	
	config := New()
	config.FilePath = testFile
	config.Exports["FOO"] = "baz"
	
	// Save should create backup
	if err := config.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Check backup exists
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("Expected backup file to be created")
	}
	
	// Verify backup content
	backupContent, err := os.ReadFile(backupFile)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}
	
	if string(backupContent) != originalContent {
		t.Error("Backup content doesn't match original")
	}
}