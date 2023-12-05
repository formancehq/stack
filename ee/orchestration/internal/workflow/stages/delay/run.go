package delay

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func RunDelay(ctx workflow.Context, delay Delay) error {
	var duration time.Duration
	switch {
	case delay.Duration != nil:
		duration = time.Duration(*delay.Duration)
	case delay.Until != nil:
		duration = time.Until(*delay.Until)
	}
	return workflow.Sleep(ctx, duration)
}
