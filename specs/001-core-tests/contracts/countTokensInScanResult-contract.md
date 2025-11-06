# Test Contract: countTokensInScanResult Function

**Function Under Test**: `countTokensInScanResult(scanResult *scanner.ScanResult, verbose bool) error`  
**Package**: `pkg/core`  
**Type**: Unexported (internal) function

## Function Signature

```go
func countTokensInScanResult(scanResult *scanner.ScanResult, verbose bool) error
```

## Contract Specifications

### Input Parameters

| Parameter | Type | Description | Constraints |
|-----------|------|-------------|-------------|
| `scanResult` | `*scanner.ScanResult` | Pointer to scan result to process | May be nil (edge case) |
| `verbose` | `bool` | Enable verbose logging output | Required |

### Output Behavior

| Condition | Expected Behavior |
|-----------|-------------------|
| Success | Returns `nil`, updates scanResult.TotalTokens and per-file TokenCounts |
| Token counter init fails | Returns error immediately, no files processed |
| All files skipped | Returns `nil`, TotalTokens remains 0 |
| Mix of valid/invalid files | Returns `nil`, processes only valid files, warns on stderr for failures |

### Side Effects

- **Writes to**: `os.Stderr` (verbose messages and warnings)
- **Reads from**: `scanResult.Files` slice
- **Modifies**: 
  - `scanResult.TotalTokens` (sets to sum of all token counts)
  - `scanResult.Files[i].TokenCount` (sets per-file counts)
- **Depends on**: 
  - `tokencounter.NewTokenCounter` (may fail)
  - `tokencounter.CountTokensWithPath` (may fail per-file)

## Test Requirements

### Functional Requirements Coverage

Tests MUST verify:

1. **FR-004**: Correctly counts tokens for all valid files in a ScanResult
   - **Test**: Create ScanResult with 2+ valid files
   - **Assert**: Each file has TokenCount > 0, TotalTokens = sum

2. **FR-005**: Skips empty files without errors
   - **Test**: Create ScanResult with empty Content
   - **Assert**: TokenCount remains 0, TotalTokens = 0, no error returned

3. **FR-006**: Skips directories without errors
   - **Test**: Create ScanResult with IsDir=true
   - **Assert**: TokenCount remains 0, TotalTokens = 0, no error returned

4. **FR-007**: Skips files with errors without failing
   - **Test**: Create ScanResult with Error != nil
   - **Assert**: TokenCount remains 0, TotalTokens excludes this file, no error returned

5. **FR-008**: Returns error if token counter initialization fails
   - **Test**: Mock scenario where NewTokenCounter fails (if possible, or document assumption)
   - **Assert**: Function returns error, no files processed

6. **FR-009**: Updates TotalTokens field in ScanResult
   - **Test**: Process valid files
   - **Assert**: scanResult.TotalTokens equals sum of all TokenCounts

7. **FR-010**: Updates TokenCount field for each processed file
   - **Test**: Check each file's TokenCount after processing
   - **Assert**: Valid files have TokenCount > 0

### Edge Cases

Tests MUST handle:

| Edge Case | Test Input | Expected Behavior |
|-----------|------------|-------------------|
| Nil ScanResult | `scanResult = nil` | Either panic (acceptable) or graceful error return |
| Empty Files slice | `scanResult.Files = []` | Return nil, TotalTokens = 0 |
| All files empty | Files with Content="" | Return nil, TotalTokens = 0 |
| Mixed valid/invalid | Some files valid, others empty/dir/error | Process valid files only, TotalTokens = sum of valid |
| Very large content | Content with 10K+ tokens | Successfully counts, no overflow |
| Special characters | Unicode, emojis in content | Successfully counts |
| Verbose enabled | `verbose = true` | stderr contains progress messages |
| Verbose disabled | `verbose = false` | stderr only has warnings (if any) |

### Test Data Specifications

**Minimum Test Cases Required**: 5 (per functional requirements)

#### Test Case 1: Valid Files
```go
Input: ScanResult{
    Files: []FileInfo{
        {Path: "file1.go", IsDir: false, Content: "package main\n", Error: nil},
        {Path: "file2.go", IsDir: false, Content: "func test() {}\n", Error: nil},
    },
}
Expected: 
- file1.TokenCount > 0
- file2.TokenCount > 0
- TotalTokens = file1.TokenCount + file2.TokenCount
- Returns nil
```

#### Test Case 2: Empty Files
```go
Input: ScanResult{
    Files: []FileInfo{
        {Path: "empty.txt", IsDir: false, Content: "", Error: nil},
    },
}
Expected:
- empty.txt TokenCount = 0
- TotalTokens = 0
- Returns nil
```

