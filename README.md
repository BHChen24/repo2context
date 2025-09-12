# repo2context

A CLI tool to convert the repository context to a structured markdown file that's easy for LLMs to understand.

## Special Notes

- For now the tool only allows 1 to 3 paths as input arguments.
- Please try to create a new branch from "**dev**" and name it as "test/what you are testing(or your username)" to implement your tests. 

## Description

Repo2Context is a CLI tool that analyzes your repository structure and generates a comprehensive markdown file containing all relevant context. This structured output is optimized for Large Language Models (LLMs) to understand your codebase quickly and effectively.

Goal Features:

- Scans repository structure and file contents
- Generates organized markdown with proper sections
- Allows selective file/directory inclusion with filtering
- Optimized for LLM consumption and understanding

## Current implementation

- Basic CLI setup using Cobra (such as "-h" and "-v" flags)
- Accepts multiple paths as input arguments
- Converts relative paths to absolute paths and prints them

## Installing

To install repo2context, ensure you have [https://go.dev/](GO) installed and follow these steps:

```bash
# Clone the repository
git clone https://github.com/BHChen24/repo2context.git
cd repo2context

# Build the binary
go build -o r2c

# (Optional) Move to PATH for global access
sudo mv r2c /usr/local/bin/
```

Alternatively, to run locally without installing globally:

```bash
# Build and run in the same directory
go build -o r2c
./r2c -h

# Or directly run without building
go run . -h
```

## Usage

After installation, you can run the tool from the command line. Here are some example usages:

```bash
r2c -h
r2c -v
r2c ./src
r2c ./src ./docs ./README.md
```

## TO-DO

- Separate implementation logic into dedicated packages
- Implement file/directory reading and filtering logic
- Generate structured markdown output
- Others

## License

repo2context is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
