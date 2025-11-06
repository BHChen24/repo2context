# Data Model: Core Module Unit Tests

**Feature**: Core Module Unit Tests  
**Date**: 2025-11-05  
**Status**: Complete

## Overview

This document describes the test data structures and mock objects used in the core module unit tests. Since this is a testing feature rather than a production feature, the "data model" consists of test fixtures and mock data structures used to validate the behavior of `verboseLog` and `countTokensInScanResult`.

## Test Data Entities

### 1. Test Case Structure (verboseLog)

**Purpose**: Define test scenarios for verboseLog function

**Fields**:
- `name` (string): Descriptive test case name
- `verbose` (bool): Whether verbose mode is enabled
- `message` (string): Format string for the log message
- `args` ([]interface{}): Arguments to format into the message
- `expectedOutput` (string): Expected stderr output

**Example**:
```go
{
    name:           "verbose enabled with formatted message",
    verbose:        true,
    message:        "Processing file: %s",
    args:           []interface{}{"example.go"},
    expectedOutput: "-> Processing file: example.go\n",
}
```

**Validation Rules**:
- `name` must be non-empty and descriptive
- `verbose` must be explicitly set (true or false)
- `expectedOutput` should match actual stderr output including prefix and newline

### 2. Mock ScanResult

**Purpose**: Create test fixtures for countTokensInScanResult function

**Source Type**: `scanner.ScanResult` (from existing package)

**Fields Used in Tests**:
- `RootPath` (string): Base path for the scan
- `Files` ([]scanner.FileInfo): List of files in the scan
- `DirectoryTree` (string): Tree representation
- `TotalFiles` (int): Count of files
- `TotalLines` (int): Total line count
- `TotalTokens` (int): Total token count (updated by function under test)
- `Errors` ([]string): List of error messages

**Test Fixtures Required**:

#### Fixture 1: Valid Files
```go
{
    RootPath: "/test/path",
    Files: []scanner.FileInfo{
        {
            Path:         "/test/path/file1.go",
            RelativePath: "file1.go",
            IsDir:        false,
            Content:      "package main\nfunc main() {}\n",
            TokenCount:   0,  // Will be set by function
            Error:        nil,
        },
        {
            Path:         "/test/path/file2.go",
            RelativePath: "file2.go",
            IsDir:        false,
            Content:      "package test\n",
            TokenCount:   0,  // Will be set by function
            Error:        nil,
        },
    },
    TotalFiles:  2,
    TotalLines:  0,
    TotalTokens: 0,  // Will be updated by function
    Errors:      []string{},
}
```

#### Fixture 2: Empty Files
```go
{
    RootPath: "/test/path",
    Files: []scanner.FileInfo{
        {
            Path:         "/test/path/empty.txt",
            RelativePath: "empty.txt",
            IsDir:        false,
            Content:      "",  // Empty content
            TokenCount:   0,
            Error:        nil,
        },
    },
    TotalFiles:  1,
    TotalTokens: 0,  // Should remain 0
    Errors:      []string{},
}
```

#### Fixture 3: Directories
```go
{
    RootPath: "/test/path",
    Files: []scanner.FileInfo{
        {
            Path:         "/test/path/subdir",
            RelativePath: "subdir",
            IsDir:        true,  // Directory flag
            Content:      "",
            TokenCount:   0,
            Error:        nil,
        },
    },
    TotalFiles:  0,
    TotalTokens: 0,  // Should remain 0
    Errors:      []string{},
}
```

#### Fixture 4: Files with Errors
```go
{
    RootPath: "/test/path",
    Files: []scanner.FileInfo{
        {
            Path:         "/test/path/error.txt",
            RelativePath: "error.txt",
            IsDir:        false,
            Content:      "",
            TokenCount:   0,
            Error:        fmt.Errorf("permission denied"),  // Error present
        },
    },
    TotalFiles:  1,
    TotalTokens: 0,  // Should remain 0 (file skipped)
    Errors:      []string{},
}
```

#### Fixture 5: Mixed Scenario
```go
{
    RootPath: "/test/path",
    Files: []scanner.FileInfo{
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
            Content:      "",  // Should be skipped
            TokenCount:   0,
            Error:        nil,
        },
        {
            Path:         "/test/path/subdir",
            RelativePath: "subdir",
            IsDir:        true,  // Should be skipped
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
            Error:        fmt.Errorf("read error"),  // Should be skipped
        },
    },
    TotalFiles:  4,
    TotalTokens: 0,  // Should only count valid.go
    Errors:      []string{},
}
```

