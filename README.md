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

## Core Functionality

repo2context generates comprehensive markdown documentation that includes:

- **File System Location**: Absolute path of analyzed directory/file
- **Git Information**: Current commit hash, branch, author, and date (or "Not a git repository")
- **Directory Structure**: Visual tree representation of files and folders
- **File Contents**: Complete content of all text files with syntax highlighting
- **Summary Statistics**: Total file count, line count, and any processing errors

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

# Process only specific project areas
r2c ./src ./tests ./docs

# Combine flags
r2c --no-gitignore . --output full-context.md

# Process multiple files (up to 5 files/directories)
r2c file1.go file2.go file3.go
```

### Flags

- `--help, -h`: Show help information
- `--version, -v`: Show version information
- `--output, -o`: Save output to file instead of stdout
- `--no-gitignore`: Disable automatic .gitignore filtering

**Important Notes:**

- The `--config` flag appears in help output but is not yet implemented. Configuration file support is planned for future releases.
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

```text
src/
  main.go
  utils/
    helper.go
docs/
  README.md
```

### 4. **File Contents**

Complete content of all text files with:

- Syntax highlighting based on file extension
- Proper code formatting
- File-by-file organization

### 5. **Summary Statistics**

- Total number of files processed
- Total lines of code counted
- Number of errors encountered (if any)

## Technical Details

### File Handling

- **Encoding Support**: Handles various text encodings
- **Path Processing**: Supports both relative and absolute paths

### Gitignore Integration

- Automatically reads and respects `.gitignore` rules from git repository root
- Excludes common build artifacts, dependencies, and temporary files
- Works correctly when scanning subdirectories of a git repository
- Override with `--no-gitignore` flag when needed

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
```

**Manual Testing Coverage** includes:

- CLI argument parsing and validation
- Single file and directory processing
- Gitignore pattern filtering
- Multiple path processing
- Output formatting and file saving
- Error handling for invalid paths and arguments

## License

repo2context is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
