package esi

import (
	"encoding/json"
	"fmt"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"log/slog"
	"time"
)

// mapToIds will convert a slice of responses to an Ids which contains all the integer ids.
func mapToIds(resp []*response) (Ids, error) {

	var response Ids
	var expiries []time.Time

	if len(resp) == 0 {
		return Ids{}, fmt.Errorf("mapToIds: expected at least 1 response, got %d", len(resp))
	}

	for i := range resp {
		var ids []int
		err := json.Unmarshal(resp[i].body, &ids)
		if err != nil {
			return Ids{}, fmt.Errorf("error decoding from json %w", err)
		}

		response.Ids = append(response.Ids, ids...)
		expiries = append(expiries, resp[i].expires)
	}

	setMetadata(&response, resp)
	return response, nil
}

// mapToSingle will convert a slice of responses to the target ESI type. It expects only one response.
func mapToSingle[T BaseEsiModel](r []*response) (T, error) {

	var t T

	if len(r) != 1 {
		return t, fmt.Errorf("mapToSingle: expected 1 response, got %d", len(r))
	}

	err := json.Unmarshal(r[0].body, &t)
	if err != nil {
		return t, fmt.Errorf("error decoding from json %w", err)
	}

	setMetadata(t, r)
	return t, err
}

// setMetadata will set the metadata for ESI types.
func setMetadata(b BaseEsiModel, r []*response) {

	var times []time.Time
	for i := range r {
		times = append(times, r[i].expires)
	}
	expiry, err := clock.FindEarliestTime(times)
	if err != nil {
		expiry = clock.GetTime().Add(time.Duration(config.EsiDateAdditionalTime()) * time.Second)
		slog.Error("error finding earliest time when setting esi metadata, set to future default", "err", err, "expiry", expiry)
	}

	b.SetPages(r[0].pages)
	b.SetExpires(expiry)
}
