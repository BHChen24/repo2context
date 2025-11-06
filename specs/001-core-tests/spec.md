# Feature Specification: Core Module Unit Tests

**Feature Branch**: `001-core-tests`  
**Created**: 2025-11-05  
**Status**: Draft  
**Input**: User description: "Add unit tests for verboseLog and countTokensInScanResult functions in pkg/core/core.go"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Test verboseLog Function (Priority: P1)

As a developer working on the repo2context codebase, I want comprehensive unit tests
for the `verboseLog` function so that I can confidently refactor or modify verbose
logging behavior without breaking existing functionality or introducing regressions.

**Why this priority**: The `verboseLog` function is used throughout the codebase for
debugging output. Testing it ensures consistent behavior and prevents silent failures
in verbose mode that could make debugging harder for users.

**Independent Test**: Can be fully tested by creating test cases that verify message
output to stderr under different verbose flag conditions. Delivers value by ensuring
the debugging infrastructure works reliably.

**Acceptance Scenarios**:

1. **Given** verbose mode is enabled (true), **When** verboseLog is called with a message, **Then** the message is written to stderr with the "-> " prefix
2. **Given** verbose mode is disabled (false), **When** verboseLog is called with a message, **Then** no output is written to stderr
3. **Given** verbose mode is enabled and a formatted message with arguments is provided, **When** verboseLog is called, **Then** the formatted message is correctly written to stderr
4. **Given** verbose mode is enabled and an empty message is provided, **When** verboseLog is called, **Then** only the "-> " prefix and newline are written to stderr

---

### User Story 2 - Test countTokensInScanResult Function (Priority: P2)

As a developer, I want comprehensive unit tests for the `countTokensInScanResult`
function so that token counting functionality is reliable and handles edge cases
correctly, preventing incorrect token counts from being reported to users.

**Why this priority**: Token counting is a key feature for users estimating LLM
context size. Incorrect counts could lead to failed LLM requests or wasted API costs.
This builds on P1 by testing a more complex function that depends on external libraries.

**Independent Test**: Can be tested by creating mock ScanResult structures with various
file configurations and verifying that token counts are correctly calculated and stored.

**Acceptance Scenarios**:

1. **Given** a ScanResult with valid text files, **When** countTokensInScanResult is called, **Then** each file has a TokenCount value and TotalTokens is the sum of all file tokens
2. **Given** a ScanResult with empty files, **When** countTokensInScanResult is called, **Then** empty files are skipped and do not contribute to the total token count
3. **Given** a ScanResult with directories, **When** countTokensInScanResult is called, **Then** directories are skipped and only regular files are processed
4. **Given** a ScanResult with files that have errors, **When** countTokensInScanResult is called, **Then** files with errors are skipped without causing the function to fail
5. **Given** the token counter initialization fails, **When** countTokensInScanResult is called, **Then** an error is returned immediately without processing any files

---

### Edge Cases

- What happens when verboseLog is called with nil or invalid format strings?
- How does countTokensInScanResult handle files with special characters or unusual encodings?
- What happens when the ScanResult is nil or empty?
- How does countTokensInScanResult handle very large files that might overflow token counts?
- What happens when verboseLog receives format arguments that don't match the format string?
- How does countTokensInScanResult behave when verbose logging is enabled during token counting?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Test suite MUST verify that verboseLog writes to stderr when verbose mode is enabled
- **FR-002**: Test suite MUST verify that verboseLog produces no output when verbose mode is disabled
- **FR-003**: Test suite MUST verify that verboseLog correctly formats messages with variable arguments
- **FR-004**: Test suite MUST verify that countTokensInScanResult correctly counts tokens for all valid files in a ScanResult
- **FR-005**: Test suite MUST verify that countTokensInScanResult skips empty files without errors
- **FR-006**: Test suite MUST verify that countTokensInScanResult skips directories without errors
- **FR-007**: Test suite MUST verify that countTokensInScanResult skips files with errors without failing
- **FR-008**: Test suite MUST verify that countTokensInScanResult returns an error if token counter initialization fails
- **FR-009**: Test suite MUST verify that countTokensInScanResult updates the TotalTokens field in the ScanResult
- **FR-010**: Test suite MUST verify that countTokensInScanResult updates the TokenCount field for each processed file
- **FR-011**: Test suite MUST use table-driven tests where applicable to cover multiple input scenarios efficiently
- **FR-012**: Test suite MUST capture and verify stderr output for verboseLog tests
- **FR-013**: Test suite MUST create appropriate mock data structures (ScanResult, FileInfo) for countTokensInScanResult tests

### Assumptions

- **Assumption 1**: The existing tokenCounter package is already tested and can be used as a dependency
- **Assumption 2**: Tests will be placed in the same package (core_test.go) to allow testing of unexported behavior if needed
- **Assumption 3**: Go's standard testing framework is sufficient (no need for external testing libraries)
- **Assumption 4**: Mock ScanResult structures can be created inline in tests without requiring a mocking framework
- **Assumption 5**: The token counting encoding (o200k_base) is fixed and doesn't need to be tested with multiple encodings
- **Assumption 6**: Tests should verify behavior, not implementation details, to allow for future refactoring

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All test cases pass successfully when run with `go test ./pkg/core/`
- **SC-002**: Code coverage for verboseLog function reaches 100%
- **SC-003**: Code coverage for countTokensInScanResult function reaches at least 90%
- **SC-004**: Test execution completes in under 2 seconds for the entire test suite
- **SC-005**: At least 8 distinct test scenarios are covered across both functions (4 per function minimum)
- **SC-006**: All edge cases identified in the specification are covered by at least one test case
- **SC-007**: Tests provide clear failure messages that identify the specific condition that failed
- **SC-008**: Zero test flakiness - tests produce consistent results across multiple runs
