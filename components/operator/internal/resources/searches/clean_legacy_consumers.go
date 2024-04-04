package searches

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/jobs"
	"github.com/formancehq/operator/internal/resources/settings"
	corev1 "k8s.io/api/core/v1"
)

func cleanConsumers(ctx Context, search *v1beta1.Search) error {

	brokerURI, err := settings.RequireURL(ctx, search.Spec.Stack, "broker.dsn")
	if err != nil {
		return err
	}

	if brokerURI == nil {
		return nil
	}

	const script = `
	set -xe
	for service in ledger payments audit; do
		for consumer in search-ledgerv2 search-payments-resets search-audit; do 
			index=$(nats --server $NATS_URI consumer ls $STACK-$service -j | jq "index(\"$consumer\")")
			if [ "$index" != "null" ]; then
				nats --server $NATS_URI consumer rm $STACK-$service $consumer -f
			fi
		done
	done
`
	return jobs.Handle(ctx, search, "clean-consumers", corev1.Container{
		Image: "natsio/nats-box:0.14.1",
		Name:  "delete-consumer",
		Args:  ShellScript(script),
		Env: []corev1.EnvVar{
			Env("NATS_URI", fmt.Sprintf("nats://%s", brokerURI.Host)),
			Env("STACK", search.Spec.Stack),
		},
	})
}
