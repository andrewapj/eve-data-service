package esi

import (
	"context"
	"fmt"
)

type Ids struct {
	Ids []int
}

func esiTypeIdsRequest() request {
	return newRequest("/latest/universe/types/", make(map[string]string))
}

func (c *Client) FetchTypeIds(ctx context.Context, pageRequest PageRequest) (Ids, error) {

	resp, err := c.getPages(ctx, esiTypeIdsRequest(), pageRequest)
	if err != nil {
		return Ids{}, fmt.Errorf("unable to fetch type ids from ESI. %w", err)
	}

	return mapToIds(resp)
}
