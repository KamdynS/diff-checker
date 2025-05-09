# Project Plan: LLM-Powered Git Diff Rule Checker

## 1. Project Goal

Create a command-line interface (CLI) tool written in Go. The user will run this tool with two flags: one specifying the path to a git repository, and another specifying the path to a directory containing markdown rule files. The tool will then:
1.  Execute a `git diff` command within the specified repository.
2.  Find and concatenate all `.md` files from the specified rules directory.
3.  Send the git diff and concatenated rules to a Language Model (LLM), specifically Gemini.
4.  Process Gemini's response and inform the user whether the diff is considered compliant with the rules.

## 2. Core Functionalities

*   **Command-Line Flags:**
    *   `--repo-path`: Path to the target git repository.
    *   `--rules-dir`: Path to the directory containing markdown rule files.
*   **Git Diff Execution:** Execute `git diff` (e.g., against `HEAD` or a configurable target) in the specified repository.
*   **Rule Loading:** Locate, read, and concatenate all markdown files from the `--rules-dir`.
*   **LLM Interaction (Gemini):**
    *   Retrieve Gemini API key (e.g., from `.env` file, specifically `GEMINI_API_KEY`).
    *   Construct a prompt containing the git diff and the concatenated rules.
    *   Send the prompt to the Gemini API (using a model with a large context window, e.g., Gemini 1.5 Flash or Pro if available via API and suitable).
    *   Receive and parse Gemini's response.
*   **Output:** Display Gemini's assessment to the user.

## 3. Key Steps & Modules

1.  **`main.go` (in `cmd/diff-checker/`)**:
    *   Main application entry point.
    *   Parse command-line flags (`--repo-path`, `--rules-dir`).
    *   Load Gemini API key.
    *   Coordinate the workflow: execute git diff, load rules, query LLM, display results.
2.  **`gitutils` package (e.g., `internal/gitutils/`)**:
    *   Function to execute `git diff` in the given repository path and return the diff output.
3.  **`rules` package (e.g., `internal/rules/`)**:
    *   Function to find, read, and concatenate all `.md` files from the given rules directory path.
4.  **`llm` package (e.g., `internal/llm/`)**:
    *   Interface for LLM interaction, specifically with Gemini.
    *   Functions to prepare the prompt.
    *   Functions to make API calls to Gemini, using the loaded API key.
5.  **Configuration**:
    *   Gemini API Key management (loading from `.env` file).
6.  **Error Handling**:
    *   Implement robust error handling and user-friendly messages throughout the application.

## 4. Project Setup and Structure

*   Initialize Go module: `go mod init <your_module_path>` (e.g., `go mod init github.com/username/diff-checker`).
*   Create directory structure:
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
    ├── .env (contains GEMINI_API_KEY, user-managed)
    ├── go.mod
    ├── go.sum (will be generated)
    └── plan.md
    ```

## 5. Input Specification

*   **Command-Line Flags**:
    *   `--repo-path <path_to_git_repo>` (required): Specifies the local file system path to the git repository.
    *   `--rules-dir <path_to_rules_directory>` (required): Specifies the local file system path to the directory containing `.md` rule files.
*   **Environment Variable**:
    *   `GEMINI_API_KEY` (required, stored in a `.env` file in the project root): The API key for accessing the Gemini LLM.

## 6. Output Specification

*   The primary output will be the LLM's assessment. This should clearly state:
    *   Whether the diff is considered compliant with the rules.
    *   Any specific reasons or sections of the diff that violate rules, as identified by the LLM.
*   Informative messages for errors or progress.

## 7. Dependencies (Initial Thoughts)

*   **Go Standard Library**: `os`, `fmt`, `flag`, `io/ioutil`, `path/filepath`, `strings`, `os/exec`.
*   **LLM Client Library**: A Go client library for Gemini (e.g., `github.com/google/generative-ai-go/genai`).
*   **`.env` file loading**: A library like `github.com/joho/godotenv` to load `GEMINI_API_KEY` from the `.env` file.
*   **CLI Argument Parsing**: Standard `flag` package or a more advanced library like `github.com/spf13/cobra`.

## 8. Milestones

1.  **Project Initialization**: Set up Go module, directory structure (including `gitutils` and `rules` again), basic `main.go`.
2.  **Flag Parsing & Config Loading**: Implement command-line flag parsing and loading `GEMINI_API_KEY` from `.env`.
3.  **Git Diff Implementation**: Implement `gitutils` to execute `git diff`.
4.  **Rule Loading Implementation**: Implement `rules` package to load and concatenate markdown files.
5.  **LLM Integration (Gemini)**: Set up Gemini client, construct prompt, and make API call.
6.  **Core Logic**: Combine all parts to get diff, rules, send to LLM, and get response.
7.  **Output Formatting**: Present LLM response clearly.
8.  **Error Handling & Refinements**: Add robust error handling, logging.

## 9. Future Enhancements (Optional)

*   Support for multiple LLM providers.
*   Caching LLM responses for identical diffs/rules.
*   More sophisticated diff parsing (e.g., targeting specific languages within the diff).
*   Interactive mode for rule selection or clarification.
*   Configuration file for advanced settings.
*   Ability to automatically fetch diffs (e.g., `git diff HEAD`).
