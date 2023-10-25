package stack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/nats-io/nats.go"
	"github.com/opensearch-project/opensearch-go"

	"github.com/formancehq/operator/internal/modules"
	"github.com/formancehq/operator/internal/storage/es"
	natsStore "github.com/formancehq/operator/internal/storage/nats"
	"github.com/formancehq/operator/internal/storage/pg"
	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type StackFinalizer struct {
	name   string
	client client.Client
	log    logr.Logger

	reconcileConf *modules.ReconciliationConfig
}

const stacksIndex = "stacks"

var ErrNotFullyDisabled = errors.New("not fully disabled")
var ErrFinalizerNotRemoved = errors.New("finalizer not removed")
var ErrFinalizerNotAdded = errors.New("finalizer not added")
var ErrCast = errors.New("cannot cast interface to string")
var natsClientId = "membership"

func NewStackFinalizer(
	client client.Client,
	log logr.Logger,
	conf *modules.ReconciliationConfig,
) *StackFinalizer {
	s := &StackFinalizer{
		name:          "delete",
		client:        client,
		log:           log,
		reconcileConf: conf,
	}

	return s
}

func (f *StackFinalizer) RemoveFinalizer(ctx context.Context) error {
	if !controllerutil.ContainsFinalizer(f.reconcileConf.Stack, f.name) {
		return nil
	}

	updated := controllerutil.RemoveFinalizer(f.reconcileConf.Stack, f.name)
	if !updated {
		return ErrFinalizerNotRemoved
	}

	if err := f.client.Update(ctx, f.reconcileConf.Stack); err != nil {
		return err
	}
	return nil
}

func (f *StackFinalizer) HandleFinalizer(ctx context.Context, reqName string) (bool, error) {

	// examine DeletionTimestamp to determine if object is under deletion
	// The object is being deleted
	if !f.reconcileConf.Stack.ObjectMeta.DeletionTimestamp.IsZero() && controllerutil.ContainsFinalizer(f.reconcileConf.Stack, f.name) {
		// Make sur that the object is disable
		f.reconcileConf.Stack.Spec.Disabled = true
		if err := f.client.Update(ctx, f.reconcileConf.Stack); err != nil {
			return false, err
		}

		// We also need to make sure that all deployements are Terminated,
		// So that no one is already accessing ours databases.
		found := &appsv1.DeploymentList{}
		opts := &client.ListOptions{
			Namespace: reqName,
		}

		if err := f.client.List(ctx, found, opts); err != nil {
			return false, err
		}

		if len(found.Items) > 0 {
			return false, ErrNotFullyDisabled
		}

		// our StackFinalizer is present, so lets handle any external dependency
		err := f.deleteStack(ctx)
		if err != nil {
			// if fail to delete the external dependency here, return with error
			// so that it can be retried
			return false, err
		}

		// remove our StackFinalizer from the list and update it.
		updated := controllerutil.RemoveFinalizer(f.reconcileConf.Stack, f.name)
		if !updated {
			return false, ErrFinalizerNotRemoved
		}

		if err := f.client.Update(ctx, f.reconcileConf.Stack); err != nil {
			return false, err
		}

		// Stop reconciliation as the item is being deleted
		return true, nil

	}

	// The object is not being deleted, so if it does not have our StackFinalizer,
	// then lets add the StackFinalizer and update the object. This is equivalent
	// registering our StackFinalizer.
	if !controllerutil.ContainsFinalizer(f.reconcileConf.Stack, f.name) {
		updated := controllerutil.AddFinalizer(f.reconcileConf.Stack, f.name)
		if !updated {
			return false, ErrFinalizerNotAdded
		}
		if err := f.client.Update(ctx, f.reconcileConf.Stack); err != nil {
			return false, err
		}
	}

	return false, nil
}

