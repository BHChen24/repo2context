package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/BHChen24/repo2context/pkg/scanner"
)

// Helper Functions

// captureStderr captures stderr output during function execution
func captureStderr(fn func()) string {
	// Save original stderr
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()

	// Create pipe to capture output
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	os.Stderr = w

	// Execute function
	fn()

	// Close writer and read captured output
	//nolint:errcheck
	w.Close()
	var buf bytes.Buffer
	//nolint:errcheck
	io.Copy(&buf, r)

	return buf.String()
}

// createMockScanResult creates a mock ScanResult for testing
func createMockScanResult(files []scanner.FileInfo) *scanner.ScanResult {
	return &scanner.ScanResult{
		RootPath:      "/test/path",
		Files:         files,
		DirectoryTree: "test/",
		TotalFiles:    len(files),
		TotalLines:    0,
		TotalTokens:   0,
		Errors:        []string{},
	}
}

// Tests for verboseLog

func TestVerboseLog_EnabledWritesToStderr(t *testing.T) {
	output := captureStderr(func() {
		verboseLog(true, "test message")
	})

	expected := "-> test message\n"
	if output != expected {
		t.Errorf("Expected stderr output %q, got %q", expected, output)
	}
}

func TestVerboseLog_DisabledProducesNoOutput(t *testing.T) {
	output := captureStderr(func() {
		verboseLog(false, "should not appear")
	})

	if output != "" {
		t.Errorf("Expected no stderr output when verbose=false, got %q", output)
	}
}

func TestVerboseLog_FormatsMessageWithArgs(t *testing.T) {
	output := captureStderr(func() {
		verboseLog(true, "File %s has %d lines", "test.go", 42)
	})

	expected := "-> File test.go has 42 lines\n"
	if output != expected {
		t.Errorf("Expected stderr output %q, got %q", expected, output)
	}
}

func TestVerboseLog_EmptyMessage(t *testing.T) {
	output := captureStderr(func() {
		verboseLog(true, "")
	})

	expected := "-> \n"
	if output != expected {
		t.Errorf("Expected stderr output %q, got %q", expected, output)
	}
}

func TestVerboseLog_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
		message string
		args    []interface{}
		want    string
	}{
		{
			name:    "verbose enabled with message",
			verbose: true,
			message: "processing %s",
			args:    []interface{}{"file.go"},
			want:    "-> processing file.go\n",
		},
		{
			name:    "verbose disabled",
			verbose: false,
			message: "should not appear",
			args:    []interface{}{},
			want:    "",
		},
		{
			name:    "verbose enabled with empty message",
			verbose: true,
			message: "",
			args:    []interface{}{},
			want:    "-> \n",
		},
		{
			name:    "verbose enabled with multiple args",
			verbose: true,
			message: "%s: %d items, %v complete",
			args:    []interface{}{"Status", 5, true},
			want:    "-> Status: 5 items, true complete\n",
		},
		{
			name:    "verbose disabled with args",
			verbose: false,
			message: "File %s processed",
			args:    []interface{}{"test.txt"},
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStderr(func() {
				verboseLog(tt.verbose, tt.message, tt.args...)
			})

			if output != tt.want {
				t.Errorf("Test %q failed: expected %q, got %q", tt.name, tt.want, output)
			}
		})
	}
}

// Tests for countTokensInScanResult

func TestCountTokensInScanResult_ValidFiles(t *testing.T) {
	scanResult := createMockScanResult([]scanner.FileInfo{
		{
			Path:         "/test/path/file1.go",
			RelativePath: "file1.go",
			IsDir:        false,
			Content:      "package main\nfunc main() {}\n",
			TokenCount:   0,
			Error:        nil,
		},
		{
			Path:         "/test/path/file2.go",
			RelativePath: "file2.go",
			IsDir:        false,
			Content:      "package test\n",
			TokenCount:   0,
			Error:        nil,
		},
	})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if scanResult.TotalTokens == 0 {
		t.Error("Expected TotalTokens > 0, got 0")
	}

	if scanResult.Files[0].TokenCount == 0 {
		t.Error("Expected file1.go TokenCount > 0, got 0")
	}

	if scanResult.Files[1].TokenCount == 0 {
		t.Error("Expected file2.go TokenCount > 0, got 0")
	}

	expectedTotal := scanResult.Files[0].TokenCount + scanResult.Files[1].TokenCount
	if scanResult.TotalTokens != expectedTotal {
		t.Errorf("Expected TotalTokens %d, got %d", expectedTotal, scanResult.TotalTokens)
	}
}

