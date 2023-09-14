package apitools

import (
	"net/http"
	"strconv"
)

const (
	limitHeader     = "X-Rate-Limit-Limit"
	remainingHeader = "X-Rate-Limit-Remaining"
	resetHeader     = "X-Rate-Limit-Reset"
)

// RateLimit holds the rat limit information from request headers.
type RateLimit struct {
	// Count of names allotted for the current time window.
	Limit int64
	// Count of names left in the current time window.
	Remaining int64
	// Seconds until the rate limit window resets.
	Reset int64
}

// RateLimitFromHeaders returns a rate limit from the response headers.
func RateLimitFromHeaders(resp *http.Response) *RateLimit {
	limit, limitErr := strconv.ParseInt(resp.Header.Get(limitHeader), 10, 64)
	remaining, remainingErr := strconv.ParseInt(resp.Header.Get(remainingHeader), 10, 64)
	reset, resetErr := strconv.ParseInt(resp.Header.Get(resetHeader), 10, 64)

	if limitErr == nil && remainingErr == nil && resetErr == nil {
		return &RateLimit{
			Limit:     limit,
			Remaining: remaining,
			Reset:     reset,
		}
	}

	return nil
}
