package brokers

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Reconcile(ctx core.Context, stack *v1beta1.Stack, broker *v1beta1.Broker) error {

	brokerURI, err := settings.RequireURL(ctx, stack.Name, "broker.dsn")
	if err != nil {
		return err
	}
	if brokerURI == nil {
		return errors.New("broker configuration not found")
	}

	switch brokerURI.Scheme {
	case "nats":
		if err := createOneStreamByTopic(ctx, stack, broker, brokerURI); err != nil {
			return err
		}
	}

	broker.Status.URI = brokerURI

	return nil
}

func deleteBroker(ctx core.Context, broker *v1beta1.Broker) error {
	if !broker.Status.Ready {
		return nil
	}
	if broker.Status.URI.Scheme != "nats" {
		return nil
	}

	const script = `
	set -xe
	for stream in $(nats --server $NATS_URI stream ls -n | grep $STACK); do
		nats stream rm -f --server $NATS_URI $stream	
	done
`
	return jobs.Handle(ctx, broker, "delete-streams", corev1.Container{
		Image: "natsio/nats-box:0.14.1",
		Name:  "delete-streams",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", broker.Status.URI.Host)),
			core.Env("STACK", broker.Spec.Stack),
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

	grp, _ := errgroup.WithContext(ctx)

	for _, item := range l.Items {
		item := item
		grp.Go(func() error {
			return createNatsTopic(ctx, stack, broker, &item, brokerURI)
		})
	}

	return grp.Wait()
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
