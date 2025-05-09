package rules

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// LoadAndConcatenateRules finds all .md files in the given directory,
// reads their content, and concatenates them into a single string.
func LoadAndConcatenateRules(rulesDir string) (string, error) {
	var concatenatedRules strings.Builder

	files, err := ioutil.ReadDir(rulesDir)
	if err != nil {
		return "", fmt.Errorf("failed to read rules directory '%s': %w", rulesDir, err)
	}

	foundMdFiles := false
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
			foundMdFiles = true
			filePath := filepath.Join(rulesDir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				// Log a warning or decide if one failed file should stop the whole process
				// For now, we'll return an error
				return "", fmt.Errorf("failed to read rule file '%s': %w", filePath, err)
			}
			concatenatedRules.Write(content)
			concatenatedRules.WriteString("\n\n") // Add some separation between rule files
		}
	}

	if !foundMdFiles {
		return "", fmt.Errorf("no .md files found in rules directory '%s'", rulesDir)
	}

	return strings.TrimSpace(concatenatedRules.String()), nil
}
