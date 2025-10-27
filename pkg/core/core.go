package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BHChen24/repo2context/pkg/flagConfig"
	"github.com/BHChen24/repo2context/pkg/formatter"
	"github.com/BHChen24/repo2context/pkg/scanner"
	tokencounter "github.com/BHChen24/repo2context/pkg/tokenCounter"
)

// verboseLog prints message to stderr if verbose mode is enabled
// arg: other arguments, type free
func verboseLog(verbose bool, info string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, "-> "+info+"\n", args...)
	}
}

// countTokensInScanResult counts tokens for all files in the scan result
func countTokensInScanResult(scanResult *scanner.ScanResult, verbose bool) error {
	verboseLog(verbose, "Starting token counting...")

	// Create token counter with default encoding (o200k_base)
	tc, err := tokencounter.NewTokenCounter("")
	if err != nil {
		return fmt.Errorf("failed to create token counter: %w", err)
	}

	totalTokens := 0
	fileCount := 0

	// Count tokens for each file
	for i := range scanResult.Files {
		file := &scanResult.Files[i]

		// Skip directories and files with errors
		if file.IsDir || file.Error != nil {
			continue
		}

		// Skip empty files
		if file.Content == "" {
			continue
		}

		// Count tokens
		count, err := tc.CountTokensWithPath(file.Content, file.Path)
		if err != nil {
			verboseLog(verbose, "Warning: failed to count tokens for %s: %v", file.RelativePath, err)
			continue
		}

		// Store per-file token count
		file.TokenCount = count
		totalTokens += count
		fileCount++
		verboseLog(verbose, "  %s: %d tokens", file.RelativePath, count)
	}

	scanResult.TotalTokens = totalTokens
	verboseLog(verbose, "Token counting completed - %d files, %d total tokens", fileCount, totalTokens)

	return nil
}

// Run processes paths and generates repository context output
func Run(paths []string, flagCfg flagConfig.FlagConfig) error {
	verboseLog(flagCfg.Verbose, "Starting repo2context with %d path(s)", len(paths))

	// Check if too many files are provided
	if len(paths) > 5 {
		return fmt.Errorf("too many files specified (%d). Maximum allowed: %d", len(paths), 5)
	}

	verboseLog(flagCfg.Verbose, "Processing paths: %v", paths)

	// Process each path provided
	for i, path := range paths {
		verboseLog(flagCfg.Verbose, "Processing path %d/%d: %s", i+1, len(paths), path)
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
		verboseLog(flagCfg.Verbose, "Processing absolute path: %s", absPath)
		err = processPath(absPath, flagCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error processing path '%s': %v\n", absPath, err)
			continue
		}
		verboseLog(flagCfg.Verbose, "Successfully processed: %s", absPath)
	}
	verboseLog(flagCfg.Verbose, "Completed processing all paths")
	return nil
}

// processPath handles a single file or directory
func processPath(absPath string, flagCfg flagConfig.FlagConfig) error {
	stat, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if stat.IsDir() {
		verboseLog(flagCfg.Verbose, "Detected directory: %s", absPath)
		return processDirectory(absPath, flagCfg)
	} else {
		verboseLog(flagCfg.Verbose, "Detected file: %s", absPath)
		return processFile(absPath, flagCfg)
	}
}

// processDirectory scans and formats directory output
func processDirectory(dirPath string, flagCfg flagConfig.FlagConfig) error {
	verboseLog(flagCfg.Verbose, "Starting directory scan: %s", dirPath)
	verboseLog(flagCfg.Verbose, "Scan options - NoGitignore: %t, DisplayLineNum: %t", flagCfg.NoGitignore, flagCfg.DisplayLineNum)

	// Scan the directory with options
	scanResult, err := scanner.ScanDirectoryWithOptions(dirPath, scanner.ScanOptions{
		NoGitignore:    flagCfg.NoGitignore,
		DisplayLineNum: flagCfg.DisplayLineNum,
	})
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	verboseLog(flagCfg.Verbose, "Directory scan completed - Found %d files, %d total lines", scanResult.TotalFiles, scanResult.TotalLines)

	// Print any errors to stderr
	for _, errMsg := range scanResult.Errors {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", errMsg)
	}

	// Count tokens if flag is enabled
	if flagCfg.CountTokens {
		if err := countTokensInScanResult(scanResult, flagCfg.Verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: token counting failed: %v\n", err)
		}
		// Regenerate directory tree with token counts
		scanResult.DirectoryTree = scanner.RegenerateDirectoryTree(scanResult)
	}

	verboseLog(flagCfg.Verbose, "Creating context data for formatting")
	// Create context data
	contextData, err := formatter.NewContextData(scanResult, dirPath)
	if err != nil {
		return fmt.Errorf("failed to create context data: %w", err)
	}

	// Handle output - either to file or stdout
	if flagCfg.OutputFile != "" {
		verboseLog(flagCfg.Verbose, "Saving output to file: %s", flagCfg.OutputFile)
		// Save to file
		err = formatter.SaveToFile(contextData, flagCfg.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to save to file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Output saved to: %s\n", flagCfg.OutputFile)
		verboseLog(flagCfg.Verbose, "File saved successfully")
	} else {
		verboseLog(flagCfg.Verbose, "Formatting output for stdout")
		// Format and output to stdout
		output, err := formatter.Format(contextData)
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
		verboseLog(flagCfg.Verbose, "Output formatted, writing to stdout")
		fmt.Print(output)
	}

	return nil
}

func countLines(content string) int {
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
	return lines
}

// processFile handles individual file output
func processFile(filePath string, flagCfg flagConfig.FlagConfig) error {
	// For individual files, treat the parent directory as the root
	parentDir := filepath.Dir(filePath)

	// Read the file content
	content, err := scanner.Peek(filePath, flagCfg.DisplayLineNum)
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
	lines := countLines(content)

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
				ModTime:      stat.ModTime(),
				Error:        nil,
			},
		},
		DirectoryTree: displayPath,
		TotalFiles:    1,
		TotalLines:    lines,
		Errors:        []string{},
	}

	// Count tokens if flag is enabled
	if flagCfg.CountTokens {
		if err := countTokensInScanResult(scanResult, flagCfg.Verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: token counting failed: %v\n", err)
		}
		// Regenerate directory tree with token counts
		scanResult.DirectoryTree = scanner.RegenerateDirectoryTree(scanResult)
	}

	// Create context data
	contextData, err := formatter.NewContextData(scanResult, parentDir)
	if err != nil {
		return fmt.Errorf("failed to create context data: %w", err)
	}

	// Handle output - either to file or stdout
	if flagCfg.OutputFile != "" {
		verboseLog(flagCfg.Verbose, "Saving output to file: %s", flagCfg.OutputFile)
		// Save to file
		err = formatter.SaveToFile(contextData, flagCfg.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to save to file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Output saved to: %s\n", flagCfg.OutputFile)
		verboseLog(flagCfg.Verbose, "File saved successfully")
	} else {
		verboseLog(flagCfg.Verbose, "Formatting output for stdout")
		// Format and output to stdout
		output, err := formatter.Format(contextData)
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
		verboseLog(flagCfg.Verbose, "Output formatted, writing to stdout")
		fmt.Print(output)
	}

	return nil
}
