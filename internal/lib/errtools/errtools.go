package errtools

import "fmt"

// WithOperation retuns new error with "op: err" format
func WithOperation(err error, op string) error {
	return fmt.Errorf("%s: %w", op, err)
}
