package brokers

import (
	"fmt"
	"sort"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(ctx core.Context, stack *v1beta1.Stack, broker *v1beta1.Broker) error {

	brokerURI, err := settings.RequireURL(ctx, stack.Name, "broker", "dsn")
	if err != nil {
		return err
	}
	if brokerURI == nil {
		return errors.New("broker configuration not found")
	}

	if broker.Status.Mode == "" {
		if err := detectBrokerMode(ctx, stack, broker, brokerURI); err != nil {
			return err
		}
	}

	switch brokerURI.Scheme {
	case "nats":
		switch broker.Status.Mode {
		case v1beta1.ModeOneStreamByService:
			if err := createOneStreamByTopic(ctx, stack, broker, brokerURI); err != nil {
				return err
			}
		case v1beta1.ModeOneStreamByStack:
			if err := createOneStreamByStack(ctx, stack, broker, brokerURI); err != nil {
				return err
			}
		}
	}

	broker.Status.URI = brokerURI

	return nil
}

func detectBrokerMode(ctx core.Context, stack *v1beta1.Stack, broker *v1beta1.Broker, uri *v1beta1.URI) error {
	if hasLegacyStream, err := detectBrokerModeByCheckingExistentStreams(ctx, stack, broker, uri); err != nil {
		return err
	} else if hasLegacyStream {
		broker.Status.Mode = v1beta1.ModeOneStreamByService
		return nil
	}
	if ok, err := hasAllVersionsGreaterThan(ctx, stack, "v2.0.0-rc.27"); err != nil {
		return err
	} else if ok {
		broker.Status.Mode = v1beta1.ModeOneStreamByStack
	} else {
		broker.Status.Mode = v1beta1.ModeOneStreamByService
	}

	return nil
}

func hasAllVersionsGreaterThan(ctx core.Context, stack *v1beta1.Stack, ref string) (bool, error) {
	switch {
	case stack.Spec.Version != "":
		if !semver.IsValid(stack.Spec.Version) {
			return true, nil
		}
		return semver.Compare(stack.Spec.Version, ref) >= 0, nil
	case stack.Spec.VersionsFromFile != "":
		versions := &v1beta1.Versions{}
		if err := ctx.GetClient().Get(ctx, types.NamespacedName{
			Name: stack.Spec.VersionsFromFile,
		}, versions); err != nil {
			return false, err
		}
		for service, v := range versions.Spec {
			// notes(gfyrag): control is stick to version v1.7.0
			if service == "control" {
				continue
			}
			if !semver.IsValid(v) {
				continue
			}
			if semver.Compare(v, ref) < 0 {
				return false, nil
			}
		}

		return true, nil
	default:
		// notes(gfyrag): Using latest, be careful, we cannot sure `latest` image on the nodes are really latest
		return true, nil
	}
}

func detectBrokerModeByCheckingExistentStreams(ctx core.Context, stack *v1beta1.Stack, broker *v1beta1.Broker, uri *v1beta1.URI) (bool, error) {
	const script = `
	# notes(gfyrag): Check if we have any stream named "$STACK-xxx"
	v=$(nats stream ls -n --server $NATS_URI | grep "$STACK-")
	# exit with code 12 if we detect any streams
	[[ -z "$v" ]] || exit 12;
`

	natsBoxImage, err := registries.GetNatsBoxImage(ctx, stack, "0.14.1")
	if err != nil {
		return false, err
	}

	hasLegacyStream := false
	if err := jobs.Handle(ctx, broker, "detect-mode", corev1.Container{
		Image: natsBoxImage,
		Name:  "detect-mode",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", uri.Host)),
			core.Env("STACK", broker.Spec.Stack),
		},
	},
		// notes(gfyrag): As the time of writing these lines, the succeedPolicy feature (https://kubernetes.io/docs/concepts/workloads/controllers/job/#success-policy)
		// is too early to be used. Keep an eye on the k8s versions to switch when appropriate.
		jobs.WithPodFailurePolicy(batchv1.PodFailurePolicy{
			Rules: []batchv1.PodFailurePolicyRule{{
				Action: batchv1.PodFailurePolicyActionFailJob,
				OnExitCodes: &batchv1.PodFailurePolicyOnExitCodesRequirement{
					Operator: batchv1.PodFailurePolicyOnExitCodesOpIn,
					Values:   []int32{12},
				},
			}},
		}),
		jobs.WithValidator(func(job *batchv1.Job) bool {
			if job.Status.Succeeded > 0 {
				return true
			}
			// notes(gfyrag): podFailurePolicy mark the job as failed
			// so, we need to watch conditions to determine if it fails because of exit code 12
			for _, condition := range job.Status.Conditions {
				if condition.Type == "Failed" &&
					condition.Reason == "PodFailurePolicy" &&
					condition.Status == "True" {
					hasLegacyStream = true
					return true
				}
			}
			return false
		}),
	); err != nil {
		return false, err
	}

	return hasLegacyStream, nil
}

