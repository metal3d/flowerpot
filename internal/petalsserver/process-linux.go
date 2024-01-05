//go:build linux

package petalsserver

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

// IsPetalsServerRunning returns true if the Petals Server
// process is running.
func IsPetalsServerRunning() bool {
	cmd := exec.Command(
		"sh",
		"-c",
		`ps aux | grep `+petalsServerProcessName+` | grep -v grep`,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return false
	}
	defer stdout.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	output := buf.String()
	if err := cmd.Wait(); err != nil {
		return false
	}

	return strings.Contains(output, petalsServerProcessName)
}
