package petalsserver

import (
	"testing"
	"time"
)

func TestLaunchPetalsServer(t *testing.T) {
	outchan := make(chan []byte)
	cancel, err := LaunchPetalsServer(nil, outchan)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Launched petals server")

	// wait for the server to start
	time.Sleep(5 * time.Second)
	cancel()
	t.Log("Stopped petals server")

	// wait for the server to stop
	time.Sleep(5 * time.Second)

	// check that the server is not running
	if IsPetalsServerRunning() {
		ForceKill()
		t.Fatal("Server is still running")
	}
}