func TestCountTokensInScanResult_SkipsEmptyFiles(t *testing.T) {
	scanResult := createMockScanResult([]scanner.FileInfo{
		{
			Path:         "/test/path/empty.txt",
			RelativePath: "empty.txt",
			IsDir:        false,
			Content:      "",
			TokenCount:   0,
			Error:        nil,
		},
	})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if scanResult.TotalTokens != 0 {
		t.Errorf("Expected TotalTokens 0 for empty file, got %d", scanResult.TotalTokens)
	}

	if scanResult.Files[0].TokenCount != 0 {
		t.Errorf("Expected empty file TokenCount 0, got %d", scanResult.Files[0].TokenCount)
	}
}

func TestCountTokensInScanResult_SkipsDirectories(t *testing.T) {
	scanResult := createMockScanResult([]scanner.FileInfo{
		{
			Path:         "/test/path/subdir",
			RelativePath: "subdir",
			IsDir:        true,
			Content:      "",
			TokenCount:   0,
			Error:        nil,
		},
	})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if scanResult.TotalTokens != 0 {
		t.Errorf("Expected TotalTokens 0 for directory, got %d", scanResult.TotalTokens)
	}

	if scanResult.Files[0].TokenCount != 0 {
		t.Errorf("Expected directory TokenCount 0, got %d", scanResult.Files[0].TokenCount)
	}
}

func TestCountTokensInScanResult_SkipsFilesWithErrors(t *testing.T) {
	scanResult := createMockScanResult([]scanner.FileInfo{
		{
			Path:         "/test/path/error.txt",
			RelativePath: "error.txt",
			IsDir:        false,
			Content:      "some content",
			TokenCount:   0,
			Error:        fmt.Errorf("permission denied"),
		},
	})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if scanResult.TotalTokens != 0 {
		t.Errorf("Expected TotalTokens 0 for file with error, got %d", scanResult.TotalTokens)
	}

	if scanResult.Files[0].TokenCount != 0 {
		t.Errorf("Expected error file TokenCount 0, got %d", scanResult.Files[0].TokenCount)
	}
}

func TestCountTokensInScanResult_MixedScenario(t *testing.T) {
	scanResult := createMockScanResult([]scanner.FileInfo{
		{
			Path:         "/test/path/valid.go",
			RelativePath: "valid.go",
			IsDir:        false,
			Content:      "package main\n",
			TokenCount:   0,
			Error:        nil,
		},
		{
			Path:         "/test/path/empty.txt",
			RelativePath: "empty.txt",
			IsDir:        false,
			Content:      "",
			TokenCount:   0,
			Error:        nil,
		},
		{
			Path:         "/test/path/subdir",
			RelativePath: "subdir",
			IsDir:        true,
			Content:      "",
			TokenCount:   0,
			Error:        nil,
		},
		{
			Path:         "/test/path/error.go",
			RelativePath: "error.go",
			IsDir:        false,
			Content:      "package test\n",
			TokenCount:   0,
			Error:        fmt.Errorf("read error"),
		},
	})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Only valid.go should be counted
	if scanResult.Files[0].TokenCount == 0 {
		t.Error("Expected valid.go TokenCount > 0, got 0")
	}

	if scanResult.TotalTokens != scanResult.Files[0].TokenCount {
		t.Errorf("Expected TotalTokens %d, got %d", scanResult.Files[0].TokenCount, scanResult.TotalTokens)
	}

	// Other files should have TokenCount = 0
	for i := 1; i < len(scanResult.Files); i++ {
		if scanResult.Files[i].TokenCount != 0 {
			t.Errorf("Expected file %d TokenCount 0, got %d", i, scanResult.Files[i].TokenCount)
		}
	}
}

