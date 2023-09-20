package client

import (
	"encoding/json"
	"fmt"
)

const genderizeURL = "https://api.genderize.io"

type GenderFetcher struct{}

func NewGenderFetcher() *GenderFetcher {
	return &GenderFetcher{}
}

// Fetch returns the response from https://api.genderize.io?name=name
func (*GenderFetcher) Fetch(name string) ([]byte, error) {
	data, err := get(genderizeURL, name)
	if err != nil {
		return nil, err
	}

	var resp genderizeResponse
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf(`{"gender": "%s"}`, resp.Gender)), nil
}

type genderizeResponse struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
}
