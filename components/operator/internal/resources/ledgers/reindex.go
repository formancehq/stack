package ledgers

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createReindexCronJob(ctx core.Context, ledger *v1beta1.Ledger) (*batchv1.CronJob, error) {
	cronJob, _, err := core.CreateOrUpdate[*batchv1.CronJob](ctx, types.NamespacedName{
		Namespace: ledger.Spec.Stack,
		Name:      fmt.Sprintf("reindex-ledger"),
	}, func(t *batchv1.CronJob) error {
		t.Spec = batchv1.CronJobSpec{
			Suspend:  pointer.For(true),
			Schedule: "* * * * *",
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					TTLSecondsAfterFinished: pointer.For(int32(30)),
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{{
								Image: "curlimages/curl:8.2.1",
								Name:  "reindex-ledger",
								Command: core.ShellScript(`
					curl http://benthos.%s.svc.cluster.local:4195/ledger_reindex_all -X POST -H 'Content-Type: application/json' -d '{}'`, ledger.Spec.Stack),
							}},
						},
					},
				},
			},
		}

		return nil
	}, core.WithController[*batchv1.CronJob](ctx.GetScheme(), ledger))
	return cronJob, err
}
