package models

type ErrorMessage struct {
	Meta  []byte `json:"meta"`
	Error string `json:"error"`
}
