package client

import (
	"encoding/json"
	"fmt"
)

const agifyURL = "https://api.agify.io"

type AgeFetcher struct{}

func NewAgeFetcher() *AgeFetcher {
	return &AgeFetcher{}
}

// Fetch returns the response from https://api.agify.io?name=name
func (*AgeFetcher) Fetch(name string) ([]byte, error) {
	data, err := get(agifyURL, name)
	if err != nil {
		return nil, err
	}

	var resp agifyResponse
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf(`{"age":%d}`, resp.Age)), nil
}

type agifyResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}
