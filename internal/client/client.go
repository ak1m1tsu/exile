package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/apitools"
)

// TODO: wrap all fetch functions in structs that implement Fetcher interface
const (
	genderizeURL   = "https://api.genderize.io"
	nationalizeURL = "https://api.nationalize.io"
	agifyURL       = "https://api.agify.io"
)

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

type agifyResponse struct {
	Age int `json:"age"`
}

// FetchAge returns the age of given name from Agify API.
func FetchAge(name string) (int, error) {
	body, err := get(agifyURL, name)
	if err != nil {
		return 0, err
	}

	var resp agifyResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}

	if resp.Age == 0 {
		return 0, ErrFindAge
	}

	return resp.Age, nil
}

type genderizeResponse struct {
	Gender string `json:"gender"`
}

// FetchGender returns the gender of given name from Genderize API.
func FetchGender(name string) (string, error) {
	body, err := get(genderizeURL, name)
	if err != nil {
		return "", err
	}

	var resp genderizeResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	if resp.Gender == "" {
		return "", ErrFindGender
	}

	return resp.Gender, nil
}

type nationalizeResponse struct {
	Country []country `json:"country"`
}

type country struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

// FetchNationality returns the nationality of the given name from Nationalize API.
func FetchNationality(name string) (string, error) {
	body, err := get(nationalizeURL, name)
	if err != nil {
		return "", err
	}

	var resp nationalizeResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	if len(resp.Country) == 0 {
		return "", ErrFindNationality
	}

	return resp.Country[0].ID, nil
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
