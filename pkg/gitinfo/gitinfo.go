package gitinfo

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// This argument is duplicated too many times, so we put it in a constant
// There will be a command collector if needed in later version
const revParse = "rev-parse"

// IsGitRepository checks if a path is within a Git repository
func IsGitRepository(path string) (bool, error) {
	cmd := exec.Command("git", "-C", path, revParse, "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil, nil
}

// GetGitRoot returns the root directory of the git repository
func GetGitRoot(path string) (string, error) {
	return runGitCommand(path, revParse, "--show-toplevel")
}

// runGitCommand executes git commands in a specific directory
func runGitCommand(path string, args ...string) (string, error) {
	gitArgs := append([]string{"-C", path}, args...)
	cmd := exec.Command("git", gitArgs...)
	// For saving dynamic output
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// GetGitInfo retrieves Git information for a repository
func GetGitInfo(path string) (string, error) {
	isRepo, err := IsGitRepository(path)
	if err != nil || !isRepo {
		return "Not a git repository or git not installed.", nil
	}

	// Get commit ref
	commit, err := runGitCommand(path, "log", "-1", "--pretty=%H")
	if err != nil {
		return "", fmt.Errorf("error getting commit: %w", err)
	}

	// Get branch name
	branch, err := runGitCommand(path, revParse, "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("error getting branch: %w", err)
	}

	// Get author name
	author, err := runGitCommand(path, "log", "-1", "--pretty=%an <%ae>")
	if err != nil {
		return "", fmt.Errorf("error getting author: %w", err)
	}

	// Get date
	date, err := runGitCommand(path, "log", "-1", "--pretty=%ad")
	if err != nil {
		return "", fmt.Errorf("error getting date: %w", err)
	}

	return fmt.Sprintf("Commit: %s\nBranch: %s\nAuthor: %s\nDate  : %s", commit, branch, author, date), nil
}
