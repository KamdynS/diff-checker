package llm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/genai"
)

// Client holds the Gemini client and related configuration.
type Client struct {
	genaiClient *genai.Client
	modelName   string
	ctx         context.Context
}

// NewClient creates and returns a new Gemini client.
// It expects the API key to be passed as an argument.
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()

	// Use "gemini-1.5-flash-latest" for a balance of cost and capability with a large context window.
	// Ensure the model chosen is available for your API key and region.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new genai client: %w", err)
	}

	modelName := "gemini-2.0-flash"

	log.Printf("LLM Client initialized with model: %s", modelName)

	return &Client{
		genaiClient: client,
		modelName:   modelName,
		ctx:         ctx,
	}, nil
}

// BuildPrompt constructs the prompt string to be sent to the LLM.
func BuildPrompt(diffContent string, rulesContent string) string {
	prompt := fmt.Sprintf("Review the following git diff based on the provided coding/style rules.\n\nRULES:\n"+
		"====================\n"+
		"%s\n"+
		"====================\n\n"+
		"GIT DIFF:\n"+
		"====================\n"+
		"%s\n"+
		"====================\n\n"+
		"Based on the rules, does the git diff comply? Provide a brief explanation for your assessment. "+
		"Start your response with either \"COMPLIANT\" or \"NON-COMPLIANT\".",
		rulesContent, diffContent)
	return prompt
}

// AssessDiff sends the prompt to Gemini and returns the assessment.
func (c *Client) AssessDiff(prompt string) (string, error) {
	log.Println("Sending prompt to Gemini...")

	// Create content with the prompt as user text.
	content := genai.NewContentFromText(prompt, "user")

	resp, err := c.genaiClient.Models.GenerateContent(c.ctx, c.modelName, []*genai.Content{content}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content from Gemini: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("received an empty response or no content parts from Gemini")
	}

	var assessmentBuilder strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		assessmentBuilder.WriteString(fmt.Sprint(part))
	}

	assessment := assessmentBuilder.String()
	if assessment == "" {
		return "", fmt.Errorf("no text content received from Gemini")
	}

	log.Println("Received assessment from Gemini.")
	return assessment, nil
}