#### Test Case 3: Directories
```go
Input: ScanResult{
    Files: []FileInfo{
        {Path: "subdir", IsDir: true, Content: "", Error: nil},
    },
}
Expected:
- subdir TokenCount = 0
- TotalTokens = 0
- Returns nil
```

#### Test Case 4: Files with Errors
```go
Input: ScanResult{
    Files: []FileInfo{
        {Path: "error.txt", IsDir: false, Content: "data", Error: fmt.Errorf("permission denied")},
    },
}
Expected:
- error.txt TokenCount = 0
- TotalTokens = 0
- Returns nil (error logged but not propagated)
```

#### Test Case 5: Mixed Scenario
```go
Input: ScanResult{
    Files: []FileInfo{
        {Path: "valid.go", IsDir: false, Content: "package main\n", Error: nil},
        {Path: "empty.txt", IsDir: false, Content: "", Error: nil},
        {Path: "dir", IsDir: true, Content: "", Error: nil},
        {Path: "error.go", IsDir: false, Content: "data", Error: fmt.Errorf("err")},
    },
}
Expected:
- Only valid.go has TokenCount > 0
- TotalTokens = valid.go TokenCount
- Returns nil
```

## Test Implementation Requirements

### Assertion Strategy

Tests MUST:
1. Create ScanResult with known file configurations
2. Call `countTokensInScanResult(scanResult, verbose)`
3. Assert return value is nil (or expected error)
4. Assert `scanResult.TotalTokens` matches expected sum
5. Assert each file's `TokenCount` matches expectations
6. Verify verbose messages appear on stderr when verbose=true

### Mock/Stub Requirements

- **ScanResult**: Create inline test fixtures
- **FileInfo**: Create inline with various configurations
- **TokenCounter dependency**: Use real implementation (already tested)
  - Assumption: `tokencounter.NewTokenCounter` succeeds with default encoding
  - If initialization can fail, tests should document this scenario

### Coverage Requirements

- **Target**: 90%+ code coverage
- **Branch coverage**: All skip conditions tested (IsDir, Error, empty Content)
- **Error paths**: Token counter init failure (if testable)

## Success Criteria Mapping

| Success Criterion | Verification Method |
|-------------------|---------------------|
| SC-003: 90%+ coverage for countTokensInScanResult | Run `go test -cover` and verify coverage report |
| SC-004: Test execution under 2 seconds | Run `go test` and measure duration |
| SC-007: Clear failure messages | Test assertions include detailed error context |
| SC-008: Zero flakiness | Tests use deterministic content, no time/random dependencies |

## State Mutation Contract

### Before Call

```go
scanResult.TotalTokens == 0
scanResult.Files[i].TokenCount == 0  // for all i
```

### After Successful Call

```go
scanResult.TotalTokens == sum(valid_file_token_counts)
scanResult.Files[i].TokenCount > 0   // for valid files
scanResult.Files[i].TokenCount == 0  // for skipped files (empty/dir/error)
```

### Invariants

- If function returns error, scanResult state is undefined (may be partially modified)
- TotalTokens is always >= 0
- TotalTokens == sum of all Files[].TokenCount
- Skipped files (IsDir, Error!=nil, Content=="") always have TokenCount==0

## Interface Stability

**Breaking Changes**:
- Changing function signature
- Changing error return behavior
- Not updating TotalTokens or TokenCount fields
- Propagating file processing errors instead of logging

**Non-Breaking Changes**:
- Internal token counting implementation
- Verbose logging format
- Performance optimizations
- Additional validation

## Dependencies

**Required Packages** (for tests):
- `testing` (Go standard library)
- `github.com/BHChen24/repo2context/pkg/scanner` (for ScanResult, FileInfo types)
- `github.com/BHChen24/repo2context/pkg/tokenCounter` (used by function, already tested)

**Test Assumptions**:
- `tokenCounter.NewTokenCounter("")` succeeds with default encoding
- `tokenCounter.CountTokensWithPath` behavior is tested elsewhere
- Token counts are deterministic for same input content

## Performance Expectations

- **Per-file processing**: < 10ms for typical files (< 10KB)
- **Total execution**: Linear O(n) in number of files
- **Memory**: O(1) additional allocation (reuses existing ScanResult)

## Notes

- Function is unexported, so tests must be in same package (`package core`)
- Tests should not verify tokenCounter implementation details
- Focus on integration: does function correctly use tokenCounter and update state?
- Verbose logging format is not part of the contract (can change)
- File skipping logic is critical and must be thoroughly tested
