<!--
Sync Impact Report
==================
Version Change: [INITIAL] → 1.0.0
Modification Type: Initial Ratification

Modified Principles:
- [NEW] I. CLI-First Architecture
- [NEW] II. Standard Streams Protocol
- [NEW] III. Gitignore Respect
- [NEW] IV. LLM-Optimized Output
- [NEW] V. Graceful Error Handling
- [NEW] VI. Test Coverage for Core Modules
- [NEW] VII. Simple Configuration Management

Added Sections:
- Core Principles (7 principles defined)
- Development Standards
- Quality Gates
- Governance

Removed Sections:
- None (initial version)

Template Consistency:
- ✅ plan-template.md: Constitution Check section exists, ready for validation
- ✅ spec-template.md: User scenarios align with CLI-first approach
- ✅ tasks-template.md: Task categorization supports testing discipline
- ✅ agent-file-template.md: Structure supports technology extraction
- ✅ checklist-template.md: Format supports validation checklists

Follow-up TODOs:
- None

Rationale for Version 1.0.0:
This is the initial ratification of the project constitution, establishing the
foundational governance principles for repo2context (r2c), a CLI tool for
packaging repository context for Large Language Models.
-->

# repo2context (r2c) Constitution

## Core Principles

### I. CLI-First Architecture

All functionality MUST be exposed through a command-line interface. The tool
MUST accept file and directory paths as arguments, support standard flags
(`--help`, `--version`), and provide intuitive command-line options for all
features. Configuration files (`.r2c-config.toml`) MAY supplement but NEVER
replace CLI flags.

**Rationale**: As a CLI tool, the command-line interface is the primary user
interaction surface. All features must be accessible via CLI to maintain tool
simplicity and scriptability.

### II. Standard Streams Protocol

Output MUST follow strict stream separation:

- Repository analysis results → stdout
- Error messages, warnings, and debug info → stderr
- Exit codes MUST reflect success (0) or failure (non-zero)

**Rationale**: Proper stream handling enables piping, redirection, and
integration with other Unix tools. This is fundamental to CLI tool design.

### III. Gitignore Respect

The tool MUST respect `.gitignore` files by default when scanning repositories.
Users MUST be able to override this behavior with the `--no-gitignore` flag.
When gitignore processing is enabled, files and directories matching gitignore
patterns MUST be excluded from analysis.

**Rationale**: Respecting .gitignore prevents inclusion of build artifacts,
dependencies, and temporary files that add noise to LLM context. The override
ensures users have control when needed.

### IV. LLM-Optimized Output

Output MUST be structured markdown optimized for Large Language Model
consumption with clear sections:

- File system location (absolute path)
- Git information (commit, branch, author, date)
- Directory structure tree
- File contents with syntax highlighting
- Summary statistics (file count, line count, token count)

**Rationale**: The tool's purpose is to package repository context for LLMs.
Output structure directly impacts LLM comprehension and must follow a
consistent, parseable format.

### V. Graceful Error Handling

The tool MUST handle errors gracefully:

- Permission errors: Skip file, log to stderr, continue processing
- Invalid paths: Report clear error message, exit with non-zero code
- Binary files: Exclude from content output (metadata only)
- Large files: Handle appropriately (future: truncation or summarization)

**Rationale**: Repository scanning encounters various error conditions.
Graceful handling ensures partial results are still useful and users receive
actionable error messages.

### VI. Test Coverage for Core Modules

All core functionality modules (tokenCounter, scanner, formatter, gitinfo,
gitignore) MUST have unit tests. Tests MUST cover:

- Happy path scenarios
- Edge cases (empty files, missing directories, invalid input)
- Error conditions

**Rationale**: Core modules are reused across features. Comprehensive testing
prevents regressions and ensures reliability as the codebase evolves.

### VII. Simple Configuration Management

Configuration MUST follow priority order:

1. CLI flags (highest priority)
2. `.r2c-config.toml` in current working directory
3. Built-in defaults (lowest priority)

Configuration files MUST use TOML format. All configuration options MUST be
overridable via CLI flags.

**Rationale**: Flexible configuration accommodates both one-off usage and
project-specific defaults while maintaining CLI flag supremacy for explicit
user intent.

## Development Standards

### Code Organization

- CLI logic: `cmd/root.go` (Cobra framework)
- Core functionality: `pkg/` directory with focused modules
- Tests: Co-located with source files (`*_test.go`)
- Entry point: `main.go` (minimal, delegates to cmd.Execute())

### Language and Dependencies

- Go 1.25.1 or later required
- Cobra CLI framework for command structure
- Viper for configuration management
- Minimal external dependencies preferred

### Documentation

- README.md MUST document all features and usage examples
- CLAUDE.md MUST provide AI assistant context
- Code comments MUST explain "why" not "what"

## Quality Gates

### Before Merging

- All unit tests pass (`go test ./...`)
- CLI help text accurate (`./r2c --help`)
- Manual smoke tests completed for changed features
- README updated if user-facing changes

### Release Criteria

- All required features functional
- No known critical bugs
- Version bumped appropriately
- Git tag created matching version

## Governance

### Amendment Process

1. Propose amendment with rationale
2. Document impact on existing features
3. Update version following semantic versioning
4. Update all dependent templates (plan, spec, tasks)
5. Commit with message: `docs: amend constitution to vX.Y.Z (change summary)`

### Versioning Policy

- MAJOR: Breaking principle changes (removal, fundamental redefinition)
- MINOR: New principles added, material expansions
- PATCH: Clarifications, typo fixes, non-semantic refinements

### Compliance

All PRs and code reviews MUST verify compliance with these principles. Any
deviation MUST be explicitly justified and documented. Complexity added MUST be
justified against simpler alternatives.

This constitution supersedes all other development practices and policies.

**Version**: 1.0.0 | **Ratified**: 2025-11-05 | **Last Amended**: 2025-11-05