func TestCountTokensInScanResult_EmptyScanResult(t *testing.T) {
	scanResult := createMockScanResult([]scanner.FileInfo{})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error for empty scan result, got %v", err)
	}

	if scanResult.TotalTokens != 0 {
		t.Errorf("Expected TotalTokens 0 for empty scan result, got %d", scanResult.TotalTokens)
	}
}

func TestCountTokensInScanResult_LargeContent(t *testing.T) {
	largeContent := strings.Repeat("word ", 10000)
	scanResult := createMockScanResult([]scanner.FileInfo{
		{
			Path:         "/test/path/large.txt",
			RelativePath: "large.txt",
			IsDir:        false,
			Content:      largeContent,
			TokenCount:   0,
			Error:        nil,
		},
	})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error for large file, got %v", err)
	}

	if scanResult.Files[0].TokenCount == 0 {
		t.Error("Expected large file TokenCount > 0, got 0")
	}

	if scanResult.TotalTokens == 0 {
		t.Error("Expected TotalTokens > 0, got 0")
	}
}

func TestCountTokensInScanResult_SpecialCharacters(t *testing.T) {
	specialContent := "Hello ?? ??\nTab:\t\nNewline:\n"
	scanResult := createMockScanResult([]scanner.FileInfo{
		{
			Path:         "/test/path/special.txt",
			RelativePath: "special.txt",
			IsDir:        false,
			Content:      specialContent,
			TokenCount:   0,
			Error:        nil,
		},
	})

	err := countTokensInScanResult(scanResult, false)
	if err != nil {
		t.Fatalf("Expected no error for special characters, got %v", err)
	}

	if scanResult.Files[0].TokenCount == 0 {
		t.Error("Expected special char file TokenCount > 0, got 0")
	}

	if scanResult.TotalTokens == 0 {
		t.Error("Expected TotalTokens > 0, got 0")
	}
}

func TestCountTokensInScanResult_TableDriven(t *testing.T) {
	tests := []struct {
		name               string
		files              []scanner.FileInfo
		expectError        bool
		expectTotalZero    bool
		expectFirstNonZero bool
	}{
		{
			name: "valid files",
			files: []scanner.FileInfo{
				{Path: "/test/file.go", RelativePath: "file.go", IsDir: false, Content: "package main\n", Error: nil},
			},
			expectError:        false,
			expectTotalZero:    false,
			expectFirstNonZero: true,
		},
		{
			name: "empty file",
			files: []scanner.FileInfo{
				{Path: "/test/empty.txt", RelativePath: "empty.txt", IsDir: false, Content: "", Error: nil},
			},
			expectError:        false,
			expectTotalZero:    true,
			expectFirstNonZero: false,
		},
		{
			name: "directory",
			files: []scanner.FileInfo{
				{Path: "/test/dir", RelativePath: "dir", IsDir: true, Content: "", Error: nil},
			},
			expectError:        false,
			expectTotalZero:    true,
			expectFirstNonZero: false,
		},
		{
			name: "file with error",
			files: []scanner.FileInfo{
				{Path: "/test/error.txt", RelativePath: "error.txt", IsDir: false, Content: "content", Error: fmt.Errorf("error")},
			},
			expectError:        false,
			expectTotalZero:    true,
			expectFirstNonZero: false,
		},
		{
			name:               "empty scan result",
			files:              []scanner.FileInfo{},
			expectError:        false,
			expectTotalZero:    true,
			expectFirstNonZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanResult := createMockScanResult(tt.files)
			err := countTokensInScanResult(scanResult, false)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tt.expectTotalZero && scanResult.TotalTokens != 0 {
				t.Errorf("Expected TotalTokens 0, got %d", scanResult.TotalTokens)
			}
			if !tt.expectTotalZero && scanResult.TotalTokens == 0 {
				t.Error("Expected TotalTokens > 0, got 0")
			}

			if len(tt.files) > 0 {
				if tt.expectFirstNonZero && scanResult.Files[0].TokenCount == 0 {
					t.Error("Expected first file TokenCount > 0, got 0")
				}
				if !tt.expectFirstNonZero && scanResult.Files[0].TokenCount != 0 {
					t.Errorf("Expected first file TokenCount 0, got %d", scanResult.Files[0].TokenCount)
				}
			}
		})
	}
}
