package petalsserver

import (
	"encoding/json"
	"net/http"
)

const (
	HealthPage = "https://health.petals.dev/api/v1/state"
)

type Peer struct {
	PeerID string `json:"peer_id"`
	State  string `json:"state"`
	Span   struct {
		ServerInfo struct {
			PublicName string `json:"public_name"`
			State      string `json:"state"`
			StartBlock int    `json:"start_block"`
			EndBlock   int    `json:"end_block"`
			UsingRelay bool   `json:"using_relay"`
			Version    string `json:"version"`
		} `json:"server_info"`
	} `json:"span"`
}

type Status struct {
	ModelReports []struct {
		Name       string `json:"name"`
		State      string `json:"state"`
		NumBlocks  int    `json:"num_blocks"`
		ServerRows []Peer `json:"server_rows"`
	} `json:"model_reports"`
}

func GetStatus() (*Status, error) {
	resp, err := http.Get(HealthPage)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	status := &Status{}
	err = json.NewDecoder(resp.Body).Decode(status)
	if err != nil {
		return nil, err
	}
	return status, nil
}
