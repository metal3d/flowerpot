package ui

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestPythonEnv(t *testing.T) {
	Path := "${HOME}/.local/share/petals-server/bin/python"
	Path = os.ExpandEnv(Path)
	cmd := exec.Command(
		Path,
		"--version",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	log.Println(cmd.Path)
	defer stdout.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	output := buf.String()
	cmd.Wait()
	t.Log(output)
}

func TestFindPython(t *testing.T) {
	python, err := findPythonLesserThan("3.11")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(python)
}
