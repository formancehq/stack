package client

import (
	"fmt"

	"github.com/stripe/stripe-go/v79"
)

func handleError(itr *stripe.Iter) error {
	err := itr.Err()
	meta := itr.Meta()
	if meta == nil {
		return fmt.Errorf("failed stripe request: %w", err)
	}
	return fmt.Errorf("stripe api failed in req to %q: %w", meta.URL, err)
}
