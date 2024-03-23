package wrappererr

import "fmt"

func Wrap(customErr, err error) error {
	return fmt.Errorf("%w: %w", customErr, err)
}
