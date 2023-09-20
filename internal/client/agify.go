package client

const agifyURL = "https://api.agify.io"

type AgeFetcher struct{}

func NewAgeFetcher() *AgeFetcher {
	return &AgeFetcher{}
}

// Fetch returns the response from https://api.agify.io?name=name
func (*AgeFetcher) Fetch(name string) ([]byte, error) {
	return get(agifyURL, name)
}
