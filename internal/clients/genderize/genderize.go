package genderize

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/apitools"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/errtools"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/validator"
)

const apiURL = "https://api.genderize.io"

// Response represents the response from the Genderize API.
type Response struct {
	Count       int     `json:"count"`
	Name        string  `json:"string"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

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

// Query represents the URL parameters for the Genderize API.
type Query struct {
	Name      string `validate:"required"`
	CountryID string `validate:"omitempty,len=2"`
}

// Get returns the gender of given name and nationality.
//
// If status code of the response is not 2xx, it returns nil, APIError.
func Get(query Query) (*Response, error) {
	const op = "clients.genderize.Get"

	if err := validator.ValidateStruct(query); err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	params := url.Values{}
	params.Add("name", query.Name)
	if query.CountryID != "" {
		params.Add("country_id", query.CountryID)
	}

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

	success := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	decoder := json.NewDecoder(resp.Body)

	if !success {
		apiErr := APIError{StatusCode: resp.StatusCode}

		if err = decoder.Decode(&apiErr); err != nil {
			return nil, errtools.WithOperation(err, op)
		}

		if rt := apitools.RateLimitFromHeaders(resp); rt != nil {
			apiErr.RateLimit = rt
		}

		return nil, apiErr
	}

	genderizeResp := &Response{}
	err = decoder.Decode(&genderizeResp)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	return genderizeResp, nil
}
