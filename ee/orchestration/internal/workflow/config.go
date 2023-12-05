package workflow

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/formancehq/orchestration/internal/schema"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type RawStage map[string]map[string]any

type Config struct {
	Name   string     `json:"name"`
	Stages []RawStage `json:"stages"`
}

func (c *Config) runStage(ctx workflow.Context, s Stage, stage RawStage, variables map[string]string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
		}
	}()
	var (
		name  string
		value map[string]any
	)
	for name, value = range stage {
	}

	stageSchema, err := schema.Resolve(schema.Context{
		Variables: variables,
	}, value, name)
	if err != nil {
		return err
	}

	if err := schema.ValidateRequirements(stageSchema); err != nil {
		return err
	}

	err = workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowID: s.TemporalWorkflowID(),
		}),
		stageSchema.GetWorkflow(),
		stageSchema,
	).Get(ctx, nil)
	if err != nil {
		var appError *temporal.ApplicationError
		if errors.As(err, &appError) {
			return errors.New(appError.Message())
		}
		var canceledError *temporal.CanceledError
		if errors.As(err, &canceledError) {
			return canceledError
		}
		return err
	}

	return nil
}

func (c *Config) run(ctx workflow.Context, db *bun.DB, instance Instance, variables map[string]string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.WithStack(fmt.Errorf("%s", e))
			debug.PrintStack()
		}
	}()

	logger := workflow.GetLogger(ctx)
	for ind, rawStage := range c.Stages {
		logger.Info("run stage", "index", ind, "workflowID", instance.ID)

		stage := NewStage(instance.ID, ind)

		if _, dbErr := db.NewInsert().
			Model(&stage).
			Exec(context.Background()); dbErr != nil {
			logger.Error("error inserting stage into database", "error", dbErr)
		}

		err := c.runStage(ctx, stage, rawStage, variables)
		stage.SetTerminated(err, workflow.Now(ctx).Round(time.Nanosecond))
		if err != nil {
			logger.Debug("error running stage", "error", stage.Error)
		}

		if _, dbErr := db.NewUpdate().
			Model(&stage).
			WherePK().
			Exec(context.Background()); dbErr != nil {
			logger.Error("error updating stage into database", "error", dbErr)
		}
		if err != nil {
			return err
		}
		logger.Info("stage terminated", "index", ind, "workflowID", instance.ID)
	}

	return nil
}

func (c *Config) Validate() error {
	for _, rawStage := range c.Stages {
		if len(rawStage) == 0 {
			return fmt.Errorf("empty specification")
		}
		if len(rawStage) > 1 {
			return fmt.Errorf("a specification should have only one name")
		}
		var (
			name  string
			value map[string]any
		)
		for name, value = range rawStage {
		}

		_, err := schema.Resolve(schema.Context{}, value, name)
		if err != nil {
			return err
		}
	}
	return nil
}
