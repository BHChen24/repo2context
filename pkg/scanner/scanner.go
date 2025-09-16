package scanner

import (
	"os/exec"
	"regexp"
)

// Walk will be responsible for scanning the directory structure.
func Walk(path string) (string, error) {
	cmd := exec.Command("tree","-F")
	output, err := cmd.Output()
	outputStr := string(output)
	re := regexp.MustCompile(`[├└│][─\s]*`)
	outputStr = re.ReplaceAllString(outputStr, "")
	if err != nil {
		return "", err
	}
	return outputStr, nil
}
