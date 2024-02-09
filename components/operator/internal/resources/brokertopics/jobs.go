package brokertopics

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/jobs"
	corev1 "k8s.io/api/core/v1"
)

func createJob(ctx core.Context, topic *v1beta1.BrokerTopic, brokerURI *v1beta1.URI) error {
	const script = `
	set -xe
	index=$(nats --server $NATS_URI stream ls -j | jq "index(\"$TOPIC\")")
	if [ "$index" = "null" ]; then
		nats stream add \
			--server $NATS_URI \
			--retention interest \
			--subjects $TOPIC \
			--defaults \
			--replicas $REPLICAS \
			--no-allow-direct \
			$TOPIC
	fi
`

	return jobs.Handle(ctx, topic, "create-topic", corev1.Container{
		Image: "natsio/nats-box:0.14.1",
		Name:  "create-topic",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", brokerURI.Host)),
			core.Env("TOPIC", topic.Name),
			core.Env("REPLICAS", func() string {
				if replicas := brokerURI.Query().Get("replicas"); replicas != "" {
					return replicas
				}
				return "1"
			}()),
		},
	})
}

func deleteJob(ctx core.Context, topic *v1beta1.BrokerTopic) error {
	const script = `
	set -xe
	index=$(nats --server $NATS_URI stream ls -j | jq "index(\"$TOPIC\")")
	if [ "$index" = "null" ]; then
		exit 0
	fi
	nats stream rm -f --server $NATS_URI $TOPIC
`
	return jobs.Handle(ctx, topic, "delete-topic", corev1.Container{
		Image: "natsio/nats-box:0.14.1",
		Name:  "delete-topic",
		Args:  core.ShellScript(script),
		Env: []corev1.EnvVar{
			core.Env("NATS_URI", fmt.Sprintf("nats://%s", topic.Status.URI.Host)),
			core.Env("TOPIC", topic.Name),
		},
	})
}
