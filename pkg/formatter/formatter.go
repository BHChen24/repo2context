package formatter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BHChen24/repo2context/pkg/gitinfo"
	"github.com/BHChen24/repo2context/pkg/scanner"
)

// ContextData contains all the data for generating repository context output
type ContextData struct {
	ScanResult *scanner.ScanResult
	GitInfo    string
}

// Format generates markdown output from repository context data
func Format(data interface{}) (string, error) {
	contextData, ok := data.(*ContextData)
	if !ok {
		return "", fmt.Errorf("expected *ContextData, got %T", data)
	}

	var output strings.Builder

	// Header
	output.WriteString("# Repository Context\n\n")

	// File System Location
	output.WriteString("## File System Location\n\n")
	output.WriteString(fmt.Sprintf("%s\n\n", contextData.ScanResult.RootPath))

	// Git Info
	output.WriteString("## Git Info\n\n")
	if contextData.GitInfo != "" {
		// Format git info with proper markdown list
		gitLines := strings.Split(contextData.GitInfo, "\n")
		for _, line := range gitLines {
			if strings.TrimSpace(line) != "" {
				output.WriteString(fmt.Sprintf("- %s\n", line))
			}
		}
	} else {
		output.WriteString("- Not a git repository\n")
	}
	output.WriteString("\n")

	// Structure
	output.WriteString("## Structure\n\n")
	output.WriteString("```\n")
	if contextData.ScanResult.DirectoryTree != "" {
		output.WriteString(contextData.ScanResult.DirectoryTree)
	} else {
		output.WriteString("(empty directory)\n")
	}
	output.WriteString("```\n\n")

	// File Contents
	output.WriteString("## File Contents\n\n")

	for _, file := range contextData.ScanResult.Files {
		// Skip directories
		if file.IsDir {
			continue
		}

		// Skip files with errors
		if file.Error != nil {
			continue
		}

		// Skip empty files
		if strings.TrimSpace(file.Content) == "" {
			continue
		}

		// Write file header
		displayPath := file.RelativePath
		if displayPath == "" {
			displayPath = filepath.Base(file.Path)
		}
		output.WriteString(fmt.Sprintf("### File: %s (%d bytes)\t", displayPath, file.Size))
		
		// Write modified time
		// Refer to: https://pkg.go.dev/time
		if !file.ModTime.IsZero() {
			output.WriteString(fmt.Sprintf("(Modified: %s)\n\n", file.ModTime.Format("2006-01-02 15:04:05")))
		} else {
			output.WriteString("(Modified: unknown)\n\n")
		}

		// Determine file extension for syntax highlighting
		ext := strings.ToLower(filepath.Ext(file.Path))
		language := getLanguageFromExtension(ext)

		// Write file content with syntax highlighting
		output.WriteString(fmt.Sprintf("```%s\n", language))
		output.WriteString(file.Content)
		if !strings.HasSuffix(file.Content, "\n") {
			output.WriteString("\n")
		}

		// Write file tail
		output.WriteString("```\n\n")
	}

	// Summary
	output.WriteString("## Summary\n\n")
	output.WriteString(fmt.Sprintf("- Total files: %d\n", contextData.ScanResult.TotalFiles))
	output.WriteString(fmt.Sprintf("- Total lines: %d\n", contextData.ScanResult.TotalLines))

	// Add errors if any
	if len(contextData.ScanResult.Errors) > 0 {
		output.WriteString(fmt.Sprintf("- Errors encountered: %d\n", len(contextData.ScanResult.Errors)))
	}

	return output.String(), nil
}

// SaveToFile saves formatted data to a file
func SaveToFile(data interface{}, path string) error {
	// First format the data
	content, err := Format(data)
	if err != nil {
		return fmt.Errorf("failed to format data: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write to file
	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// NewContextData creates ContextData from scan results
func NewContextData(scanResult *scanner.ScanResult, rootPath string) (*ContextData, error) {
	// Get git information
	gitInfo, err := gitinfo.GetGitInfo(rootPath)
	if err != nil {
		// If git info fails, continue without it
		gitInfo = "Not a git repository or git not installed."
	}

	return &ContextData{
		ScanResult: scanResult,
		GitInfo:    gitInfo,
	}, nil
}

// getLanguageFromExtension returns syntax highlighting language for file extensions
// Will be used in markdown code blocks
func getLanguageFromExtension(ext string) string {
	languageMap := map[string]string{
		".go":     "go",
		".js":     "javascript",
		".ts":     "typescript",
		".py":     "python",
		".java":   "java",
		".c":      "c",
		".cpp":    "cpp",
		".h":      "c",
		".hpp":    "cpp",
		".rs":     "rust",
		".php":    "php",
		".rb":     "ruby",
		".sh":     "bash",
		".bash":   "bash",
		".zsh":    "bash",
		".fish":   "bash",
		".ps1":    "powershell",
		".html":   "html",
		".css":    "css",
		".scss":   "scss",
		".sass":   "sass",
		".json":   "json",
		".xml":    "xml",
		".yaml":   "yaml",
		".yml":    "yaml",
		".toml":   "toml",
		".ini":    "ini",
		".cfg":    "ini",
		".conf":   "ini",
		".md":     "markdown",
		".txt":    "text",
		".sql":    "sql",
		".r":      "r",
		".m":      "matlab",
		".swift":  "swift",
		".kt":     "kotlin",
		".scala":  "scala",
		".clj":    "clojure",
		".hs":     "haskell",
		".lua":    "lua",
		".vim":    "vim",
		".dockerfile": "dockerfile",
		".makefile":  "makefile",
	}

	if language, exists := languageMap[ext]; exists {
		return language
	}

	// Default to text for unknown extensions
	return "text"
}
