package client

const genderizeURL = "https://api.genderize.io"

type GenderFetcher struct{}

func NewGenderFetcher() *GenderFetcher {
	return &GenderFetcher{}
}

// Fetch returns the response from https://api.genderize.io?name=name
func (*GenderFetcher) Fetch(name string) ([]byte, error) {
	return get(genderizeURL, name)
}
