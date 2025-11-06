# Tasks: Core Module Unit Tests

**Input**: Design documents from `/specs/001-core-tests/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `pkg/core/` at repository root
- Tests co-located with source: `pkg/core/core_test.go`

## Phase 1: Setup

**Purpose**: Project initialization and test file structure

- [x] T001 Create test file pkg/core/core_test.go with package core declaration
- [x] T002 Add import statements (testing, os, io, bytes, strings, fmt) to pkg/core/core_test.go
- [x] T003 Add scanner package import (github.com/BHChen24/repo2context/pkg/scanner) to pkg/core/core_test.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core test infrastructure that MUST be complete before ANY user story tests

**?? CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 [P] Implement captureStderr helper function in pkg/core/core_test.go
- [x] T005 [P] Implement createMockScanResult helper function in pkg/core/core_test.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Test verboseLog Function (Priority: P1) ?? MVP

**Goal**: Create comprehensive unit tests for verboseLog function to ensure reliable debugging output

**Independent Test**: Run `go test -v -run TestVerboseLog ./pkg/core/` - all verboseLog tests pass with 100% coverage for verboseLog function

### Implementation for User Story 1

- [x] T006 [P] [US1] Write TestVerboseLog_EnabledWritesToStderr test in pkg/core/core_test.go
- [x] T007 [P] [US1] Write TestVerboseLog_DisabledProducesNoOutput test in pkg/core/core_test.go
- [x] T008 [P] [US1] Write TestVerboseLog_FormatsMessageWithArgs test in pkg/core/core_test.go
- [x] T009 [P] [US1] Write TestVerboseLog_EmptyMessage test in pkg/core/core_test.go
- [x] T010 [US1] Create table-driven TestVerboseLog_TableDriven test with all scenarios in pkg/core/core_test.go
- [x] T011 [US1] Verify verboseLog coverage reaches 100% using go test -cover ./pkg/core/
- [x] T012 [US1] Run verboseLog tests 10 times to verify zero flakiness: go test -count=10 -run TestVerboseLog ./pkg/core/

**Checkpoint**: At this point, User Story 1 should be fully functional - verboseLog function has comprehensive test coverage and all tests pass

---

## Phase 4: User Story 2 - Test countTokensInScanResult Function (Priority: P2)

**Goal**: Create comprehensive unit tests for countTokensInScanResult to ensure reliable token counting

**Independent Test**: Run `go test -v -run TestCountTokensInScanResult ./pkg/core/` - all countTokensInScanResult tests pass with 90%+ coverage

### Implementation for User Story 2

- [x] T013 [P] [US2] Write TestCountTokensInScanResult_ValidFiles test in pkg/core/core_test.go
- [x] T014 [P] [US2] Write TestCountTokensInScanResult_SkipsEmptyFiles test in pkg/core/core_test.go
- [x] T015 [P] [US2] Write TestCountTokensInScanResult_SkipsDirectories test in pkg/core/core_test.go
- [x] T016 [P] [US2] Write TestCountTokensInScanResult_SkipsFilesWithErrors test in pkg/core/core_test.go
- [x] T017 [P] [US2] Write TestCountTokensInScanResult_InitializationFailure test (if testable) in pkg/core/core_test.go
- [x] T018 [US2] Create table-driven TestCountTokensInScanResult_TableDriven test with all scenarios in pkg/core/core_test.go
- [x] T019 [US2] Add edge case tests (nil ScanResult, empty ScanResult, large content) in pkg/core/core_test.go
- [x] T020 [US2] Add edge case tests (special characters, mixed scenarios) in pkg/core/core_test.go
- [x] T021 [US2] Verify countTokensInScanResult coverage reaches 90%+ using go test -cover ./pkg/core/
- [x] T022 [US2] Run countTokensInScanResult tests 10 times to verify zero flakiness: go test -count=10 -run TestCountTokensInScanResult ./pkg/core/

**Checkpoint**: All user stories should now be independently functional - both verboseLog and countTokensInScanResult have comprehensive test coverage

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and documentation

- [x] T023 Run full test suite and verify all tests pass: go test ./pkg/core/
- [x] T024 Generate coverage report and verify targets met: go test -coverprofile=coverage.out ./pkg/core/
- [x] T025 Verify test execution time under 2 seconds
- [x] T026 Verify at least 8 distinct test scenarios are covered (4 per function minimum)
- [x] T027 Verify all 6 edge cases from spec are covered by tests
- [x] T028 Review test failure messages for clarity and specificity
- [x] T029 [P] Update documentation if needed (README.md, quickstart.md)
- [x] T030 Run constitution check: verify Principle VI compliance (test coverage for core modules)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 ? P2)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of US1 but implements next priority

### Within Each User Story

- Tests can be written in any order (all marked [P] for parallel)
- Table-driven tests typically written after individual tests for consolidation
- Coverage verification and flakiness checks done after all tests written
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks (T001-T003) can run sequentially (same file)
- All Foundational tasks (T004-T005) marked [P] can run in parallel (different functions)
- Once Foundational phase completes, both user stories can start in parallel if team capacity allows
- Within US1: All individual tests (T006-T009) marked [P] can run in parallel
- Within US2: All individual tests (T013-T017) marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all individual tests for User Story 1 together:
Task T006: "Write TestVerboseLog_EnabledWritesToStderr test in pkg/core/core_test.go"
Task T007: "Write TestVerboseLog_DisabledProducesNoOutput test in pkg/core/core_test.go"
Task T008: "Write TestVerboseLog_FormatsMessageWithArgs test in pkg/core/core_test.go"
Task T009: "Write TestVerboseLog_EmptyMessage test in pkg/core/core_test.go"

# Then consolidate:
Task T010: "Create table-driven TestVerboseLog_TableDriven test with all scenarios"
```

