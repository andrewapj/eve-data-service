package esi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func mapToIds(resp []*response) (Ids, error) {

	var response Ids

	buff := new(bytes.Buffer)
	for i := range resp {
		_, err := buff.Write(resp[i].body)
		if err != nil {
			return Ids{}, fmt.Errorf("error writing to buffer %w", err)
		}

		var ids []int
		err = json.NewDecoder(buff).Decode(&ids)
		if err != nil {
			return Ids{}, fmt.Errorf("error decoding from json %w", err)
		}
		response.Ids = append(response.Ids, ids...)
		buff.Reset()
	}

	return response, nil
}

func mapToSingle[T BaseEsiModel](r []*response) (T, error) {

	var t T

	if len(r) != 1 {
		return t, fmt.Errorf("mapToSingle: expected 1 result, got %d", len(r))
	}

	buff := new(bytes.Buffer)
	_, err := buff.Write(r[0].body)
	if err != nil {
		return t, fmt.Errorf("error writing to buffer %w", err)
	}

	err = json.NewDecoder(buff).Decode(&t)
	if err != nil {
		return t, fmt.Errorf("error decoding from json %w", err)
	}

	setMetadata(t, r[0])
	return t, err
}

func setMetadata(b BaseEsiModel, r *response) {

	b.SetPages(r.pages)
	b.SetExpires(r.expires)
}
