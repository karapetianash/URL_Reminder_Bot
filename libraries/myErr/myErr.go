package myErr

import (
	"fmt"
)

func Wrap(str string, err error) error {
	return fmt.Errorf("%s: %w", str, err)
}

func WrapIfErr(str string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(str, err)
}