### 3. Mock FileInfo

**Purpose**: Individual file test data within ScanResult

**Source Type**: `scanner.FileInfo` (from existing package)

**Fields**:
- `Path` (string): Absolute file path
- `RelativePath` (string): Path relative to root
- `IsDir` (bool): Whether this is a directory
- `Size` (int64): File size in bytes
- `Content` (string): File content
- `ModTime` (time.Time): Modification time
- `TokenCount` (int): Token count (set by function under test)
- `Error` (error): Error during file processing (if any)

**Validation Rules**:
- If `IsDir` is true, `Content` should be empty
- If `Error` is non-nil, file should be skipped
- `TokenCount` starts at 0 and is set by `countTokensInScanResult`

## Test Helper Functions

### createMockScanResult

**Purpose**: Factory function to create ScanResult test fixtures

**Signature**:
```go
func createMockScanResult(files []scanner.FileInfo) *scanner.ScanResult
```

**Implementation**:
- Takes array of FileInfo objects
- Constructs complete ScanResult with defaults
- Returns pointer for mutation by function under test

### captureStderr

**Purpose**: Capture stderr output for verboseLog tests

**Signature**:
```go
func captureStderr(fn func()) string
```

**Implementation**:
- Saves original os.Stderr
- Replaces with pipe
- Executes provided function
- Captures and returns output
- Restores original stderr

## Edge Case Test Data

### Nil ScanResult
```go
var nilScanResult *scanner.ScanResult = nil
// Should handle gracefully without panic
```

### Empty ScanResult
```go
{
    RootPath:      "/test/path",
    Files:         []scanner.FileInfo{},  // Empty slice
    TotalFiles:    0,
    TotalTokens:   0,
    Errors:        []string{},
}
```

### Large Content File
```go
{
    Path:         "/test/path/large.txt",
    RelativePath: "large.txt",
    IsDir:        false,
    Content:      strings.Repeat("word ", 10000),  // Large content
    TokenCount:   0,
    Error:        nil,
}
```

### Special Characters in Content
```go
{
    Path:         "/test/path/special.txt",
    RelativePath: "special.txt",
    IsDir:        false,
    Content:      "Hello ?? ??\nTab:\t\nNewline:\n",  // Unicode & special chars
    TokenCount:   0,
    Error:        nil,
}
```

## State Transitions

### CountTokensInScanResult State Changes

**Initial State**:
- `scanResult.TotalTokens = 0`
- `scanResult.Files[i].TokenCount = 0` for all files

**After Processing**:
- `scanResult.TotalTokens = sum(valid_file_token_counts)`
- `scanResult.Files[i].TokenCount = actual_count` for valid files
- Invalid/empty/error files remain at `TokenCount = 0`

**State Diagram**:
```
[ScanResult Created]
    |
    v
[TotalTokens = 0]
[All TokenCounts = 0]
    |
    v
[countTokensInScanResult called]
    |
    v
[For each file]
    |-- IsDir? --> Skip (TokenCount remains 0)
    |-- Error != nil? --> Skip (TokenCount remains 0)
    |-- Content empty? --> Skip (TokenCount remains 0)
    |-- Valid file --> Count tokens, set TokenCount
    |
    v
[TotalTokens = sum of all non-zero TokenCounts]
```

## Relationships

```
ScanResult (1) ---contains---> (N) FileInfo
    |
    |-- TotalTokens: sum of all FileInfo.TokenCount
    |-- TotalFiles: count of FileInfo items
    |-- Errors: collection of error strings
    
FileInfo
    |-- TokenCount: set by countTokensInScanResult
    |-- Error: if present, file is skipped
    |-- IsDir: if true, file is skipped
    |-- Content: if empty, file is skipped
```

## Test Coverage Matrix

| Scenario | verboseLog | countTokensInScanResult |
|----------|-----------|------------------------|
| Happy path | verbose=true | Valid files with content |
| Disabled/Skip | verbose=false | Empty files, directories |
| Error handling | Invalid format string | Files with errors |
| Edge cases | Empty message | Nil/empty ScanResult |
| Formatting | Args mismatch | Special characters |
| Performance | N/A | Large files |

## Notes

- All test data is created inline in test functions
- No external test data files required
- Mock structures use existing scanner package types
- Tests are deterministic (no random data)
- Helper functions reduce duplication across test cases
