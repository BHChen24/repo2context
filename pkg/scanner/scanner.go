package scanner

import (
	"os/exec"
	"regexp"
)

func GetEntryPoint(path string) (string, error) {
	// TODO
	return path, nil
}

func Walk(path string) (string, error) {
	// TODO: Replace implementation with std.
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

func Peek(path string) (string, error) {
	// TODO: Implement file content peeking logic.
	return "", nil
}