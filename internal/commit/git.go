package commit

import (
	"os/exec"
	"strings"
)

func getChangedFiles() ([]string, error) {
	diffFilesCmd := exec.Command("git", "diff", "HEAD", "--name-only")
	output, err := diffFilesCmd.Output()
	if err != nil {
		return []string{}, err
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	trimmedOutput := strings.TrimSpace(string(output))
	return strings.Split(trimmedOutput, "\n"), nil
}

func getDiff(file string) (string, error) {
	diffCmd := exec.Command("git", "diff", "HEAD", file)
	output, err := diffCmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