func (f *StackFinalizer) retrieveModuleTopic() []string {
	subjectSet := collectionutils.NewSet[string]()
	values := reflect.ValueOf(f.reconcileConf.Configuration.Spec.Services)

	// TODO: iterate over registered modules ?
	for i := 0; i < values.NumField(); i++ {
		serviceName := strings.ToLower(values.Field(i).Type().Name())
		mod := modules.Get(strings.ToLower(serviceName))
		if mod == nil {
			continue
		}

		for _, v := range mod.Versions() {
			services := v.Services(*f.reconcileConf)
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
		subjects = append(subjects, fmt.Sprintf("%s-%s", f.reconcileConf.Stack.Name, k))
	}

	return subjects
}

func (f *StackFinalizer) deleteStack(ctx context.Context) error {
	f.log.Info("start deleting databases " + f.reconcileConf.Stack.Name)
	if err := f.deleteByService(ctx); err != nil {
		return err
	}

	f.log.Info("start deleting brokers subjects " + f.reconcileConf.Stack.Name)
	if err := f.deleteByBrokers(ctx); err != nil {
		return err
	}

	return nil
}

func (f *StackFinalizer) deleteByBrokers(ctx context.Context) error {
	values := reflect.ValueOf(f.reconcileConf.Configuration.Spec.Broker)
	for i := 0; i < values.NumField(); i++ {
		switch t := values.Field(i).Interface().(type) {
		case v1beta3.NatsConfig:
			if err := f.deleleNatsSubjects(ctx, &t, f.reconcileConf.Stack.Name, f.retrieveModuleTopic()); err != nil {
				return err
			}
		}

	}
	return nil
}

func (f *StackFinalizer) deleleNatsSubjects(ctx context.Context, config *v1beta3.NatsConfig, stackName string, subjectService []string) error {
	client, err := natsStore.NewClient(config, natsClientId)
	if err != nil {
		f.log.Error(err, "NATS: client")
		return err
	}
	defer client.Close()

	jsCtx, err := client.JetStream()
	if err != nil {
		return err
	}

	for _, service := range subjectService {
		stackServiceSubject := fmt.Sprintf("%s-%s", stackName, service)
		exist, err := existSubject(ctx, jsCtx, stackServiceSubject)
		if err != nil {
			return err
		}

		// It meens it has already been deleted, and just not exists anymore
		if !exist {
			continue
		}

		err = jsCtx.DeleteStream(stackServiceSubject, nats.Context(ctx))
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *StackFinalizer) deleteByService(ctx context.Context) error {

	values := reflect.ValueOf(f.reconcileConf.Configuration.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			switch t := servicesValues.Field(j).Interface().(type) {
			case v1beta3.PostgresConfig:
				serviceName := strings.ToLower(values.Type().Field(i).Name)
				if err := f.deletePostgresDb(ctx, serviceName, t); err != nil {
					return err
				}

			case v1beta3.ElasticSearchConfig:
				client, err := es.NewElasticSearchClient(t)
				if err != nil {
					f.log.Error(err, "ELK: client")
					return err
				}

				if err := f.dropESIndex(ctx, client); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (f *StackFinalizer) dropESIndex(ctx context.Context, client *opensearch.Client) error {
	var (
		buf bytes.Buffer
	)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"stack": f.reconcileConf.Stack.Name,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		f.log.Error(err, "ELK: Error during json encoding")

		return err
	}
	body := bytes.NewReader(buf.Bytes())
	response, err := client.DeleteByQuery([]string{stacksIndex}, body, client.DeleteByQuery.WithContext(ctx))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		return fmt.Errorf("ELK status: %d", response.StatusCode)
	}

	return nil
}

func (f *StackFinalizer) deletePostgresDb(
	ctx context.Context,
	serviceName string,
	postgresConfig v1beta3.PostgresConfig,
) error {
	client, err := pg.OpenClient(postgresConfig)
	if err != nil {
		f.log.Error(err, "PG: Cannot open pg client")
		return err
	}
	defer client.Close()

	if err := pg.DropDB(client, f.reconcileConf.Stack.Name, serviceName, ctx); err != nil {
		f.log.Error(err, "PG: Error during drop")
		return err
	}

	f.log.Info(fmt.Sprintf("PG: database \"%s-%s\" droped", f.reconcileConf.Stack.Name, serviceName))

	return nil
}

func existSubject(ctx context.Context, natsContext nats.JetStreamContext, subject string) (bool, error) {
	_, err := natsContext.StreamNameBySubject(subject, nats.Context(ctx))
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}
