package gitinfo

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

/* 
IsGitRepository checks if the given path is within a Git repository.

git-rev-parse command is usually used to check git repo status.
Reference: https://git-scm.com/docs/git-rev-parse

Note: `-C` is for running git command in the given path parameter
*/

func IsGitRepository(path string) (bool, error) {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil, nil
}

// Git command template
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

/* 
GetGitInfo retrieves Git information for the repository at the given path.

Git-log has some special placeholders for printing pretty information.
Reference: https://git-scm.com/docs/git-log
Reference2: https://git-scm.com/docs/git-log#_pretty_formats
*/
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
	branch, err := runGitCommand(path, "rev-parse", "--abbrev-ref", "HEAD")
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

	// Sprintf still returns a string after printing it to stdout.
	return fmt.Sprintf("Commit: %s\nBranch: %s\nAuthor: %s\nDate  : %s", commit, branch, author, date), nil
}