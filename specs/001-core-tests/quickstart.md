# Quickstart: Running Core Module Unit Tests

**Feature**: Core Module Unit Tests  
**Last Updated**: 2025-11-05

## Prerequisites

- Go 1.25.1 or later installed
- Repository cloned locally
- Working directory at repository root

## Running Tests

### Run All Core Tests

```bash
# From repository root
go test ./pkg/core/

# Expected output:
# ok      github.com/BHChen24/repo2context/pkg/core      0.XXXs
```

### Run with Verbose Output

```bash
go test -v ./pkg/core/

# Shows individual test results:
# === RUN   TestVerboseLog_EnabledWritesToStderr
# --- PASS: TestVerboseLog_EnabledWritesToStderr (0.00s)
# === RUN   TestCountTokensInScanResult_ValidFiles
# --- PASS: TestCountTokensInScanResult_ValidFiles (0.01s)
# ...
```

### Run with Coverage

```bash
# Generate coverage report
go test -cover ./pkg/core/

# Example output:
# ok      github.com/BHChen24/repo2context/pkg/core      0.XXXs  coverage: 95.2% of statements
```

### Generate Detailed Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./pkg/core/

# View coverage in browser
go tool cover -html=coverage.out

# Or view in terminal
go tool cover -func=coverage.out
```

### Run Specific Test

```bash
# Run single test by name
go test -v -run TestVerboseLog_EnabledWritesToStderr ./pkg/core/

# Run all verboseLog tests
go test -v -run TestVerboseLog ./pkg/core/

# Run all countTokensInScanResult tests
go test -v -run TestCountTokensInScanResult ./pkg/core/
```

### Run Tests Multiple Times (Verify No Flakiness)

```bash
# Run tests 10 times to check for flakiness
go test -count=10 ./pkg/core/

# All runs should pass with identical results
```

## Expected Results

### Success Criteria Verification

| Criterion | Verification Command | Expected Result |
|-----------|---------------------|-----------------|
| SC-001: All tests pass | `go test ./pkg/core/` | Exit code 0, "PASS" status |
| SC-002: 100% verboseLog coverage | `go test -cover ./pkg/core/` | verboseLog function at 100% |
| SC-003: 90%+ countTokensInScanResult coverage | `go test -cover ./pkg/core/` | countTokensInScanResult at ?90% |
| SC-004: Under 2 seconds | `go test ./pkg/core/` | Time < 2.000s |
| SC-005: 8+ test scenarios | `go test -v ./pkg/core/` | Count test functions |
| SC-008: Zero flakiness | `go test -count=10 ./pkg/core/` | All 10 runs pass |

### Sample Output

```bash
$ go test -v -cover ./pkg/core/
=== RUN   TestVerboseLog_EnabledWritesToStderr
--- PASS: TestVerboseLog_EnabledWritesToStderr (0.00s)
=== RUN   TestVerboseLog_DisabledProducesNoOutput
--- PASS: TestVerboseLog_DisabledProducesNoOutput (0.00s)
=== RUN   TestVerboseLog_FormatsMessageWithArgs
--- PASS: TestVerboseLog_FormatsMessageWithArgs (0.00s)
=== RUN   TestVerboseLog_TableDriven
=== RUN   TestVerboseLog_TableDriven/verbose_enabled
--- PASS: TestVerboseLog_TableDriven/verbose_enabled (0.00s)
=== RUN   TestVerboseLog_TableDriven/verbose_disabled
--- PASS: TestVerboseLog_TableDriven/verbose_disabled (0.00s)
--- PASS: TestVerboseLog_TableDriven (0.00s)
=== RUN   TestCountTokensInScanResult_ValidFiles
--- PASS: TestCountTokensInScanResult_ValidFiles (0.02s)
=== RUN   TestCountTokensInScanResult_SkipsEmptyFiles
--- PASS: TestCountTokensInScanResult_SkipsEmptyFiles (0.01s)
=== RUN   TestCountTokensInScanResult_SkipsDirectories
--- PASS: TestCountTokensInScanResult_SkipsDirectories (0.01s)
=== RUN   TestCountTokensInScanResult_SkipsFilesWithErrors
--- PASS: TestCountTokensInScanResult_SkipsFilesWithErrors (0.01s)
=== RUN   TestCountTokensInScanResult_MixedScenario
--- PASS: TestCountTokensInScanResult_MixedScenario (0.02s)
PASS
coverage: 95.2% of statements in github.com/BHChen24/repo2context/pkg/core
ok      github.com/BHChen24/repo2context/pkg/core      0.123s
```

## Troubleshooting

### Tests Fail to Run

**Problem**: `no Go files in /path/to/pkg/core`

**Solution**: Ensure you're in the repository root and `pkg/core/core.go` exists

```bash
pwd  # Should show repo root
ls pkg/core/core.go  # Should exist
ls pkg/core/core_test.go  # Test file should exist
```

### Import Errors

**Problem**: `cannot find package "github.com/BHChen24/repo2context/pkg/scanner"`

**Solution**: Run `go mod tidy` to download dependencies

```bash
go mod tidy
go test ./pkg/core/
```

### Coverage Below Target

**Problem**: Coverage is < 100% for verboseLog or < 90% for countTokensInScanResult

**Solution**: Check which lines are not covered

```bash
go test -coverprofile=coverage.out ./pkg/core/
go tool cover -func=coverage.out | grep -E '(verboseLog|countTokensInScanResult)'