---

## Parallel Example: User Story 2

```bash
# Launch all individual tests for User Story 2 together:
Task T013: "Write TestCountTokensInScanResult_ValidFiles test in pkg/core/core_test.go"
Task T014: "Write TestCountTokensInScanResult_SkipsEmptyFiles test in pkg/core/core_test.go"
Task T015: "Write TestCountTokensInScanResult_SkipsDirectories test in pkg/core/core_test.go"
Task T016: "Write TestCountTokensInScanResult_SkipsFilesWithErrors test in pkg/core/core_test.go"
Task T017: "Write TestCountTokensInScanResult_InitializationFailure test in pkg/core/core_test.go"

# Then add edge cases and consolidate:
Task T018: "Create table-driven TestCountTokensInScanResult_TableDriven test"
Task T019: "Add edge case tests (nil ScanResult, empty ScanResult, large content)"
Task T020: "Add edge case tests (special characters, mixed scenarios)"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T005) - CRITICAL blocking phase
3. Complete Phase 3: User Story 1 (T006-T012)
4. **STOP and VALIDATE**: Run `go test -v -run TestVerboseLog ./pkg/core/` and verify 100% coverage
5. Commit and potentially deploy/demo

### Incremental Delivery

1. Complete Setup + Foundational ? Foundation ready
2. Add User Story 1 ? Test independently ? Commit (MVP with verboseLog tests!)
3. Add User Story 2 ? Test independently ? Commit
4. Complete Polish phase ? Full test suite validated
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (T001-T005)
2. Once Foundational is done:
   - Developer A: User Story 1 (T006-T012)
   - Developer B: User Story 2 (T013-T022)
3. Stories complete and integrate independently
4. Team completes Polish phase together (T023-T030)

---

## Success Criteria Verification

### Per Task Validation

| Task Range | Verification Command | Expected Result |
|------------|---------------------|-----------------|
| T006-T012 (US1) | `go test -v -run TestVerboseLog ./pkg/core/` | All verboseLog tests pass |
| T011 (US1 coverage) | `go test -cover -run TestVerboseLog ./pkg/core/` | verboseLog at 100% coverage |
| T013-T022 (US2) | `go test -v -run TestCountTokensInScanResult ./pkg/core/` | All countTokensInScanResult tests pass |
| T021 (US2 coverage) | `go test -cover ./pkg/core/` | countTokensInScanResult at 90%+ coverage |
| T023 (All tests) | `go test ./pkg/core/` | Exit code 0, all tests pass |
| T025 (Performance) | Check test output time | Total time < 2.000s |

### Final Success Criteria (from spec.md)

| Criterion | Validation Task | How to Verify |
|-----------|----------------|---------------|
| SC-001: All tests pass | T023 | `go test ./pkg/core/` returns 0 |
| SC-002: 100% verboseLog coverage | T011 | `go test -cover` shows verboseLog at 100% |
| SC-003: 90%+ countTokensInScanResult | T021 | `go test -cover` shows countTokensInScanResult at ?90% |
| SC-004: Under 2 seconds | T025 | Test output shows time < 2.000s |
| SC-005: 8+ test scenarios | T026 | Count test functions (should be 10+) |
| SC-006: All edge cases covered | T027 | Map 6 edge cases from spec to test code |
| SC-007: Clear failure messages | T028 | Review test assertions for descriptive errors |
| SC-008: Zero flakiness | T012, T022 | Multiple runs all pass consistently |

---

## Notes

- [P] tasks = different test functions, can be written in parallel
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- All test tasks modify same file (pkg/core/core_test.go) but test different functions
- Commit after completing each user story phase
- Stop at any checkpoint to validate story independently
- Follow research.md patterns for implementation (white-box testing, table-driven tests, helper functions)
- Avoid: vague tasks, unclear file paths, missing story labels in US phases

---

## Test Implementation Guidance

### From research.md

- Use `package core` declaration (white-box testing)
- Implement captureStderr helper using os.Pipe pattern
- Implement createMockScanResult helper for test fixtures
- Use table-driven tests with `t.Run` for subtests
- Name tests: `Test<FunctionName>_<Scenario>`
- Include descriptive error messages in assertions

### From contracts/

- verboseLog contract: 4 minimum test cases defined
- countTokensInScanResult contract: 5 minimum test cases defined
- Review contracts for detailed input/output specifications
- Follow assertion strategies outlined in contracts

### From data-model.md

- Use test fixtures: Valid Files, Empty Files, Directories, Files with Errors, Mixed Scenario
- Reference edge case test data for special scenarios
- Follow state transition validation for countTokensInScanResult

### File Paths Reference

All tasks operate on: `/Users/parkerchen/Documents/seneca-projects/OSD/repo2context/pkg/core/core_test.go`

Functions under test located in: `/Users/parkerchen/Documents/seneca-projects/OSD/repo2context/pkg/core/core.go`
