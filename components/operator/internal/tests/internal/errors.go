package internal

import "fmt"

func mismatchTypeError(expected, got any) error {
	return fmt.Errorf("expected object of type %T, got %T", expected, got)
}
