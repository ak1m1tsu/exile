package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type NationalizeResponse struct {
	Count   int           `json:"count"`
	Country []nationality `json:"country"`
	Name    string        `json:"name"`
}

type nationality struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type NationalityFetcher struct {
	endpoint string
}

func NewNationalityFetcher(endpoint string) *NationalityFetcher {
	return &NationalityFetcher{endpoint}
}

func (nf *NationalityFetcher) Fetch(ctx context.Context, name string) ([]byte, error) {
	data, err := get(ctx, nf.endpoint, name)
	if err != nil {
		return nil, fmt.Errorf("error fetching nationality: %w", err)
	}

	resp := new(NationalizeResponse)
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling nationality: %w", err)
	}

	nationality := map[string]any{
		"nationality": resp.Country[0].ID,
		"name":        resp.Name,
	}

	data, err = json.Marshal(nationality)
	if err != nil {
		return nil, fmt.Errorf("error marshalling nationality: %w", err)
	}

	return data, nil
}
