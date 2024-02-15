package jobs

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type handleJobConfiguration struct {
	preCreate func() error
	mutators  []core.ObjectMutator[*batchv1.Job]
}

type HandleJobOption func(configuration *handleJobConfiguration)

func PreCreate(preCreate func() error) HandleJobOption {
	return func(configuration *handleJobConfiguration) {
		configuration.preCreate = preCreate
	}
}

func Mutator(mutator core.ObjectMutator[*batchv1.Job]) HandleJobOption {
	return func(configuration *handleJobConfiguration) {
		configuration.mutators = append(configuration.mutators, mutator)
	}
}

func WithServiceAccount(serviceAccountName string) HandleJobOption {
	return Mutator(func(t *batchv1.Job) error {
		t.Spec.Template.Spec.ServiceAccountName = serviceAccountName
		return nil
	})
}

func Handle(ctx core.Context, owner v1beta1.Dependent, jobName string, container v1.Container, options ...HandleJobOption) error {

	configuration := &handleJobConfiguration{}
	for _, option := range options {
		option(configuration)
	}

	jobName = fmt.Sprintf("%s-%s", owner.GetUID(), jobName)
	job := &batchv1.Job{}
	err := ctx.GetClient().Get(ctx, types.NamespacedName{
		Namespace: owner.GetStack(),
		Name:      jobName,
	}, job)
	if err != nil && client.IgnoreNotFound(err) != nil {
		return err
	}

	if job.Status.Succeeded > 0 {
		return nil
	}

	if err == nil { // Job found
		if !equality.Semantic.DeepDerivative(container, job.Spec.Template.Spec.Containers[0]) {
			if err := ctx.GetClient().Delete(ctx, job, &client.DeleteOptions{
				GracePeriodSeconds: pointer.For(int64(0)),
				PropagationPolicy:  pointer.For(metav1.DeletePropagationForeground),
			}); err != nil {
				return err
			}
		} else {
			return core.NewPendingError()
		}
	}

	if configuration.preCreate != nil {
		err := configuration.preCreate()
		if err != nil {
			return errors.Wrap(err, "in precreate")
		}
	}

	job = &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: owner.GetStack(),
			Name:      jobName,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            pointer.For(int32(10000)),
			TTLSecondsAfterFinished: pointer.For(int32(30)),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyOnFailure,
					Containers:    []v1.Container{container},
				},
			},
		},
	}

	for _, mutator := range configuration.mutators {
		if err := mutator(job); err != nil {
			return err
		}
	}

	if err := controllerutil.SetControllerReference(owner, job, ctx.GetScheme()); err != nil {
		return err
	}
	if err := ctx.GetClient().Create(ctx, job); err != nil {
		return err
	}

	return core.NewPendingError()
}
