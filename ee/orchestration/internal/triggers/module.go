package triggers

import (
	"net/http"
	"strings"

	"github.com/formancehq/orchestration/internal/temporalworker"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/uptrace/bun"
	"go.temporal.io/sdk/client"
	"go.uber.org/fx"
)

func NewModule(taskQueue string) fx.Option {
	return fx.Options(
		fx.Provide(NewManager),
		fx.Provide(func(httpClient *http.Client) *expressionEvaluator {
			return NewExpressionEvaluator(httpClient)
		}),
		fx.Provide(func() *triggerWorkflow {
			return NewWorkflow(taskQueue, true)
		}),
		fx.Provide(fx.Annotate(func(workflow *triggerWorkflow) temporalworker.DefinitionSet {
			return workflow.DefinitionSet()
		}, fx.ResultTags(`group:"workflows"`))),
		fx.Provide(func(db *bun.DB, manager *workflow.WorkflowManager,
			expressionEvaluator *expressionEvaluator, publisher message.Publisher) Activities {
			return NewActivities(db, manager, expressionEvaluator, publisher)
		}),
		fx.Provide(fx.Annotate(func(activities Activities) temporalworker.DefinitionSet {
			return activities.DefinitionSet()
		}, fx.ResultTags(`group:"activities"`))),
	)
}

func NewListenerModule(taskIDPrefix, taskQueue string, topics []string) fx.Option {
	return fx.Options(
		fx.Invoke(func(logger logging.Logger, r *message.Router, s message.Subscriber, temporalClient client.Client) {
			logger.Infof("Listening events from topics: %s", strings.Join(topics, ","))
			registerListener(r, s, temporalClient, taskIDPrefix, taskQueue, topics)
		}),
	)
}
