package modules

import (
	"context"
	"fmt"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	batchv1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type JobRunner interface {
	RunJob(ctx context.Context, jobName string, preRun func() error, modifier func(t *batchv1.Job)) (bool, error)
}

type defaultJobRunner struct {
	Deployer
	jobNamePrefix string
}

func (d defaultJobRunner) RunJob(ctx context.Context, jobName string, preRun func() error, modifier func(t *batchv1.Job)) (bool, error) {
	logger := log.FromContext(ctx)
	if job, err := d.Jobs().Get(ctx, d.jobNamePrefix+jobName); err != nil {
		if !apierrors.IsNotFound(err) {
			return false, err
		}

		if preRun != nil {
			if err := preRun(); err != nil {
				return false, err
			}
		}

		_, err := d.Jobs().CreateOrUpdate(ctx, d.jobNamePrefix+jobName, modifier)
		if err != nil {
			return false, err
		}
		return false, nil
	} else {
		logger.Info(fmt.Sprintf("Job found, succeded: %d", job.Status.Succeeded))

		return job.Status.Succeeded > 0, nil
	}
}

func NewJobRunner(client client.Client, scheme *runtime.Scheme, stack *v1beta3.Stack, owner client.Object, jobNamePrefix string) *defaultJobRunner {
	return &defaultJobRunner{
		Deployer:      NewScopedDeployer(client, scheme, stack, owner),
		jobNamePrefix: jobNamePrefix,
	}
}

var _ JobRunner = &defaultJobRunner{}
