package petalsserver

import "testing"

func TestInstallPetals(t *testing.T) {
	c, err := PipInstallPetals("python3.11")
	if err != nil {
		t.Fatal(err)
	}
	for line := range c {
		t.Log(line)
	}
	t.Log("Installed petals")
}

func TestIsInstalled(t *testing.T) {
	state := IsPetalsServerInstalled()
	t.Logf("Installed: %v", state)
}
