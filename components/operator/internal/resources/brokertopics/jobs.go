package brokertopics

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	v1 "k8s.io/api/batch/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createJob(ctx core.Context, topic *v1beta1.BrokerTopic, configuration v1beta1.BrokerConfiguration) (*v1.Job, error) {

	job, _, err := core.CreateOrUpdate[*v1.Job](ctx, types.NamespacedName{
		Namespace: topic.Spec.Stack,
		Name:      fmt.Sprintf("%s-create-topic", topic.Spec.Service),
	},
		func(t *v1.Job) {
			args := []string{"nats", "stream", "add",
				"--server", fmt.Sprintf("nats://%s", configuration.Nats.URL),
				"--retention", "interest",
				"--subjects", topic.Name,
				"--defaults",
			}
			if configuration.Nats.Replicas > 0 {
				args = append(args, "--replicas", fmt.Sprint(configuration.Nats.Replicas))
			}
			args = append(args, topic.Name)

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = v12.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []v12.Container{{
				Image: "natsio/nats-box:0.14.1",
				Name:  "create-topic",
				Args:  args,
			}}
		},
		core.WithController[*v1.Job](ctx.GetScheme(), topic),
	)
	return job, err
}

func deleteJob(ctx core.Context, topic *v1beta1.BrokerTopic) (*v1.Job, error) {
	job, _, err := core.CreateOrUpdate[*v1.Job](ctx, types.NamespacedName{
		Namespace: topic.Spec.Stack,
		Name:      fmt.Sprintf("%s-delete-topic", topic.Spec.Service),
	},
		func(t *v1.Job) {
			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = v12.RestartPolicyOnFailure
			t.Spec.Template.Spec.Containers = []v12.Container{{
				Image: "natsio/nats-box:0.14.1",
				Name:  "create-topic",
				Args: []string{"nats", "stream", "rm", "-f", "--server",
					fmt.Sprintf("nats://%s", topic.Status.Configuration.Nats.URL), topic.Name},
			}}
		},
		core.WithController[*v1.Job](ctx.GetScheme(), topic),
	)
	return job, err
}
