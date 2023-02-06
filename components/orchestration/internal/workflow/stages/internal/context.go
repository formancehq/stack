package internal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func InfiniteRetryContext(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumInterval:    100 * time.Second,
			NonRetryableErrorTypes: []string{
				//"INTERNAL",
				//"INSUFFICIENT_FUND",
				"VALIDATION",
				"CONFLICT",
				"NO_SCRIPT",
				"COMPILATION_FAILED",
				//"METADATA_OVERRIDE",
			},
		},
	})
}
