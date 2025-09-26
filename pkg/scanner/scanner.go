package scanner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BHChen24/repo2context/pkg/gitignore"
	"github.com/BHChen24/repo2context/pkg/gitinfo"
)

// FileInfo represents a single file or directory
type FileInfo struct {
	Path         string
	RelativePath string
	IsDir        bool
	Size         int64
	Content      string
	ModTime      time.Time
	Error        error
}

// ScanResult contains directory scan results
type ScanResult struct {
	RootPath      string
	Files         []FileInfo
	DirectoryTree string
	TotalFiles    int
	TotalLines    int
	Errors        []string
}

// ScanOptions configures directory scanning
type ScanOptions struct {
	NoGitignore    bool
	DisplayLineNum bool
}

// GetEntryPoint validates a need-to-be-processed target and returns its absolute path
func GetEntryPoint(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return "", fmt.Errorf("path does not exist: %w", err)
	}

	return absPath, nil
}

// ScanDirectory scans a directory recursively
// Ignores files/directories specified in .gitignore by default
func ScanDirectory(rootPath string) (*ScanResult, error) {
	return ScanDirectoryWithOptions(rootPath, ScanOptions{NoGitignore: false})
}

// ScanDirectoryWithOptions scans a directory with custom options
func ScanDirectoryWithOptions(rootPath string, options ScanOptions) (*ScanResult, error) {
	absRoot, err := GetEntryPoint(rootPath)
	if err != nil {
		return nil, err
	}

	result := &ScanResult{
		RootPath: absRoot,
		Files:    make([]FileInfo, 0),
		Errors:   make([]string, 0),
	}

	var gi *gitignore.GitIgnore
	var gitignoreBasePath string
	// Initialize gitignore instance if requested
	if !options.NoGitignore {
		// Try to find git repository root first
		gitRoot, gitErr := gitinfo.GetGitRoot(absRoot)
		if gitErr == nil {
			// Use git repository root if we're in a git repo
			gitignoreBasePath = gitRoot
		} else {
			// Fall back to scan directory if not in git repo
			gitignoreBasePath = absRoot
		}

		gi, err = gitignore.NewGitIgnore(gitignoreBasePath)
		if err != nil {
			// Log warning but continue without gitignore
			result.Errors = append(result.Errors, fmt.Sprintf("warning: could not load .gitignore: %v", err))
		}
	}

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errMsg := fmt.Sprintf("error accessing %s: %v", path, err)
			result.Errors = append(result.Errors, errMsg)
			return nil // Continue walking
		}

		relPath, _ := filepath.Rel(absRoot, path)
		if relPath == "." {
			relPath = ""
		}

		// Check gitignore rules if enabled
		// !options.NoGitignore => enable gitignore
		if !options.NoGitignore && gi != nil && relPath != "" {
			// Calculate relative path from gitignore base path (git root or scan directory)
			gitignoreRelPath, gitignoreRelErr := filepath.Rel(gitignoreBasePath, path)
			if gitignoreRelErr == nil && gitignoreRelPath != "." && gitignoreRelPath != "" {
				if gi.IsIgnored(gitignoreRelPath, d.IsDir()) {
					// Skip this file/directory
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
		}
 
		info, infoErr := d.Info()

		fileInfo := FileInfo{
			Path:         path,
			RelativePath: relPath,
			IsDir:        d.IsDir(),
		}

		if infoErr != nil {
			fileInfo.Error = infoErr
			result.Errors = append(result.Errors, fmt.Sprintf("error getting file info for %s: %v", path, infoErr))
		} else {
			fileInfo.ModTime = info.ModTime()
		}

		if !d.IsDir() && infoErr == nil {
			fileInfo.Size = info.Size()

			// Read file content
			content, lines, err := readFileContent(path, options.DisplayLineNum)
			if err != nil {
				fileInfo.Error = err
				result.Errors = append(result.Errors, fmt.Sprintf("error reading %s: %v", path, err))
			} else {
				fileInfo.Content = content
				result.TotalLines += lines
			}

			result.TotalFiles++
		}

		result.Files = append(result.Files, fileInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	// Generate directory tree
	result.DirectoryTree = generateDirectoryTree(result.Files, absRoot)

	return result, nil
}

// Walk returns just the directory tree structure for a path
func Walk(path string) (string, error) {
	result, err := ScanDirectory(path)
	if err != nil {
		return "", err
	}
	return result.DirectoryTree, nil
}

// Peek reads a single file's content
func Peek(path string, displayLineNum bool) (string, error) {
	absPath, err := GetEntryPoint(path)
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("error checking file: %w", err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("path is a directory, not a file")
	}

	content, _, err := readFileContent(absPath, displayLineNum)
	return content, err
}

// readFileContent reads a file's content and counts lines
func readFileContent(path string, displayLineNum bool) (string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	var content strings.Builder
	bufScanner := bufio.NewScanner(file)
	lineCount := 0

	for bufScanner.Scan() {
		// Display line number (Use tab instead of space for alignment)
		if displayLineNum {
			content.WriteString(fmt.Sprintf("%d:\t", lineCount + 1))
		}
		content.WriteString(bufScanner.Text())
		content.WriteByte('\n')
		lineCount++
	}

	if err := bufScanner.Err(); err != nil {
		return "", 0, err
	}

	return content.String(), lineCount, nil
}


// This creates a text representation of directory structure
// I am considering to move tree functions to a separate package
// (after I peeked other's repo)
func generateDirectoryTree(files []FileInfo, rootPath string) string { 
	// Build a map of all paths for easy lookup
	pathMap := make(map[string]bool) // path -> isDir
	for _, file := range files {
		if file.RelativePath != "" {
			pathMap[file.RelativePath] = file.IsDir
		}
	}

	// Get all unique directory paths and sort them
	var allPaths []string
	for path := range pathMap {
		allPaths = append(allPaths, path)
	}
	sort.Strings(allPaths)

	var result strings.Builder

	// Process each path and generate the tree structure
	processedDirs := make(map[string]bool)

	for _, path := range allPaths {
		parts := strings.Split(path, string(filepath.Separator))

		// Build each directory level
		for i := 0; i < len(parts); i++ {
			currentPath := strings.Join(parts[:i+1], string(filepath.Separator))

			if processedDirs[currentPath] {
				continue
			}

			processedDirs[currentPath] = true

			// Calculate indentation
			indent := strings.Repeat("  ", i)

			// Add directory or file
			if i == len(parts)-1 {
				// This is the actual file/directory
				if pathMap[currentPath] {
					result.WriteString(fmt.Sprintf("%s%s/\n", indent, parts[i]))
				} else {
					result.WriteString(fmt.Sprintf("%s%s\n", indent, parts[i]))
				}
			} else {
				// This is a parent directory
				result.WriteString(fmt.Sprintf("%s%s/\n", indent, parts[i]))
			}
		}
	}

	return result.String()
}
