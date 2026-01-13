package utils

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// DetectGitProject returns the current git project name, or "no-project" if not in a git repo
func DetectGitProject() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "no-project"
	}

	// Get the base name of the project directory
	projectPath := strings.TrimSpace(string(output))
	projectName := filepath.Base(projectPath)

	if projectName == "" || projectName == "." {
		return "no-project"
	}

	return projectName
}
