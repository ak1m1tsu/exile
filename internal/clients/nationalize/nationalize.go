package nationalize

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/apitools"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/errtools"
)

const apiURL = "https://api.nationalize.io"

// Response represents the response from the Nationalize API
type Response struct {
	Count   int       `json:"count"`
	Name    string    `json:"name"`
	Country []country `json:"country"`
}

type country struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

// APIError represents an error returned from the Nationalize API.
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

// Get returns the nationality of the given name.
//
// If the name is empty, it returns nil, nil.
// If status code of the response is not 2xx, it returns nil, APIError.
func Get(name string) (*Response, error) {
	const op = "clients.nationalize.Get"

	if name == "" {
		return nil, nil
	}

	params := url.Values{}
	params.Add("name", name)

	endpointURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}
	endpointURL.RawQuery = params.Encode()

	req := &http.Request{
		Method: http.MethodGet,
		URL:    endpointURL,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	decoder := json.NewDecoder(resp.Body)

	if !success {
		apiErr := APIError{
			StatusCode: resp.StatusCode,
		}

		if err = decoder.Decode(&apiErr); err != nil {
			return nil, errtools.WithOperation(err, op)
		}

		if rt := apitools.RateLimitFromHeaders(resp); rt != nil {
			apiErr.RateLimit = rt
		}

		return nil, apiErr
	}

	nationalizeResp := &Response{}
	err = decoder.Decode(&nationalizeResp)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	return nationalizeResp, nil
}
