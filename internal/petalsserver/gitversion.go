package petalsserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

func GetInstalledGITSHA() (string, error) {
	command := []string{
		pipPath(),
		"freeze",
	}
	cmd := exec.Command(
		"sh",
		"-c",
		strings.Join(command, " "),
	)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "petals @") {
			// the version is the last part of the line
			parts := strings.Split(line, "@")
			return parts[len(parts)-1], nil
		}
	}
	return "", fmt.Errorf("no petals version found")
}

func GetLatestGitCommitSHA() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/bigscience-workshop/petals/branches/main")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var commit struct {
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}
	err = json.NewDecoder(resp.Body).Decode(&commit)
	if err != nil {
		return "", err
	}
	return commit.Commit.SHA, nil
}