func deleteBroker(ctx core.Context, broker *v1beta1.Broker) error {
	if !broker.Status.Ready {
		return nil
	}
	if broker.Status.URI.Scheme != "nats" {
		return nil
	}

	script := ""
	switch broker.Status.Mode {
	case v1beta1.ModeOneStreamByService:
		script = `
			for stream in $(nats --server $NATS_URI stream ls -n | grep $STACK); do
				nats stream info --server $NATS_URI $stream && nats stream rm -f --server $NATS_URI $stream || true	
			done
		`
	case v1beta1.ModeOneStreamByStack:
		script = `
			nats stream info --server $NATS_URI $STACK && nats stream rm -f --server $NATS_URI $STACK || true
		`
	}

	stack := &v1beta1.Stack{}
	if err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Name: broker.Spec.Stack,
	}, stack); err != nil {
		return err
	}

	natsBoxImage, err := registries.GetNatsBoxImage(ctx, stack, "0.14.1")
	if err != nil {
		return err
	}

	return jobs.Handle(ctx, broker, "delete-streams", corev1.Container{
		Image: natsBoxImage,
		Name:  "delete-streams",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", broker.Status.URI.Host)),
			core.Env("STACK", broker.Spec.Stack),
		},
	})
}

func createOneStreamByStack(ctx core.Context, stack *v1beta1.Stack, broker *v1beta1.Broker, uri *v1beta1.URI) error {

	if broker.Status.Ready {
		return nil
	}

	const script = `
	index=$(nats --server $NATS_URI stream ls -n | grep "$STREAM")
	[[ -z "$index" ]] && {
		nats stream add \
			--server $NATS_URI \
			--retention interest \
			--subjects $STREAM.* \
			--defaults \
			--replicas $REPLICAS \
			--no-allow-direct \
			$STREAM
	} || true`

	natsBoxImage, err := registries.GetNatsBoxImage(ctx, stack, "0.14.1")
	if err != nil {
		return err
	}

	return jobs.Handle(ctx, broker, "create-stream", corev1.Container{
		Image: natsBoxImage,
		Name:  "create-topic",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", uri.Host)),
			core.Env("STREAM", stack.Name),
			core.Env("REPLICAS", func() string {
				if replicas := uri.Query().Get("replicas"); replicas != "" {
					return replicas
				}
				return "1"
			}()),
		},
	})
}

func createOneStreamByTopic(ctx core.Context, stack *v1beta1.Stack, broker *v1beta1.Broker, brokerURI *v1beta1.URI) error {
	l := &v1beta1.BrokerTopicList{}
	if err := ctx.GetClient().List(ctx, l, client.MatchingFields{
		"stack": stack.Name,
	}); err != nil {
		return err
	}

	for _, item := range l.Items {
		item := item
		if !collectionutils.Contains(broker.Status.Streams, item.Spec.Service) {
			if err := createNatsTopic(ctx, stack, broker, &item, brokerURI); err != nil {
				return err
			}
			broker.Status.Streams = append(broker.Status.Streams, item.Spec.Service)
		}
	}

	sort.Strings(broker.Status.Streams)

	return nil
}

func createNatsTopic(ctx core.Context, stack *v1beta1.Stack, broker *v1beta1.Broker, topic *v1beta1.BrokerTopic, brokerURI *v1beta1.URI) error {
	const script = `
	index=$(nats --server $NATS_URI stream ls -j | jq "index(\"$SUBJECT\")")
	if [ "$index" = "null" ]; then
		nats stream add \
			--server $NATS_URI \
			--retention interest \
			--subjects $SUBJECT \
			--defaults \
			--replicas $REPLICAS \
			--no-allow-direct \
			$STREAM
	fi`

	natsBoxImage, err := registries.GetNatsBoxImage(ctx, stack, "0.14.1")
	if err != nil {
		return err
	}

	return jobs.Handle(ctx, broker, "create-topic-"+topic.Spec.Service, corev1.Container{
		Image: natsBoxImage,
		Name:  "create-topic",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", brokerURI.Host)),
			core.Env("SUBJECT", fmt.Sprintf("%s-%s", stack.Name, topic.Spec.Service)),
			core.Env("STREAM", fmt.Sprintf("%s-%s", stack.Name, topic.Spec.Service)),
			core.Env("REPLICAS", func() string {
				if replicas := brokerURI.Query().Get("replicas"); replicas != "" {
					return replicas
				}
				return "1"
			}()),
		},
	})
}
