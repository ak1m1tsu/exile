package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var ErrEmptyName = errors.New("the name is empty")

func get(ctx context.Context, endpoint, name string) ([]byte, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	params := url.Values{
		"name": {name},
	}

	uri, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	uri.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer res.Body.Close()

	success := res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices
	if !success {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}