# Add test cases for uncovered lines
```

### Tests Take Too Long

**Problem**: Test execution time > 2 seconds

**Solution**: Check for inefficient test data or unnecessary sleeps

```bash
go test -v ./pkg/core/ 2>&1 | grep -E '(PASS|FAIL)' | tail -1
# Shows total time

# Profile to find slow tests
go test -v -run TestCountTokens ./pkg/core/
```

### Flaky Tests

**Problem**: Tests pass sometimes but fail other times

**Solution**: Verify tests are deterministic

```bash
# Run tests multiple times
for i in {1..20}; do
    go test ./pkg/core/ || echo "FAILED on run $i"
done

# Check for:
# - Time dependencies (time.Now(), sleeps)
# - Random data (rand.Intn())
# - Uninitialized variables
# - Race conditions (run: go test -race ./pkg/core/)
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
- name: Run Core Tests
  run: go test -v -cover ./pkg/core/

- name: Check Coverage
  run: |
    go test -coverprofile=coverage.out ./pkg/core/
    go tool cover -func=coverage.out
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running core module tests..."
go test ./pkg/core/
if [ $? -ne 0 ]; then
    echo "Core tests failed. Commit aborted."
    exit 1
fi

echo "Tests passed!"
```

## Next Steps

After tests pass:

1. **Review Coverage**: `go tool cover -html=coverage.out`
2. **Check Edge Cases**: Verify all 6 edge cases from spec are tested
3. **Run Full Test Suite**: `go test ./...`
4. **Commit Changes**: 
   ```bash
   git add pkg/core/core_test.go
   git commit -m "test: add unit tests for verboseLog and countTokensInScanResult"
   ```

## Useful Commands Reference

```bash
# Quick test
go test ./pkg/core/

# Verbose with coverage
go test -v -cover ./pkg/core/

# HTML coverage report
go test -coverprofile=coverage.out ./pkg/core/ && go tool cover -html=coverage.out

# Run specific test
go test -v -run TestVerboseLog ./pkg/core/

# Check for race conditions
go test -race ./pkg/core/

# Benchmark (if benchmarks added)
go test -bench=. ./pkg/core/

# Test with short mode (skip slow tests if marked)
go test -short ./pkg/core/

# Clean test cache and re-run
go clean -testcache && go test ./pkg/core/
```

## Support

If tests fail unexpectedly:

1. Check Go version: `go version` (should be 1.25.1+)
2. Update dependencies: `go mod tidy`
3. Clean build cache: `go clean -cache -testcache`
4. Review test output for specific failure messages
5. Consult [spec.md](./spec.md) for expected behavior
6. Review [contracts/](./contracts/) for function contracts

## Documentation

- [spec.md](./spec.md) - Feature specification
- [plan.md](./plan.md) - Implementation plan
- [research.md](./research.md) - Testing patterns research
- [data-model.md](./data-model.md) - Test data structures
- [contracts/](./contracts/) - Function test contracts
