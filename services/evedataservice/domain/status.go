package domain

import (
	"encoding/json"
	"fmt"
	"github.com/andrewapj/arcturus/esi"
	"time"
)

type Status struct {
	Players       int       `json:"players"`
	ServerVersion string    `json:"server_version"`
	StartTime     time.Time `json:"start_time"`
	VIP           bool      `json:"vip"`
}

func (s Status) Id() int {
	return 1
}

func (s Status) Data() (*[]byte, error) {
	bytes, err := json.Marshal(s)
	return &bytes, err
}

func MapStatusFromEsi(status *esi.Status) (*Status, error) {

	if status == nil {
		return nil, fmt.Errorf("unable to map status, received a nil esi status")
	}

	return &Status{
		Players:       status.Players,
		ServerVersion: status.ServerVersion,
		StartTime:     status.StartTime,
		VIP:           status.VIP,
	}, nil
}
