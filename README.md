# Diff Checker

A command-line tool that uses Google's Gemini LLM to analyze git diffs against a set of markdown rules.

## Description

Diff Checker is a Go-based CLI tool that helps ensure code changes comply with specified rules. It works by:
1. Taking a git repository path and a directory of markdown rules as input
2. Executing a git diff in the specified repository
3. Analyzing the diff against the provided rules using Google's Gemini LLM
4. Providing feedback on whether the changes comply with the rules

## Prerequisites

- Go 1.16 or higher
- Git
- A Google Gemini API key

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd diff-checker

# Download dependencies
go mod download

# Build the project
go build -o diff-checker ./cmd/diff-checker

# (Optional) Install globally
sudo mv diff-checker /usr/local/bin/
```

Note: If you don't want to install globally, you can run the binary directly from the project directory using `./diff-checker`.

## Configuration

Create a `.env` file in the project root with your Gemini API key:
```
GEMINI_API_KEY=your_api_key_here
```

## Usage

```bash
./diff-checker --repo-path /path/to/repo --rules-dir /path/to/rules
```

### Flags

- `--repo-path`: Path to the target git repository (required)
- `--rules-dir`: Path to the directory containing markdown rule files (required)

## Project Structure

```
diff-checker/
├── cmd/
│   └── diff-checker/
│       └── main.go
├── internal/
│   ├── gitutils/
│   │   └── gitutils.go
│   ├── llm/
│   │   └── llm.go
│   └── rules/
│       └── rules.go
├── .env
├── go.mod
└── go.sum
```

## License

MIT License

Copyright (c) 2025

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

