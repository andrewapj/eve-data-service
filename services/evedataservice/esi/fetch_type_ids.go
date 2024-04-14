package esi

import (
	"context"
	"fmt"
)

// Ids represents a list of ids supplied by the ESI for various kinds of data.
type Ids struct {
	Ids []int
	baseEsiModel
}

// esiTypeIdsRequest represents a request that can be used to fetch the Type ids from the ESI.
func esiTypeIdsRequest() request {
	return newRequest("/latest/universe/types/", make(map[string]string))
}

// FetchTypeIds will retrieve Type ids from the ESI.
func (c *Client) FetchTypeIds(ctx context.Context, pageRequest PageRequest) (Ids, error) {

	resp, err := c.getPages(ctx, esiTypeIdsRequest(), pageRequest)
	if err != nil {
		return Ids{}, fmt.Errorf("unable to fetch type ids from ESI. %w", err)
	}

	return mapToIds(resp)
}
