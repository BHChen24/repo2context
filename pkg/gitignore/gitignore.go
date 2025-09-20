package gitignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GitIgnore represents a parsed .gitignore file
type GitIgnore struct {
	patterns []string
	basePath string
}

// NewGitIgnore creates a GitIgnore instance from a .gitignore file
func NewGitIgnore(basePath string) (*GitIgnore, error) {
	gi := &GitIgnore{
		basePath: basePath,
		patterns: make([]string, 0),
	}

	gitignorePath := filepath.Join(basePath, ".gitignore")

	// Check if .gitignore exists
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		// Return empty GitIgnore if no .gitignore file
		return gi, nil
	}

	// Read and parse .gitignore file
	file, err := os.Open(gitignorePath)
	if err != nil {
		return gi, err // Return empty GitIgnore on error
	}
	defer file.Close()

	bufScanner := bufio.NewScanner(file)
	for bufScanner.Scan() {
		line := strings.TrimSpace(bufScanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove leading/trailing slashes for simpler matching
		line = strings.Trim(line, "/")
		if line != "" {
			gi.patterns = append(gi.patterns, line)
		}
	}

	return gi, bufScanner.Err()
}

// IsIgnored checks if a path should be ignored based on gitignore rules
func (gi *GitIgnore) IsIgnored(relativePath string, isDir bool) bool {
	if relativePath == "" || relativePath == "." {
		return false
	}

	// Normalize path separators
	relativePath = filepath.ToSlash(relativePath)

	// Check each pattern
	for _, pattern := range gi.patterns {
		// Check exact match
		if matched, _ := filepath.Match(pattern, relativePath); matched {
			return true
		}

		// Check if filename matches pattern
		filename := filepath.Base(relativePath)
		if matched, _ := filepath.Match(pattern, filename); matched {
			return true
		}

		// Check if any path segment matches
		pathParts := strings.Split(relativePath, "/")
		for _, part := range pathParts {
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
		}
	}

	return false
}

// ShouldIgnoreFile is a convenience function that checks if a file should be ignored
// based on its absolute path and the base directory containing .gitignore
func ShouldIgnoreFile(basePath, filePath string, isDir bool) (bool, error) {
	gi, err := NewGitIgnore(basePath)
	if err != nil {
		return false, err
	}

	relativePath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return false, err
	}

	return gi.IsIgnored(relativePath, isDir), nil
}