package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BHChen24/repo2context/pkg/formatter"
	"github.com/BHChen24/repo2context/pkg/scanner"
)

// The maximum number of files/directories that can be processed at once
const MaxFileLimit = 5

// Run processes the given filesystem paths and generates repository context output.
// It enforces a maximum path count (MaxFileLimit) and returns an error if that limit is exceeded.
// For each path it converts to an absolute path, verifies existence, and dispatches processing
// while forwarding respectGitignore, outputFile, and displayLineNum. Non-fatal path errors are
// written to stderr and processing continues for remaining paths; the function returns nil after
// handling all provided paths unless the initial path count check fails.
func Run(paths []string, respectGitignore bool, outputFile string, displayLineNum bool) error {
	// Check if too many files are provided
	if len(paths) > MaxFileLimit {
		return fmt.Errorf("too many files specified (%d). Maximum allowed: %d", len(paths), MaxFileLimit)
	}
	// Process each path provided
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting absolute path for '%s': %v\n", path, err)
			continue
		}

		// Check if the path exists
		if _, err := os.Stat(absPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "path does not exist: %s\n", absPath)
				continue
			}
			fmt.Fprintf(os.Stderr, "error checking path '%s': %v\n", absPath, err)
			continue
		}

		// Process the path based on whether it's a file or directory
		err = processPath(absPath, respectGitignore, outputFile, displayLineNum)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error processing path '%s': %v\n", absPath, err)
			continue
		}
	}
	return nil
}

// processPath determines whether absPath refers to a directory or a file and dispatches
// processing to processDirectory or processFile, forwarding respectGitignore, outputFile,
// and displayLineNum.
// It returns any error encountered while stating the path or from the invoked handler.
func processPath(absPath string, respectGitignore bool, outputFile string, displayLineNum bool) error {
	stat, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if stat.IsDir() {
		return processDirectory(absPath, respectGitignore, outputFile, displayLineNum)
	} else {
		return processFile(absPath, outputFile, displayLineNum)
	}
}

// processDirectory scans dirPath to build repository context data and writes the result
// to outputFile if provided or to stdout otherwise. It writes non-fatal scan warnings to
// stderr. respectGitignore controls whether .gitignore rules are honored and displayLineNum
// toggles including line numbers in file contents. The function returns an error if the
// scan, context creation, formatting, or file save operation fails.
func processDirectory(dirPath string, respectGitignore bool, outputFile string, displayLineNum bool) error {
	// Scan the directory with options
	scanResult, err := scanner.ScanDirectoryWithOptions(dirPath, scanner.ScanOptions{
		RespectGitignore: respectGitignore,
		DisplayLineNum:   displayLineNum,
	})
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	// Print any errors to stderr
	for _, errMsg := range scanResult.Errors {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", errMsg)
	}

	// Create context data
	contextData, err := formatter.NewContextData(scanResult, dirPath)
	if err != nil {
		return fmt.Errorf("failed to create context data: %w", err)
	}

	// Handle output - either to file or stdout
	if outputFile != "" {
		// Save to file
		err = formatter.SaveToFile(contextData, outputFile)
		if err != nil {
			return fmt.Errorf("failed to save to file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Output saved to: %s\n", outputFile)
	} else {
		// Format and output to stdout
		output, err := formatter.Format(contextData)
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
		fmt.Print(output)
	}

	return nil
}

// processFile builds a context for a single file and writes the formatted output to the specified output file or to stdout.
// 
// It treats the file's parent directory as the scan root, reads the file content (honoring displayLineNum), computes line
// counts and a display path, constructs a scanner.ScanResult containing one FileInfo entry, and then produces context data
// that is either saved to outputFile or printed to stdout.
// 
// Errors are returned if reading the file, obtaining file info, creating context data, saving to a file, or formatting the
// output fails.
func processFile(filePath string, outputFile string, displayLineNum bool) error {
	// For individual files, treat the parent directory as the root
	// TODO: Can be improved, don't have a clear idea now
	parentDir := filepath.Dir(filePath)

	// Read the file content
	content, err := scanner.Peek(filePath, displayLineNum)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get file info
	stat, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Create a simple scan result for this single file

	// Since this variable is calculated from the parent directory
	// and the file path, it basically is obtaining the file name
	// Keep it here for possible future enhancement
	relPath, _ := filepath.Rel(parentDir, filePath)

	// Count lines in content
	lines := 0
	if content != "" {
		for _, char := range content {
			if char == '\n' {
				lines++
			}
		}
		// Add 1 if content doesn't end with newline but has content
		if len(content) > 0 && content[len(content)-1] != '\n' {
			lines++
		}
	}

	// For single files, show the full path structure from current directory
	cwd, _ := os.Getwd()
	displayPath, _ := filepath.Rel(cwd, filePath)
	if displayPath == "" || filepath.IsAbs(displayPath) {
		// If we can't get a relative path, use the filename
		displayPath = relPath
	}

	// Construct the scan result
	scanResult := &scanner.ScanResult{
		RootPath: parentDir,
		Files: []scanner.FileInfo{
			{
				Path:         filePath,
				RelativePath: relPath,
				IsDir:        false,
				Size:         stat.Size(),
				Content:      content,
				Error:        nil,
			},
		},
		DirectoryTree: displayPath,
		TotalFiles:    1,
		TotalLines:    lines,
		Errors:        []string{},
	}

	// Create context data
	contextData, err := formatter.NewContextData(scanResult, parentDir)
	if err != nil {
		return fmt.Errorf("failed to create context data: %w", err)
	}

	// Handle output - either to file or stdout
	if outputFile != "" {
		// Save to file
		err = formatter.SaveToFile(contextData, outputFile)
		if err != nil {
			return fmt.Errorf("failed to save to file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Output saved to: %s\n", outputFile)
	} else {
		// Format and output to stdout
		output, err := formatter.Format(contextData)
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
		fmt.Print(output)
	}

	return nil
}
