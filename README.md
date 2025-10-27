# repo2context (r2c)

A CLI tool to convert the repository context to a structured markdown file that's easy for LLMs to understand.

## Description

Repo2Context is a CLI tool that analyzes your repository structure and generates a comprehensive markdown file containing all relevant context. This structured output is optimized for Large Language Models (LLMs) to understand your codebase quickly and effectively.

Goal Features:

- Scans repository structure and file contents
- Generates organized markdown with proper sections
- Allows selective file/directory inclusion with filtering
- Optimized for LLM consumption and understanding

## Implementation Status

### Required Features

- **Project Setup**: GitHub repository with LICENSE and comprehensive README
- **CLI Interface**: `--help`, `--version` flags and multiple path arguments
- **File Discovery**: Support for current directory, specific directories, and individual files
- **Output Format**: File system location, git info, directory tree, file contents, and statistics
- **Standard Streams**: Output to stdout, errors to stderr
- **File Reading**: Complete file content inclusion with proper formatting
- **Error Handling**: Graceful handling of permissions and inaccessible files

### Optional Features

- **Output to File**: Save results using `--output/-o` flag instead of stdout redirection
- **Gitignore Integration**: Automatic `.gitignore` respect with `--no-gitignore` override
- **TOML Configuration File**: Support for `.r2c-config.toml` in the current directory for default options
- **Token Counting**: Count tokens using OpenAI's tiktoken encoding with `--count-tokens/-t` flag

## Core Functionality

repo2context generates comprehensive markdown documentation that includes:

- **File System Location**: Absolute path of analyzed directory/file
- **Git Information**: Current commit hash, branch, author, and date (or "Not a git repository")
- **Directory Structure**: Visual tree representation of files and folders with optional per-file token counts
- **File Contents**: Complete content of all text files with syntax highlighting
- **Summary Statistics**: Total file count, line count, token count (when enabled), and any processing errors

The tool respects `.gitignore` files by default and handles permission errors gracefully.

## Installation

### Prerequisites

- [Go](https://go.dev/) installed on your system

```bash
# Clone the repository
git clone https://github.com/BHChen24/repo2context.git
cd repo2context

# Build the binary
go build -o r2c

# Run locally
./r2c [arguments]

# (Optional) Install globally
sudo mv r2c /usr/local/bin/

# Or run directly without building
go run . [arguments]
```

#### Windows Users

On Windows, you may need to build with `.exe` extension to avoid file association issues:

```bash
# Build with .exe extension
go build -o r2c.exe

# Run the executable
./r2c.exe [arguments]

# Alternative: run directly without building
go run . [arguments]
```

**Note**: If Windows shows a "How do you want to open this file?" dialog, it means the system doesn't recognize the executable. Use the `.exe` extension or run with `go run .` instead.

### Verify Installation

```bash
r2c --version
r2c --help
```

## Usage

### Basic Commands

```bash
# Show help and available flags
r2c -h

# Show version information
r2c -v

# Analyze current directory (respects .gitignore by default)
r2c .

# Analyze specific directory
r2c ./src

# Analyze specific files
r2c main.go README.md config.json

# Analyze multiple paths at once
r2c ./src ./docs ./README.md
```

### Advanced Usage

```bash
# Include files normally ignored by .gitignore
r2c --no-gitignore .

# Save output to a file
r2c ./src --output project-context.md
r2c . -o my-repo-context.md

# Count tokens in repository (useful for LLM context estimation)
r2c --count-tokens ./src
r2c -t .

# Process only specific project areas
r2c ./src ./tests ./docs

# Combine flags
r2c --no-gitignore . --output full-context.md

# Process multiple files (up to 5 files/directories)
r2c file1.go file2.go file3.go

# Use configuration file for defaults (CLI flags override)
r2c .
```

### Configuration File

Repo2context supports a TOML configuration file  which can be placed in the root directory that allows you to set default options for the tool.

### Creating a Configuration File

Create a file named `.r2c-config.toml` in the root of your project directory with the following structure:

```toml
output = "output.md"
no_gitignore = true
display_line_num = true
count_tokens = true
verbose = true
```

### Config File Behavior

- Values in the configuration file will be used as defaults if no CLI flags are provided.
- CLI flags will always override the configuration file settings.
- If no configuration file is found, the tool will use its built-in defaults.
- The file must be located in the current working directory where you run the command.

### Example usage with Config File

```bash

# Uses config defaults (e.g., verbose=true, output=output.md)
r2c .

# Overrides config (verbose=false despite config setting)
r2c --verbose=false .
```

### Flags

- `--help, -h`: Show help information
- `--version, -v`: Show version information
- `--output, -o`: Save output to file instead of stdout
- `--no-gitignore`: Disable automatic .gitignore filtering
- `--line-numbers, -l`: Include line numbers in file contents
- `--count-tokens, -t`: Count tokens using OpenAI's tiktoken 
- `--verbose`: Display detailed processing information (useful with token counting)

**Important Notes:**

- **File Limit:** Maximum of 5 files/directories can be processed in a single command to prevent performance issues and duplicate outputs. Use directory scanning for larger projects.

## Output Format

### 1. **File System Location**

Displays the absolute path of the analyzed directory or file location.

### 2. **Git Information**

- Commit hash (latest)
- Current branch name
- Author name and email
- Commit date
- Shows "Not a git repository" if outside git repo

### 3. **Directory Structure**

Visual tree representation of the file and folder hierarchy:

**Without token counting:**
```text
src/
  main.go
  utils/
    helper.go
docs/
  README.md
```

**With token counting (`-t` flag):**
```text
src/
  main.go (245 tokens)
  utils/
    helper.go (156 tokens)
docs/
  README.md (89 tokens)
```

### 4. **File Contents**

Complete content of all text files with:

- Syntax highlighting based on file extension
- Proper code formatting
- File-by-file organization

### 5. **Summary Statistics**

- Total number of files processed
- Total lines of code counted
- Total tokens (when `--count-tokens` flag is enabled)
- Number of errors encountered (if any)

**Example Summary:**
```markdown
## Summary

- Total files: 15
- Total lines: 1247
- Total tokens: 3542 (o200k_base encoding)
- Errors encountered: 0
```

## Technical Details

### File Handling

- **Encoding Support**: Handles various text encodings
- **Path Processing**: Supports both relative and absolute paths

### Gitignore Integration

- Automatically reads and respects `.gitignore` rules from git repository root
- Excludes common build artifacts, dependencies, and temporary files
- Works correctly when scanning subdirectories of a git repository
- Override with `--no-gitignore` flag when needed

### Token Counting

- Default encoding: `o200k_base`

## Testing

The project has been manually tested with comprehensive scenarios:

```bash
# Build and test the CLI tool
go build -o r2c

# Test basic functionality
./r2c --help
./r2c --version

# Test directory processing
./r2c tests/playground_for_testing/

# Test file processing
./r2c pkg/core/core.go

# Test gitignore functionality
./r2c . --no-gitignore

# Test output to file
./r2c . -o output.md

# Test token counting
./r2c -t ./src
```

**Manual Testing Coverage** includes:

- CLI argument parsing and validation
- Single file and directory processing
- Gitignore pattern filtering
- Multiple path processing
- Output formatting and file saving
- Token counting (per-file and total)
- Error handling for invalid paths and arguments

**Token Counter Testing:**

The token counter module has comprehensive unit tests:

```bash
# Run token counter tests
go test ./tests/pkg/tokenCounter/... -v

# Run with benchmarks
go test ./tests/pkg/tokenCounter/... -bench=. -benchmem

```

## License

repo2context is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
