package petalsserver

import "testing"

func TestSMICountComputeProcesses(t *testing.T) {
	c := GetComputeProcessCount()
	t.Logf("Compute processes: %d", c)
}

func TestSMIGPUInfo(t *testing.T) {
	info, err := getGPUInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", info)
}
