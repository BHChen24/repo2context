// Package scanner provides file system scanning and analysis functionality
// for the repo2context tool. It handles directory traversal, file content
// reading, binary file detection, and directory tree generation.
package scanner

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MaxFileSize defines the maximum file size (16KB) before content truncation
const MaxFileSize = 16 * 1024

// FileInfo represents information about a single file or directory
type FileInfo struct {
	Path         string // Absolute path to the file
	RelativePath string // Path relative to the scan root
	IsDir        bool   // Whether this is a directory
	Size         int64  // File size in bytes
	Content      string // File content (empty for directories and binary files)
	IsBinary     bool   // Whether the file is detected as binary
	Error        error  // Any error encountered while processing this file
}

// ScanResult contains the complete results of a directory scan operation
type ScanResult struct {
	RootPath      string     // Absolute path of the scanned root directory
	Files         []FileInfo // Information about all discovered files and directories
	DirectoryTree string     // Text representation of the directory structure
	TotalFiles    int        // Count of files (excludes directories)
	TotalLines    int        // Total lines across all text files
	Errors        []string   // List of errors encountered during scanning
}

// GetEntryPoint validates a file or directory path and returns its absolute path.
// It checks that the path exists and is accessible.
//
// Parameters:
//   - path: The file or directory path to validate
//
// Returns:
//   - string: The absolute path if valid
//   - error: Error if path is invalid or inaccessible
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

// ScanDirectory performs a complete scan of a directory and its contents.
// It recursively traverses all subdirectories, reads file contents, detects
// binary files, and generates directory structure information.
//
// Parameters:
//   - rootPath: The root directory path to scan
//
// Returns:
//   - *ScanResult: Complete scan results including files, directory tree, and statistics
//   - error: Error if the scan cannot be completed
func ScanDirectory(rootPath string) (*ScanResult, error) {
	absRoot, err := GetEntryPoint(rootPath)
	if err != nil {
		return nil, err
	}

	result := &ScanResult{
		RootPath: absRoot,
		Files:    make([]FileInfo, 0),
		Errors:   make([]string, 0),
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

		fileInfo := FileInfo{
			Path:         path,
			RelativePath: relPath,
			IsDir:        d.IsDir(),
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				fileInfo.Error = err
				result.Errors = append(result.Errors, fmt.Sprintf("error getting file info for %s: %v", path, err))
			} else {
				fileInfo.Size = info.Size()

				// Check if file is binary
				fileInfo.IsBinary = isBinaryFile(path)

				// Read content for text files
				if !fileInfo.IsBinary {
					if fileInfo.Size > MaxFileSize {
						// For large files, read first few lines
						content, lines, err := readPartialContent(path)
						if err != nil {
							fileInfo.Error = err
							result.Errors = append(result.Errors, fmt.Sprintf("error reading %s: %v", path, err))
						} else {
							fileInfo.Content = content + "\n\n[File truncated - too large to display completely]"
							result.TotalLines += lines
						}
					} else {
						// Read full content for small files
						content, lines, err := readFileContent(path)
						if err != nil {
							fileInfo.Error = err
							result.Errors = append(result.Errors, fmt.Sprintf("error reading %s: %v", path, err))
						} else {
							fileInfo.Content = content
							result.TotalLines += lines
						}
					}
				}

				result.TotalFiles++
			}
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

// Walk generates a directory tree structure for the given path.
// This is a convenience function that performs a full scan but only
// returns the directory tree string.
//
// Parameters:
//   - path: The directory path to scan
//
// Returns:
//   - string: Text representation of the directory tree
//   - error: Error if the scan fails
func Walk(path string) (string, error) {
	result, err := ScanDirectory(path)
	if err != nil {
		return "", err
	}
	return result.DirectoryTree, nil
}

// Peek reads and returns the content of a single file.
// For binary files, it returns a descriptive message instead of content.
// For large files (>16KB), it returns truncated content with a notice.
//
// Parameters:
//   - path: The file path to read
//
// Returns:
//   - string: File content or descriptive message
//   - error: Error if the file cannot be read or is a directory
func Peek(path string) (string, error) {
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

	if isBinaryFile(absPath) {
		return fmt.Sprintf("[Binary file: %s, Size: %d bytes]", filepath.Base(absPath), stat.Size()), nil
	}

	if stat.Size() > MaxFileSize {
		content, _, err := readPartialContent(absPath)
		if err != nil {
			return "", err
		}
		return content + "\n\n[File truncated - too large to display completely]", nil
	}

	content, _, err := readFileContent(absPath)
	return content, err
}

// isBinaryFile detects whether a file is binary by examining its content.
// It reads the first 512 bytes and checks for null bytes and non-printable characters.
// If more than 1% are null bytes or more than 30% are non-printable, it's considered binary.
//
// Parameters:
//   - path: The file path to check
//
// Returns:
//   - bool: true if the file is detected as binary, false otherwise
func isBinaryFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return true // Assume binary if can't open
	}
	defer file.Close()

	// Read first 512 bytes to check for binary content
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return true
	}

	// Check for null bytes or high percentage of non-printable characters
	nullBytes := 0
	nonPrintable := 0

	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			nullBytes++
		}
		if buffer[i] < 32 && buffer[i] != 9 && buffer[i] != 10 && buffer[i] != 13 {
			nonPrintable++
		}
	}

	// If more than 1% null bytes or more than 30% non-printable, consider binary
	return nullBytes > n/100 || nonPrintable > n*3/10
}

// readFileContent reads the complete content of a text file and counts lines.
// It uses a buffered scanner for efficient reading of large files.
//
// Parameters:
//   - path: The file path to read
//
// Returns:
//   - string: The complete file content
//   - int: Number of lines in the file
//   - error: Error if the file cannot be read
func readFileContent(path string) (string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteByte('\n')
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", 0, err
	}

	return content.String(), lineCount, nil
}

// readPartialContent reads the first part of a large file (up to 50 lines).
// This is used for files larger than MaxFileSize to provide a preview
// without consuming excessive memory.
//
// Parameters:
//   - path: The file path to read
//
// Returns:
//   - string: The partial file content (first 50 lines)
//   - int: Number of lines read
//   - error: Error if the file cannot be read
func readPartialContent(path string) (string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	lineCount := 0
	maxLines := 50 // Show first 50 lines for large files

	for scanner.Scan() && lineCount < maxLines {
		content.WriteString(scanner.Text())
		content.WriteByte('\n')
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", 0, err
	}

	return content.String(), lineCount, nil
}

// generateDirectoryTree creates a textual representation of the directory structure.
// It uses simple indentation (2 spaces per level) to show the hierarchy.
// Directories are marked with a trailing slash.
//
// Example output:
//   src/
//     main.go
//     utils/
//       helper.go
//   README.md
//
// Parameters:
//   - files: Slice of FileInfo containing all scanned files and directories
//   - rootPath: The root path being scanned (used for reference)
//
// Returns:
//   - string: Text representation of the directory tree
func generateDirectoryTree(files []FileInfo, rootPath string) string {
	// Build a map of all paths for easy lookup
	pathMap := make(map[string]bool)
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