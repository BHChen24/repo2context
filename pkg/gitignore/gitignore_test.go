package gitignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGitIgnore_WithGitignoreFile(t *testing.T) {
	// Given: A directory with .gitignore
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	content := "*.log\nnode_modules/"
	if err := os.WriteFile(gitignorePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// When: Creating GitIgnore
	gi, err := NewGitIgnore(tempDir)

	// Then: Should parse patterns correctly
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if gi == nil {
		t.Fatal("Expected non-nil GitIgnore")
	}
	if len(gi.patterns) != 2 {
		t.Errorf("Expected 2 patterns, got %d", len(gi.patterns))
	}
}

func TestIsIgnored_FilenameMatch(t *testing.T) {
	// Given: GitIgnore with a filename pattern
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("*.log"), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}
	gi, err := NewGitIgnore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create GitIgnore: %v", err)
	}

	// When/Then: Should match filename pattern
	if !gi.IsIgnored("test.log", false) {
		t.Error("Expected 'test.log' to be ignored")
	}
	if gi.IsIgnored("test.txt", false) {
		t.Error("Expected 'test.txt' to not be ignored")
	}
}
