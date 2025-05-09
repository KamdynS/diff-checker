package gitutils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GetDiff executes the 'git diff <target>' command in the specified repository directory
// and returns the diff output as a string.
func GetDiff(repoPath string, diffTarget string) (string, error) {
	if diffTarget == "" {
		diffTarget = "HEAD~1..HEAD"
	}

	cmd := exec.Command("git", "diff", diffTarget)
	cmd.Dir = repoPath // Set the working directory for the command

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff command failed: %w\nStderr: %s", err, stderr.String())
	}

	return out.String(), nil
}
