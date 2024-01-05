package petalsserver

import "testing"

func TestGetStatus(t *testing.T) {
	status, err := GetStatus()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", status)
}
