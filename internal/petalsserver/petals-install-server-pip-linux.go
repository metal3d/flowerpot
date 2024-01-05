//go:build linux

package petalsserver

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	gitPackage  = "git+https://github.com/bigscience-workshop/petals"
	packageName = "petals"
)

func PipInstallPetals(pythonPath string) (chan string, error) {

	if err := createVirtualEnv(pythonPath); err != nil {
		return nil, err
	}

	// install petals
	cmd := exec.Command(
		pipPath(),
		"install",
		gitPackage,
	)

	outchan := make(chan string)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		outchan <- err.Error()
		defer close(outchan)
		return outchan, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		outchan <- err.Error()
		defer close(outchan)
		return outchan, err
	}

	if err := cmd.Start(); err != nil {
		outchan <- err.Error()
		defer close(outchan)
		return outchan, err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outchan <- scanner.Text()
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			outchan <- scanner.Text()
		}
	}()
	go func() {
		defer close(outchan)
		err := cmd.Wait()
		if err != nil {
			outchan <- err.Error()
		}
	}()

	return outchan, nil

}

func createVirtualEnv(python string) error {
	// create virtualenv on ~/.local/share/petalsserver
	path := os.ExpandEnv(shareDir)
	// create shareDir if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	path = filepath.Join(path, installDir)

	cmd := exec.Command(
		python,
		"-m",
		"venv",
		path,
	)
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
	cmd.Wait()
	return nil
}

func IsPetalsServerInstalled() bool {

	installDir := os.ExpandEnv(filepath.Join(shareDir, installDir))
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		log.Println("petals server not installed", shareDir, installDir)
		return false
	}

	// check if petals is installed
	cmd := exec.Command(pipPath(), "show", packageName)
	cmd.Path = filepath.Join(installDir, "bin", "pip")
	cmd.Path = os.ExpandEnv(cmd.Path)
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
	return strings.Contains(output, "Name: "+packageName)
}

func UpdatePetals() error {
	cmd := exec.Command(
		pipPath(),
		"install",
		"--upgrade",
		gitPackage,
	)
	cmd.Path = os.ExpandEnv(cmd.Path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	defer stdout.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	return cmd.Wait()

}
