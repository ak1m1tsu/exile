package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/apitools"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Fetcher --output ./mocks --outpkg mocks
type Fetcher interface {
	Fetch(name string) ([]byte, error)
}

var (
	ErrNameEmpty       = errors.New("the name param is empty")
	ErrFindNationality = errors.New("could not find nationality for the name")
	ErrFindGender      = errors.New("could not find gender for the name")
	ErrFindAge         = errors.New("could not find age for the name")
)

// APIError represents an error returned from the Genderize API.
type APIError struct {
	// Error from API request
	Message string `json:"error"`
	// HTTP status code from API request
	StatusCode int
	// Rate limit information from API response headers
	RateLimit *apitools.RateLimit
}

// Error returns the error message from the APIError.
func (e APIError) Error() string {
	return e.Message
}

func get(apiURL string, name string) ([]byte, error) {
	if name == "" {
		return nil, ErrNameEmpty
	}

	params := url.Values{
		"name": []string{name},
	}

	endpoint, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	endpoint.RawQuery = params.Encode()

	req := &http.Request{
		Method: http.MethodGet,
		URL:    endpoint,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	decoder := json.NewDecoder(resp.Body)

	if !success {
		apiErr := APIError{
			StatusCode: resp.StatusCode,
		}

		if err = decoder.Decode(&apiErr); err != nil {
			return nil, err
		}

		if rt := apitools.RateLimitFromHeaders(resp); rt != nil {
			apiErr.RateLimit = rt
		}

		return nil, apiErr
	}

	return io.ReadAll(resp.Body)
}
