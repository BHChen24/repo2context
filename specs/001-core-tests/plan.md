# Implementation Plan: Core Module Unit Tests

**Branch**: `001-core-tests` | **Date**: 2025-11-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-core-tests/spec.md`

## Summary

Add comprehensive unit tests for two critical functions in `pkg/core/core.go`:
- `verboseLog`: Handles conditional stderr output for debugging
- `countTokensInScanResult`: Processes token counting across file scan results

**Technical Approach**: Use Go's standard testing framework with table-driven tests
for comprehensive scenario coverage. Capture stderr output for verboseLog validation.
Create inline mock ScanResult structures for countTokensInScanResult testing without
external mocking frameworks.

## Technical Context

**Language/Version**: Go 1.25.1  
**Primary Dependencies**: 
- Go standard library `testing` package
- Existing `tokenCounter` package (already tested)
- Existing `scanner` package types (ScanResult, FileInfo)

**Storage**: N/A (testing only)  
**Testing**: Go standard testing framework (`go test`)  
**Target Platform**: Cross-platform (Linux, macOS, Windows)
**Project Type**: Single project (CLI tool)  
**Performance Goals**: Test execution under 2 seconds for entire suite  
**Constraints**: 
- 100% coverage for verboseLog
- 90%+ coverage for countTokensInScanResult
- Zero test flakiness (deterministic results)

**Scale/Scope**: 
- 2 functions to test
- Minimum 8 test scenarios total
- 6 edge cases to cover

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle VI: Test Coverage for Core Modules

✅ **PASS** - This feature directly implements the constitutional requirement for test
coverage of core modules. The constitution states: "All core functionality modules...
MUST have unit tests. Tests MUST cover: Happy path scenarios, Edge cases, Error
conditions."

This implementation:
- Adds unit tests for `pkg/core` module functions
- Covers happy path (verbose enabled/disabled, valid files)
- Covers edge cases (nil inputs, empty data, format mismatches)
- Covers error conditions (initialization failures, files with errors)

### Principle II: Standard Streams Protocol

✅ **RELEVANT** - Tests for `verboseLog` will verify correct stderr output, ensuring
compliance with the streams protocol principle.

### Quality Gates: Before Merging

✅ **ALIGNED** - Success criteria (SC-001 to SC-008) ensure:
- All unit tests pass (`go test ./pkg/core/`)
- Coverage targets met (100% for verboseLog, 90%+ for countTokensInScanResult)
- Test execution performant (<2 seconds)

**Gate Result**: ✅ **APPROVED** - No constitution violations. Feature directly
supports constitutional requirements.

## Project Structure

### Documentation (this feature)

```text
specs/001-core-tests/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md              # Feature specification
├── research.md          # Phase 0 output (testing patterns research)
├── data-model.md        # Phase 1 output (test data structures)
├── quickstart.md        # Phase 1 output (how to run tests)
├── contracts/           # Phase 1 output (test interface contracts)
└── checklists/
    └── requirements.md  # Spec quality validation
```

### Source Code (repository root)

```text
pkg/
└── core/
    ├── core.go          # Functions under test
    └── core_test.go     # NEW: Test file to be created
```

**Structure Decision**: Single project structure with co-located tests. Tests will be
placed in `pkg/core/core_test.go` using `package core` declaration (white-box testing)
to enable direct testing of unexported functions `verboseLog` and
`countTokensInScanResult`.

## Complexity Tracking

> No constitution violations - this section is not needed.
