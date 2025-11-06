# Research: Go Testing Best Practices for Core Module Tests

**Feature**: Core Module Unit Tests  
**Date**: 2025-11-05  
**Status**: Complete

## Research Questions

1. What are the best practices for testing unexported functions in Go?
2. How should we capture and verify stderr output in Go tests?
3. What patterns should we use for table-driven tests?
4. How do we create mock data structures without external mocking frameworks?
5. What are the best practices for test organization and naming?

## Findings

### 1. Testing Unexported Functions

**Decision**: Use same-package testing (`package core` in test file)

**Rationale**:
- Both `verboseLog` and `countTokensInScanResult` are unexported (lowercase) functions
- Go supports two testing approaches:
  - White-box: `package core` - can access unexported functions
  - Black-box: `package core_test` - only access exported API
- For unit testing internal functions, white-box testing is standard practice
- File location: `pkg/core/core_test.go`

**Alternatives Considered**:
- Black-box testing (`package core_test`) - Rejected because both functions are
  unexported and creating exported wrappers just for testing would pollute the API
- Separate test directory - Rejected because Go convention is co-located tests

**References**:
- Go testing documentation recommends co-located tests with same package name for
  internal unit tests
- Standard library uses this pattern extensively (e.g., `strings/strings_test.go`)

### 2. Capturing stderr Output

**Decision**: Use `bytes.Buffer` to capture `os.Stderr` writes

**Implementation Pattern**:
```go
func TestVerboseLog(t *testing.T) {
    // Save original stderr
    oldStderr := os.Stderr
    defer func() { os.Stderr = oldStderr }()
    
    // Create pipe to capture output
    r, w, _ := os.Pipe()
    os.Stderr = w
    
    // Call function under test
    verboseLog(true, "test message")
    
    // Close writer and read captured output
    w.Close()
    var buf bytes.Buffer
    io.Copy(&buf, r)
    
    // Assert output
    if !strings.Contains(buf.String(), "-> test message") {
        t.Errorf("expected message not found in stderr")
    }
}
```

**Rationale**:
- `verboseLog` writes directly to `os.Stderr` using `fmt.Fprintf`
- Go allows temporarily replacing `os.Stderr` with a pipe for testing
- This approach is standard in Go standard library tests

**Alternatives Considered**:
- Refactor `verboseLog` to accept `io.Writer` - Rejected because it would change the
  function signature just for testing
- Use `testing.T.Log` - Rejected because it doesn't capture actual stderr output

### 3. Table-Driven Test Pattern

**Decision**: Use struct slices with test case definitions

**Implementation Pattern**:
```go
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
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

**Rationale**:
- Table-driven tests reduce code duplication
- Each test case has descriptive name for clear failure messages
- `t.Run` creates subtests that can run in parallel
- Easy to add new test cases without duplicating setup/teardown logic

**Best Practices**:
- Name field comes first for readability
- Use descriptive test names that explain the scenario
- Group related test cases in same table
- Consider separate tables for different test categories (happy path vs edge cases)

### 4. Mock Data Structures

**Decision**: Create inline mock structures; no mocking framework needed

**Implementation Pattern for ScanResult**:
```go
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

func TestCountTokens_ValidFiles(t *testing.T) {
    scanResult := createMockScanResult([]scanner.FileInfo{
        {
            Path:         "/test/file1.go",
            RelativePath: "file1.go",
            IsDir:        false,
            Content:      "package main\n",
            Error:        nil,
        },
    })
    
    err := countTokensInScanResult(scanResult, false)
    // Assertions...
}
```

**Rationale**:
- `ScanResult` and `FileInfo` are simple data structures
- No complex interfaces or behaviors to mock
- Inline creation is clearer and more maintainable than mocking frameworks
- Helper functions reduce duplication across tests

**Alternatives Considered**:
- Testify mock library - Rejected as overkill for simple data structures
- Mockery code generation - Rejected as unnecessary complexity

### 5. Test Organization and Naming

**Decision**: Group tests by function, use descriptive subtest names

**File Structure**:
```go
// pkg/core/core_test.go
package core

import "testing"

// Tests for verboseLog
func TestVerboseLog_EnabledWritesToStderr(t *testing.T) { }
func TestVerboseLog_DisabledProducesNoOutput(t *testing.T) { }
func TestVerboseLog_FormatsMessageWithArgs(t *testing.T) { }
func TestVerboseLog_TableDriven(t *testing.T) { }

// Tests for countTokensInScanResult
func TestCountTokensInScanResult_ValidFiles(t *testing.T) { }
func TestCountTokensInScanResult_SkipsEmptyFiles(t *testing.T) { }
func TestCountTokensInScanResult_SkipsDirectories(t *testing.T) { }
func TestCountTokensInScanResult_SkipsFilesWithErrors(t *testing.T) { }
func TestCountTokensInScanResult_TableDriven(t *testing.T) { }

// Helper functions
func createMockScanResult(files []scanner.FileInfo) *scanner.ScanResult { }
func captureStderr(fn func()) string { }
```

**Naming Convention**:
- `Test<FunctionName>_<Scenario>` for individual tests
- Use underscores to separate function name from scenario
- Scenario names should be descriptive: "SkipsEmptyFiles" not "EmptyFiles"
- Helper functions don't need `Test` prefix

**Rationale**:
- Clear naming makes test purpose obvious
- Grouped by function for easy navigation
- Helper functions reduce duplication

## Implementation Checklist

Based on research, the implementation should:

- [ ] Create `pkg/core/core_test.go` with `package core` declaration
- [ ] Implement stderr capture helper function
- [ ] Create mock ScanResult helper function
- [ ] Write table-driven test for verboseLog (4+ scenarios)
- [ ] Write table-driven test for countTokensInScanResult (5+ scenarios)
- [ ] Write individual edge case tests for special conditions
- [ ] Ensure all tests have descriptive names
- [ ] Verify coverage meets targets (100% verboseLog, 90%+ countTokensInScanResult)
- [ ] Confirm test execution under 2 seconds
- [ ] Run tests multiple times to verify zero flakiness

## Dependencies

- No external test libraries required
- Existing packages used:
  - `testing` (Go standard library)
  - `bytes` (for buffer operations)
  - `os` (for stderr manipulation)
  - `io` (for pipe operations)
  - `strings` (for string assertions)
  - `github.com/BHChen24/repo2context/pkg/scanner` (for ScanResult types)
  - `github.com/BHChen24/repo2context/pkg/tokenCounter` (existing, already tested)
