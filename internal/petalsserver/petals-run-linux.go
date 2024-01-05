//go:build linux

package petalsserver

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func LaunchPetalsServer(options *RunOptions, outchan chan []byte) (context.CancelFunc, *exec.Cmd, error) {
	// python -m petals.cli.run_server petals-team/StableBeluga2
	if options == nil {
		options = &RunOptions{}
	}
	if options.ModelName == "" {
		options.ModelName = "petals-team/StableBeluga2"
	}

	// construct command
	command := []string{
		pythonPath(),
		"-m",
		"petals.cli.run_server",
		options.ModelName,
	}
	if options.PublicName != "" {
		command = append(command, fmt.Sprintf("--public_name=%q", strings.TrimSpace(options.PublicName)))
	}
	if options.MaxDiskSize > 0 {
		command = append(command, fmt.Sprintf("--max_disk_size=%dGiB", options.MaxDiskSize))
	}
	if options.NumBlocks > 0 {
		command = append(command, fmt.Sprintf("--num_blocks=%d", options.NumBlocks))
	}

	// launch command
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx,
		"sh",
		"-c",
		strings.Join(command, " "),
	)

	// launch petals server in background and print output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		defer cancel()
		outchan <- []byte(fmt.Sprintf("Error: %s", err))
		return nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		defer cancel()
		outchan <- []byte(fmt.Sprintf("Error: %s", err))
		return nil, nil, err
	}

	// create a multireader to read from both stdout and stderr
	//scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	//go func() {
	//	for scanner.Scan() {
	//		outchan <- scanner.Bytes()
	//	}
	//}()

	// send output to the channel
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {
			outchan <- scanner.Bytes()
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {
			outchan <- scanner.Bytes()
		}
	}()

	go func() {
		<-ctx.Done()
		stdout.Close()
		stderr.Close()
		close(outchan)
	}()

	if err := cmd.Start(); err != nil {
		defer cancel()
		return nil, nil, err
	}

	return cancel, cmd, nil
}

func ForceKill() {
	cmd := exec.Command(
		"sh",
		"-c",
		`kill $(ps ax | grep '`+petalsServerProcessName+`' | grep -v "grep" | awk '{print $1}')`,
	)
	err := cmd.Run()
	if err != nil {
		log.Println("error killing petals server:", err)
	}
}
