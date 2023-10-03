package stack

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/formancehq/operator/internal/modules"
	"github.com/formancehq/operator/internal/storage/es"
	"github.com/formancehq/operator/internal/storage/nats"
	"github.com/formancehq/operator/internal/storage/pg"
	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type StackFinalizer struct {
	name   string
	ctx    context.Context
	client client.Client
	log    logr.Logger

	reconcileConf modules.ReconciliationConfig
}

func (f *StackFinalizer) retrieveModuleTopic() []string {
	subjectSet := collectionutils.NewSet[string]()
	values := reflect.ValueOf(f.reconcileConf.Configuration.Spec.Services)

	// TODO: iterate over registered modules ?
	for i := 0; i < values.NumField(); i++ {
		serviceName := strings.ToLower(values.Field(i).Type().Name())
		mod := modules.Get(serviceName)
		for _, v := range mod.Versions() {
			services := v.Services(f.reconcileConf)
			for _, service := range services {
				if service.NeedTopic {
					subjectSet.Put(service.Name)
					break
				}
			}
		}

	}

	subjects := []string{}
	for k := range subjectSet {
		subjects = append(subjects, k)
	}

	return subjects

}

func (f *StackFinalizer) deleteStack(ctx context.Context, stack *stackv1beta3.Stack, conf *stackv1beta3.Configuration, log logr.Logger) (*ctrl.Result, error) {
	log.Info("start deleting databases " + stack.Name)
	if err := f.DeleteByService(conf, stack.Name, log); err != nil {
		return &ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	log.Info("start deleting brokers subjects " + stack.Name)
	if err := f.DeleteByBrokers(conf, stack.Name, f.retrieveModuleTopic(), log); err != nil {
		return &ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	return nil, nil
}

var ErrNotFullyDisabled = errors.New("not fully disabled")

func (f *StackFinalizer) HandleFinalizer(ctx context.Context, log logr.Logger, stack *stackv1beta3.Stack, conf *stackv1beta3.Configuration, req ctrl.Request) (bool, error) {

	// examine DeletionTimestamp to determine if object is under deletion
	// The object is being deleted
	if !stack.ObjectMeta.DeletionTimestamp.IsZero() && controllerutil.ContainsFinalizer(stack, f.name) {
		// Make sur that the object is disable
		stack.Spec.Disabled = true
		if err := f.client.Update(ctx, stack); err != nil {
			return false, err
		}

		// We also need to make sure that all deployements are Terminated,
		// So that no one is already accessing ours databases.
		found := &appsv1.DeploymentList{}
		opts := &client.ListOptions{
			Namespace: req.Name,
		}

		if err := f.client.List(ctx, found, opts); err != nil {
			return false, err
		}

		if len(found.Items) > 0 {
			return false, ErrNotFullyDisabled
		}

		// our StackFinalizer is present, so lets handle any external dependency
		_, err := f.deleteStack(ctx, stack, conf, log)
		if err != nil {
			// if fail to delete the external dependency here, return with error
			// so that it can be retried
			return false, err
		}

		// remove our StackFinalizer from the list and update it.
		controllerutil.RemoveFinalizer(stack, f.name)
		if err := f.client.Update(ctx, stack); err != nil {
			return false, err
		}

		// Stop reconciliation as the item is being deleted
		return true, nil

	}

	// The object is not being deleted, so if it does not have our StackFinalizer,
	// then lets add the StackFinalizer and update the object. This is equivalent
	// registering our StackFinalizer.
	if !controllerutil.ContainsFinalizer(stack, f.name) {
		controllerutil.AddFinalizer(stack, f.name)
		if err := f.client.Update(ctx, stack); err != nil {
			return false, err
		}
	}

	return false, nil
}

var (
	natsClientId = "membership"
)
var (
	ErrCast = errors.New("cannot cast interface to string")
)

func (f *StackFinalizer) DeleteByBrokers(c *v1beta3.Configuration, stackName string, subjectService []string, logger logr.Logger) error {
	values := reflect.ValueOf(c.Spec.Broker)
	for i := 0; i < values.NumField(); i++ {
		switch t := values.Field(i).Interface().(type) {
		case v1beta3.NatsConfig:
			if err := f.deleleNatsSubjects(&t, stackName, subjectService, logger); err != nil {
				return err
			}
		}

	}
	return nil
}

func (f *StackFinalizer) deleleNatsSubjects(config *v1beta3.NatsConfig, stackName string, subjectService []string, logger logr.Logger) error {
	client, err := nats.NewClient(config, natsClientId)
	if err != nil {
		logger.Error(err, "NATS: client")
		return err
	}
	defer client.Close()

	jsCtx, err := client.JetStream()
	if err != nil {
		return err
	}

	for _, service := range subjectService {
		stackServiceSubject := fmt.Sprintf("%s-%s", stackName, service)
		exist, err := nats.ExistSubject(jsCtx, stackServiceSubject, logger)
		if err != nil {
			return err
		}

		// It meens it has already been deleted, and just not exists anymore
		if !exist {
			continue
		}

		err = jsCtx.DeleteStream(stackServiceSubject)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *StackFinalizer) DeleteByService(c *v1beta3.Configuration, stackName string, logger logr.Logger) error {

	values := reflect.ValueOf(c.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			switch t := servicesValues.Field(j).Interface().(type) {
			case v1beta3.PostgresConfig:
				serviceName := strings.ToLower(values.Type().Field(i).Name)
				if err := f.deletePostgresDb(t, stackName, serviceName, logger); err != nil {
					return err
				}

			case v1beta3.ElasticSearchConfig:
				client, err := es.NewElasticSearchClient(&t)
				if err != nil {
					logger.Error(err, "ELK: client")
					return err
				}

				if err := es.DropESIndex(client, logger, stackName); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (f *StackFinalizer) deletePostgresDb(
	postgresConfig v1beta3.PostgresConfig,
	stackName string,
	serviceName string,
	logger logr.Logger,
) error {
	client, err := pg.OpenClient(postgresConfig)
	if err != nil {
		logger.Error(err, "PG: Cannot open pg client")
		return err
	}
	defer client.Close()

	if err := pg.DropDB(client, stackName, serviceName); err != nil {
		logger.Error(err, "PG: Error during drop")
		return err
	}

	logger.Info(fmt.Sprintf("PG: database \"%s-%s\" droped", stackName, serviceName))

	return nil
}

func NewStackFinalizer(
	context context.Context,
	client client.Client,
	log logr.Logger,
	conf modules.ReconciliationConfig,
) *StackFinalizer {
	return &StackFinalizer{
		name:          "delete",
		ctx:           context,
		client:        client,
		log:           log,
		reconcileConf: conf,
	}
}
