package esi

import (
	"context"
	"fmt"
	"time"
)

// Status represents a status received from the ESI.
type Status struct {
	Players       uint32    `json:"players"`
	ServerVersion string    `json:"server_version"`
	StartTime     time.Time `json:"start_time"`
	VIP           bool      `json:"vip"`
	baseEsiModel
}

// esiStatusRequest represents a request for an ESI status.
func esiStatusRequest() request {
	return newRequest("/latest/status/", map[string]string{})
}

// FetchStatus will retrieve a Status from the ESI.
func (c *Client) FetchStatus(ctx context.Context) (Status, error) {

	resp, err := c.getPages(ctx, esiStatusRequest(), NewPageRequest(0, 1))
	if err != nil {
		return Status{}, fmt.Errorf("unable to fetch status from ESI. %w", err)
	}

	status, err := mapToSingle[*Status](resp)
	if err != nil {
		return Status{}, fmt.Errorf("unable to parse status from ESI. %w", err)
	}

	return *status, nil
}
