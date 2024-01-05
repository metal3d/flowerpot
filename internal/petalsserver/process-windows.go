//go:build windows
// TODO: Make it work on Windows.

package petalsserver

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

// IsPetalsServerRunning returns true if the Petals Server process is running.
func IsPetalsServerRunning() bool {
	cmd := exec.Command('wsl', '...') // BUG: this is not the right command, should use wsl
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	output := buf.String()
	cmd.Wait()
	return strings.Contains(output, petalsServerProcessName)
}
