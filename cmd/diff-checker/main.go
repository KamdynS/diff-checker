package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	flag.Parse()

	if *rulesDir == "" {
		log.Fatal("Error: --rules-dir flag is required.")
	}

	// Load .env file. Program will still run if it doesn't exist,
	// but GEMINI_API_KEY might not be set.
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables for GEMINI_API_KEY")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: GEMINI_API_KEY environment variable not set. Please create a .env file or set it in your environment.")
	}

	log.Printf("Configuration:\n")
	log.Printf("  Repository Path: %s\n", *repoPath)
	log.Printf("  Rules Directory: %s\n", *rulesDir)
	log.Printf("  Diff Target: %s\n", *diffTarget)
	log.Printf("  Gemini API Key Loaded: %t\n", apiKey != "")

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

	// Get assessment from LLM
	log.Println("Sending request to LLM for assessment...")
	assessment, err := llmClient.AssessDiff(prompt)
	if err != nil {
		log.Fatalf("Error getting assessment from LLM: %v", err)
	}

	fmt.Printf("\n--- LLM Assessment ---\n%s\n----------------------\n", assessment)
	log.Println("Process completed.")
}
