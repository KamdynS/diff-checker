package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/KamdynS/diff-checker/internal/gitutils"
	"github.com/KamdynS/diff-checker/internal/llm"
	"github.com/KamdynS/diff-checker/internal/rules"

	"github.com/joho/godotenv"
	// "github.com/google/generative-ai-go/genai" // To be used in internal/llm
)

func main() {
	// Define command-line flags
	repoPath := flag.String("repo-path", ".", "Path to the git repository.")
	rulesDir := flag.String("rules-dir", "", "Path to the directory containing markdown rule files.")
	diffTarget := flag.String("diff-target", "HEAD~1..HEAD", "Git diff target for comparison (e.g., HEAD~1..HEAD, main..branch).")
	logFile := flag.String("log-file", "diff-checker.log", "Path to log file (default diff-checker.log in current directory).")

	flag.Parse()

	if *rulesDir == "" {
		fmt.Fprintln(os.Stderr, "Error: --rules-dir flag is required.")
		os.Exit(1)
	}

	// Set up logging to file _before_ any log.Printf calls
	logPath := *logFile
	if !filepath.IsAbs(logPath) {
		cwd, _ := os.Getwd()
		logPath = filepath.Join(cwd, logPath)
	}

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open log file %s, falling back to stderr: %v\n", logPath, err)
	} else {
		log.SetOutput(lf)
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables for GEMINI_API_KEY")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: GEMINI_API_KEY environment variable not set.")
		os.Exit(1)
	}

	log.Printf("Repository Path: %s", *repoPath)
	log.Printf("Rules Directory: %s", *rulesDir)
	log.Printf("Diff Target: %s", *diffTarget)
	log.Printf("Gemini API Key Loaded: %t", apiKey != "")

	// Get git diff
	log.Println("Fetching git diff...")
	diffContent, err := gitutils.GetDiff(*repoPath, *diffTarget)
	if err != nil {
		log.Fatalf("Error getting git diff: %v", err)
	}
	if diffContent == "" {
		log.Println("Git diff is empty. Nothing to check.")
		// If there's no diff, you may decide to exit early.
	}
	// For brevity, only log a snippet of the diff if it's very long
	log.Printf("Git Diff (first 200 chars if long):\n%.200s...\n", diffContent)

	// Load and concatenate rules
	log.Println("Loading rules...")
	concatenatedRules, err := rules.LoadAndConcatenateRules(*rulesDir)
	if err != nil {
		log.Fatalf("Error loading rules: %v", err)
	}
	log.Printf("Concatenated Rules (first 200 chars if long):\n%.200s...\n", concatenatedRules)

	// Initialize LLM client
	log.Println("Initializing LLM client...")
	llmClient, err := llm.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Error creating LLM client: %v", err)
	}
	// llm.Client has no explicit Close; resources are handled by the SDK.

	// Build the prompt
	log.Println("Building prompt...")
	prompt := llm.BuildPrompt(diffContent, concatenatedRules)
	// Uncomment the line below for verbose prompt debugging.
	// log.Printf("Full prompt to LLM:\n%s\n", prompt)

	fmt.Println("Analyzing diff...")
	// Get assessment from LLM
	log.Println("Sending request to LLM for assessment...")
	assessment, err := llmClient.AssessDiff(prompt)
	if err != nil {
		log.Fatalf("Error getting assessment from LLM: %v", err)
	}

	fmt.Printf("\n--- LLM Assessment ---\n%s\n----------------------\n", assessment)
	log.Println("Process completed.")
}
