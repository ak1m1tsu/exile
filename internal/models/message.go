package models

type ErrorMessage struct {
	Meta  map[string]any `json:"meta"`
	Error string         `json:"error"`
}
