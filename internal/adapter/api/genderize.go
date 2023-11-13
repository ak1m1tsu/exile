package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type GenderResponse struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

type GenderFetcher struct {
	endpoint string
}

func NewGenderFetcher(endpoint string) *GenderFetcher {
	return &GenderFetcher{endpoint: endpoint}
}

func (gf *GenderFetcher) Fetch(ctx context.Context, name string) ([]byte, error) {
	data, err := get(ctx, gf.endpoint, name)
	if err != nil {
		return nil, fmt.Errorf("error fetching gender: %w", err)
	}

	resp := new(GenderResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("error decoding gender: %w", err)
	}

	data, err = json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error encoding gender: %w", err)
	}

	return data, nil
}
