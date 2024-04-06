package esi

import (
	"github.com/andrewapj/arcturus/clock"
	"time"
)

var esiRequestTime = clock.GetTime().Add(5 * time.Minute).Truncate(time.Second)

const (
	esiStatusResponse = `
{
  "players": 25953,
  "server_version": "2539399",
  "start_time": "2024-04-05T11:03:48Z"
}`
)
