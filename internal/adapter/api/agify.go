package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type AgeResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type AgeFetcher struct {
	endpoint string
}

func NewAgeFetcher(endpoint string) AgeFetcher {
	return AgeFetcher{endpoint: endpoint}
}

func (a AgeFetcher) Fetch(ctx context.Context, name string) ([]byte, error) {
	data, err := get(ctx, a.endpoint, name)
	if err != nil {
		return nil, fmt.Errorf("error fetching age for %s: %w", name, err)
	}

	resp := new(AgeResponse)
	if err = json.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling age response: %w", err)
	}

	data, err = json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshalling age response: %w", err)
	}

	return data, nil
}
