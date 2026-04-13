package apperr

import "fmt"

func MessageError(action string, err error) error {
	return fmt.Errorf("%s: %w", action, err)
}