package client

const nationalizeURL = "https://api.nationalize.io"

type NationalityFetcher struct{}

func NewNationalityFetcher() *NationalityFetcher {
	return &NationalityFetcher{}
}

// Fetch retuns the response from https://api.nationalize.io?name=name
func (*NationalityFetcher) Fetch(name string) ([]byte, error) {
	return get(nationalizeURL, name)
}
