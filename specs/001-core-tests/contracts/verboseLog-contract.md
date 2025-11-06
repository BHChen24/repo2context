# Test Contract: verboseLog Function

**Function Under Test**: `verboseLog(verbose bool, info string, args ...interface{})`  
**Package**: `pkg/core`  
**Type**: Unexported (internal) function

## Function Signature

```go
func verboseLog(verbose bool, info string, args ...interface{})
```

## Contract Specifications

### Input Parameters

| Parameter | Type | Description | Constraints |
|-----------|------|-------------|-------------|
| `verbose` | `bool` | Flag to enable/disable logging | Required, no nil |
| `info` | `string` | Format string for message | May be empty |
| `args` | `...interface{}` | Arguments for format string | Optional, variadic |

### Output Behavior

| Condition | Expected Behavior |
|-----------|-------------------|
| `verbose == true` | MUST write formatted message to `os.Stderr` with "-> " prefix and newline |
| `verbose == false` | MUST NOT write any output to `os.Stderr` |
| Empty `info` string | MUST write only "-> \n" to stderr when verbose=true |
| Format args present | MUST format message using `fmt.Fprintf` semantics |

### Side Effects

- **Writes to**: `os.Stderr` (when verbose=true)
- **Reads from**: None
- **Modifies**: None (pure output function)
- **Depends on**: `fmt.Fprintf`, `os.Stderr`

## Test Requirements

### Functional Requirements Coverage

Tests MUST verify:

1. **FR-001**: verboseLog writes to stderr when verbose mode is enabled
   - **Test**: Call with `verbose=true` and capture stderr output
   - **Assert**: Output contains expected message with "-> " prefix

2. **FR-002**: verboseLog produces no output when verbose mode is disabled
   - **Test**: Call with `verbose=false` and capture stderr
   - **Assert**: Stderr output is empty string

3. **FR-003**: verboseLog correctly formats messages with variable arguments
   - **Test**: Call with format string containing %s, %d, etc. and matching args
   - **Assert**: Output contains properly formatted message

### Edge Cases

Tests MUST handle:

| Edge Case | Test Input | Expected Output |
|-----------|------------|-----------------|
| Empty message | `verbose=true, info="", args=[]` | `"-> \n"` |
| Format mismatch | `verbose=true, info="%s %d", args=["text"]` | Standard Go fmt behavior (may show %!d(MISSING)) |
| Nil format args | `verbose=true, info="test", args=nil` | `"-> test\n"` |
| Multiple args | `verbose=true, info="%s %s %d", args=["a","b",3]` | `"-> a b 3\n"` |

### Test Data Specifications

**Minimum Test Cases Required**: 4

#### Test Case 1: Verbose Enabled
```go
Input: {
    verbose: true,
    info: "Processing file",
    args: []interface{}{},
}
Expected: "-> Processing file\n" written to stderr
```

#### Test Case 2: Verbose Disabled
```go
Input: {
    verbose: false,
    info: "Should not appear",
    args: []interface{}{},
}
Expected: "" (empty stderr)
```

#### Test Case 3: Formatted Message
```go
Input: {
    verbose: true,
    info: "File %s has %d lines",
    args: []interface{}{"test.go", 42},
}
Expected: "-> File test.go has 42 lines\n" written to stderr
```

#### Test Case 4: Empty Message
```go
Input: {
    verbose: true,
    info: "",
    args: []interface{}{},
}
Expected: "-> \n" written to stderr
```

## Test Implementation Requirements

### Assertion Strategy

Tests MUST:
1. Capture `os.Stderr` output using pipe mechanism
2. Compare captured output with expected string
3. Restore original `os.Stderr` after test (cleanup)
4. Use `t.Helper()` in helper functions to get correct line numbers

### Mock/Stub Requirements

- **No mocks needed**: Function only depends on standard library
- **Stderr replacement**: Tests must temporarily replace `os.Stderr` with pipe
- **Cleanup**: Tests must defer restoration of original stderr

### Coverage Requirements

- **Target**: 100% code coverage
- **Branch coverage**: Both `if verbose` paths must be tested
- **Line coverage**: All lines including format and write operations

## Success Criteria Mapping

| Success Criterion | Verification Method |
|-------------------|---------------------|
| SC-002: 100% coverage for verboseLog | Run `go test -cover` and verify coverage report |
| SC-007: Clear failure messages | Test assertions must include descriptive error messages |
| SC-008: Zero flakiness | Tests must be deterministic (no time dependencies, no random data) |

## Interface Stability

**Breaking Changes**:
- Changing function signature (adding/removing parameters)
- Changing output format (removing "-> " prefix)
- Changing output destination (from stderr to stdout)

**Non-Breaking Changes**:
- Internal implementation details
- Performance optimizations
- Additional error handling

## Dependencies

**Required Packages** (for tests):
- `testing` (Go standard library)
- `os` (for Stderr manipulation)
- `io` (for pipe operations)
- `bytes` (for buffer operations)
- `strings` (for string assertions)

**No External Dependencies**: Tests use only Go standard library

## Notes

- Function is unexported, so tests must be in same package (`package core`)
- Tests should not verify internal fmt.Fprintf implementation
- Focus on behavior: correct output given inputs, not implementation details
