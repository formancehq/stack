package brokertopics

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func checkExistingJob(ctx core.Context, topic *v1beta1.BrokerTopic, brokerURI *v1beta1.URI) (bool, error) {

	const script = `
	index=$(nats --server %s stream ls -j | jq 'index("%s")')
	if [ "$index" = "null" ]; then
		exit 28
	fi
`
	job, _, err := core.CreateOrUpdate[*v1.Job](ctx, types.NamespacedName{
		Namespace: topic.Spec.Stack,
		Name:      fmt.Sprintf("%s-check-topic", topic.Spec.Service),
	},
		func(t *v1.Job) error {
			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyNever
			if len(t.Spec.Template.Spec.Containers) == 0 {
				t.Spec.Template.Spec.Containers = []corev1.Container{{}}
			}
			t.Spec.Template.Spec.Containers[0].Image = "natsio/nats-box:0.14.1"
			t.Spec.Template.Spec.Containers[0].Name = "check-topic-existence"
			t.Spec.Template.Spec.Containers[0].Args = core.ShellScript(script, fmt.Sprintf("nats://%s", brokerURI.Host), topic.Name)
			t.Spec.PodFailurePolicy = &v1.PodFailurePolicy{
				Rules: []v1.PodFailurePolicyRule{
					{
						Action: "FailJob",
						OnExitCodes: &v1.PodFailurePolicyOnExitCodesRequirement{
							Operator: "In",
							Values:   []int32{28},
						},
					},
					{
						Action: "Ignore",
						OnPodConditions: []v1.PodFailurePolicyOnPodConditionsPattern{{
							Type: "DisruptionTarget",
						}},
					},
				},
			}

			return nil
		},
		core.WithController[*v1.Job](ctx.GetScheme(), topic),
	)
	if err != nil {
		return false, err
	}

	switch {
	case job.Status.Succeeded == 1:
		return true, nil
	case job.Status.Failed > 0:
		for _, condition := range job.Status.Conditions {
			if condition.Type != "Failed" {
				continue
			}
			if condition.Reason != "PodFailurePolicy" {
				continue
			}
			if condition.Status != "True" {
				continue
			}

			// notes(gfyrag): the container exited with status code 28, indicating than the topic does not exist
			return false, nil
		}
	}

	return false, core.NewPendingError()
}

func createJob(ctx core.Context, topic *v1beta1.BrokerTopic, brokerURI *v1beta1.URI) error {

	exists, err := checkExistingJob(ctx, topic, brokerURI)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	job, _, err := core.CreateOrUpdate[*v1.Job](ctx, types.NamespacedName{
		Namespace: topic.Spec.Stack,
		Name:      fmt.Sprintf("%s-create-topic", topic.Spec.Service),
	},
		func(t *v1.Job) error {
			args := []string{"nats", "stream", "add",
				"--server", fmt.Sprintf("nats://%s", brokerURI.Host),
				"--retention", "interest",
				"--subjects", topic.Name,
				"--defaults",
				"--no-allow-direct",
			}
			if replicas := brokerURI.Query().Get("replicas"); replicas != "" {
				args = append(args, "--replicas", replicas)
			}
			args = append(args, topic.Name)

			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyOnFailure
			if len(t.Spec.Template.Spec.Containers) == 0 {
				t.Spec.Template.Spec.Containers = []corev1.Container{{}}
			}
			t.Spec.Template.Spec.Containers[0].Image = "natsio/nats-box:0.14.1"
			t.Spec.Template.Spec.Containers[0].Name = "create-topic"
			t.Spec.Template.Spec.Containers[0].Args = args

			return nil
		},
		core.WithController[*v1.Job](ctx.GetScheme(), topic),
	)
	if job.Status.Succeeded == 0 {
		return core.NewPendingError()
	}
	return nil
}

func deleteJob(ctx core.Context, topic *v1beta1.BrokerTopic) (*v1.Job, error) {
	job, _, err := core.CreateOrUpdate[*v1.Job](ctx, types.NamespacedName{
		Namespace: topic.Spec.Stack,
		Name:      fmt.Sprintf("%s-delete-topic", topic.Spec.Service),
	},
		func(t *v1.Job) error {
			t.Spec.BackoffLimit = pointer.For(int32(10000))
			t.Spec.TTLSecondsAfterFinished = pointer.For(int32(30))
			t.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyOnFailure
			if len(t.Spec.Template.Spec.Containers) == 0 {
				t.Spec.Template.Spec.Containers = []corev1.Container{{}}
			}
			t.Spec.Template.Spec.Containers[0].Image = "natsio/nats-box:0.14.1"
			t.Spec.Template.Spec.Containers[0].Name = "delete-topic"
			t.Spec.Template.Spec.Containers[0].Args = []string{"nats", "stream", "rm", "-f", "--server",
				fmt.Sprintf("nats://%s", topic.Status.URI.Host), topic.Name}

			return nil
		},
		core.WithController[*v1.Job](ctx.GetScheme(), topic),
	)
	return job, err
}
