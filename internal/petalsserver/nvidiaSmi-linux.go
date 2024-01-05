//go:build linux

package petalsserver

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"
)

var GPUStatus *SMILogs
var lock = &sync.Mutex{}

func init() {
	log.Println("init")
	var err error
	GPUStatus, err = getGPUInfo()
	if err != nil {
		log.Println("error getting gpu info:", err)
	}
	go func() {
		for {
			select {
			case <-time.Tick(1 * time.Second):
				lock.Lock()
				var err error
				GPUStatus, err = getGPUInfo()
				if err != nil {
					log.Println("error getting gpu info:", err)
				}
				lock.Unlock()
			}
		}
	}()

}

const petalsServerProcessName = "petals.cli.run_server"

type Utilization struct {
	GPUUtil    string `xml:"gpu_util"`
	MemoryUtil string `xml:"memory_util"`
}

type ProcessInfo struct {
	Type string `xml:"type"`
	Name string `xml:"process_name"`
}

type MemoryInfo struct {
	Total string `xml:"total"`
	Used  string `xml:"used"`
	Free  string `xml:"free"`
}

// GPU is the representeation of the GPU in the nvidia-smi -q -x (XML) command.
type GPU struct {
	ID            string        `xml:"id,attr"`
	Utilization   Utilization   `xml:"utilization"`
	Processes     []ProcessInfo `xml:"processes>process_info"`
	FBMemoryUsage MemoryInfo    `xml:"fb_memory_usage"`
}

// SMILogs is the representeation of the logs of the nvidia-smi -q -x (XML) command.
type SMILogs struct {
	XMLName       string `xml:"nvidia_smi_log"`
	CUDAVersion   string `xml:"cuda_version"`
	DriverVersion string `xml:"driver_version"`
	GPU           []GPU  `xml:"gpu"`
}

func getGPUInfo() (*SMILogs, error) {
	cmd := exec.Command(
		"nvidia-smi",
		"-q",
		"-x",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	defer stdout.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	output := buf.String()
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	smiLogs := &SMILogs{}
	err = xml.Unmarshal([]byte(output), smiLogs)
	if err != nil {
		return nil, err
	}
	return smiLogs, nil
}

func GetComputeProcessCount() int {
	lock.Lock()
	defer lock.Unlock()
	count := 0
	for _, gpu := range GPUStatus.GPU {
		for _, process := range gpu.Processes {
			if process.Type == "C" {
				count++
			}
		}
	}
	return count
}

func GetFreeMemory() float64 {
	lock.Lock()
	defer lock.Unlock()
	total := 0.0
	used := 0.0

	for _, gpu := range GPUStatus.GPU {
		total += float64(sizeToInt(gpu.FBMemoryUsage.Total))
		used += float64(sizeToInt(gpu.FBMemoryUsage.Used))
	}
	return 100.0 * (used / total)
}

func sizeToInt(size string) int {
	var sizeInt int
	var unit string
	fmt.Sscanf(size, "%d %s", &sizeInt, &unit)
	switch unit {
	case "KiB":
		sizeInt *= 1024
	case "MiB":
		sizeInt *= 1024 * 1024
	case "GiB":
		sizeInt *= 1024 * 1024 * 1024
	case "TiB":
		sizeInt *= 1024 * 1024 * 1024 * 1024
	}
	return sizeInt
}
