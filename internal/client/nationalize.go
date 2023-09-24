package client

import (
	"encoding/json"
	"fmt"
)

const nationalizeURL = "https://api.nationalize.io"

type NationalityFetcher struct{}

func NewNationalityFetcher() *NationalityFetcher {
	return &NationalityFetcher{}
}

// Fetch retuns the response from https://api.nationalize.io?name=name
func (*NationalityFetcher) Fetch(name string) ([]byte, error) {
	data, err := get(nationalizeURL, name)
	if err != nil {
		return nil, err
	}

	resp := nationalizeResponse{Country: make([]nationality, 0)}
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf(`{"nationality":"%s"}`, resp.Country[0].ID)), nil
}

type nationalizeResponse struct {
	Count   int           `json:"count"`
	Country []nationality `json:"country"`
	Name    string        `json:"name"`
}

type nationality struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
