//nolint:nosnakecase
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	"github.com/alitto/pond"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -source=/src/ee/agent/internal/membership_listener.go -destination=/src/ee/agent/internal/membership_client_generated.go -package=internal . MembershipClient
type MembershipClient interface {
	Orders() chan *generated.Order
	Send(message *generated.Message) error
}

type MembershipClientMock struct {
	mu       sync.Mutex
	orders   chan *generated.Order
	messages []*generated.Message
}

func (m MembershipClientMock) Orders() chan *generated.Order {
	return m.orders
}

func (m *MembershipClientMock) Send(message *generated.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, message)
	return nil
}

func (m *MembershipClientMock) GetMessages() []*generated.Message {
	return m.messages
}

func NewMembershipClientMock() *MembershipClientMock {
	return &MembershipClientMock{
		orders: make(chan *generated.Order),
	}
}

type ClientInfo struct {
	ID         string
	BaseUrl    *url.URL
	Production bool
	Version    string
}

type membershipListener struct {
	clientInfo ClientInfo
	client     K8SClient

	restMapper meta.RESTMapper
	orders     MembershipClient
	wp         *pond.WorkerPool
}

func (c *membershipListener) Start(ctx context.Context) {
	defer c.wp.StopAndWait()
	for {
		select {
		case msg, ok := <-c.orders.Orders():
			if !ok {
				return
			}

			c.wp.Submit(func() {
				ctx, span := otel.GetTracerProvider().Tracer("com.formance.agent").Start(ctx, "newOrder", trace.WithNewRoot())
				defer span.End()
				logging.FromContext(ctx).
					WithField("traceId", span.SpanContext().TraceID()).
					WithField("spanId", span.SpanContext().SpanID()).
					Infof("Got message from membership: %T", msg.GetMessage())
				switch msg := msg.Message.(type) {
				case *generated.Order_ExistingStack:
					span.SetName("syncExistingStack")
					c.syncExistingStack(ctx, msg.ExistingStack)
				case *generated.Order_DeletedStack:
					span.SetName("deleteStack")
					c.deleteStack(ctx, msg.DeletedStack)
				case *generated.Order_DisabledStack:
					span.SetName("disableStack")
					c.disableStack(ctx, msg.DisabledStack)
				case *generated.Order_EnabledStack:
					span.SetName("enableStack")
					c.enableStack(ctx, msg.EnabledStack)
				}
			})
		}
	}
}

func (c *membershipListener) syncExistingStack(ctx context.Context, membershipStack *generated.Stack) {
	versions := membershipStack.Versions
	if versions == "" {
		versions = "default"
	}

	metadata := c.generateMetadata(membershipStack)

	stack, err := c.createOrUpdate(ctx, v1beta1.GroupVersion.WithKind("Stack"), membershipStack.ClusterName, membershipStack.ClusterName, nil, map[string]any{
		"metadata": metadata,
		"spec": map[string]any{
			"versionsFromFile": versions,
			"disabled":         membershipStack.Disabled,
			"enableAudit":      membershipStack.EnableAudit,
		},
	})
	if err != nil {
		logging.FromContext(ctx).Errorf("Unable to create stack cluster side: %s", err)
		return
	}

	c.syncModules(ctx, metadata, stack, membershipStack)
	c.syncStargate(ctx, metadata, stack, membershipStack)
	c.syncAuthClients(ctx, metadata, stack, membershipStack.StaticClients)

	logging.FromContext(ctx).Infof("Stack %s updated cluster side", stack.GetName())
}

func (c *membershipListener) generateMetadata(membershipStack *generated.Stack) map[string]any {
	additionalLabels := map[string]any{}
	for key, value := range membershipStack.AdditionalLabels {
		additionalLabels["formance.com/"+key] = value
	}

	additionalAnnotations := map[string]any{}
	for key, value := range membershipStack.AdditionalAnnotations {
		additionalAnnotations["formance.com/"+key] = value
	}

	return map[string]any{
		"annotations": additionalAnnotations,
		"labels":      additionalLabels,
	}

}
func (c *membershipListener) syncModules(ctx context.Context, metadata map[string]any, stack *unstructured.Unstructured, membershipStack *generated.Stack) {
	modules := collectionutils.Map(membershipStack.Modules, func(module *generated.Module) string {
		return strings.ToLower(module.Name)
	})
	logger := logging.FromContext(ctx).WithField("stack", membershipStack.ClusterName)
	logger.Infof("Syncing modules for stack %s", membershipStack.Modules)
	for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
		object := reflect.New(rtype).Interface()
		if _, ok := object.(v1beta1.Module); !ok {
			continue
		}

		if gvk.Kind == "Stargate" {
			continue
		}
		resources, err := c.restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			logger.Errorf("Unable to get resources for %s: %s", gvk.Kind, err)
			continue
		}
		singular, err := c.restMapper.ResourceSingularizer(resources.Resource.Resource)
		if err != nil {
			logger.Errorf("Unable to get singular for %s: %s", gvk.Kind, err)
			continue
		}
		logger.Debugf("Resource: checking module resource %s, singular: %s", resources.Resource.Resource, singular)
		if !slices.Contains(modules, singular) {
			if err := c.deleteModule(ctx, logger, resources.Resource.Resource, stack.GetName()); err != nil {
				logger.Errorf("Unable to get and delete module %s cluster side: %s", gvk.Kind, err)
			}
			continue
		}

		switch gvk.Kind {
		case "Auth":
			if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, gvk, map[string]any{
				"metadata": metadata,
				"spec": map[string]any{
					"delegatedOIDCServer": map[string]any{
						"issuer":       membershipStack.AuthConfig.Issuer,
						"clientID":     membershipStack.AuthConfig.ClientId,
						"clientSecret": membershipStack.AuthConfig.ClientSecret,
					},
				},
			}); err != nil {
				logger.Errorf("Unable to create module Auth cluster side: %s", err)
			}
		case "Gateway":
			if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, gvk, map[string]any{
				"metadata": metadata,
				"spec": map[string]any{
					"ingress": map[string]any{
						"host":   fmt.Sprintf("%s.%s", stack.GetName(), c.clientInfo.BaseUrl.Host),
						"scheme": c.clientInfo.BaseUrl.Scheme,
					},
				},
			}); client.IgnoreNotFound(err) != nil {
				logger.Errorf("Unable to create module Stargate cluster side: %s", err)
			}
		default:
			if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, gvk, map[string]any{
				"metadata": metadata,
			}); err != nil {
				logger.Errorf("Unable to create module %s cluster side: %s", gvk.Kind, err)
			}
		}

	}
}

func (c *membershipListener) deleteModule(ctx context.Context, logger logging.Logger, resource string, stackName string) error {
	logger.Debugf("Deleting module %s", resource)

	return c.client.EnsureNotExistsBySelector(ctx, resource, stackLabels(stackName))
}

func (c *membershipListener) syncStargate(ctx context.Context, metadata map[string]any, stack *unstructured.Unstructured, membershipStack *generated.Stack) {
	stargateName := fmt.Sprintf("%s-stargate", membershipStack.ClusterName)
	if membershipStack.StargateConfig != nil && membershipStack.StargateConfig.Enabled {
		parts := strings.Split(stack.GetName(), "-")

		if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, v1beta1.GroupVersion.WithKind("Stargate"), map[string]any{
			"metadata": metadata,
			"spec": map[string]any{
				"organizationID": parts[0],
				"stackID":        parts[1],
				"serverURL":      membershipStack.StargateConfig.Url,
				"auth": map[string]any{
					"issuer":       membershipStack.AuthConfig.Issuer,
					"clientID":     membershipStack.AuthConfig.ClientId,
					"clientSecret": membershipStack.AuthConfig.ClientSecret,
				},
			},
		}); err != nil {
			logging.FromContext(ctx).Errorf("Unable to create module Stargate cluster side: %s", err)
		}
	} else {
		if err := c.client.EnsureNotExists(ctx, "Stargates", stargateName); err != nil {
			logging.FromContext(ctx).Errorf("Unable to delete module Stargate cluster side: %s", err)
		}
	}
}

func (c *membershipListener) syncAuthClients(ctx context.Context, metadata map[string]any, stack *unstructured.Unstructured, staticClients []*generated.AuthClient) {
	expectedAuthClients := make([]*unstructured.Unstructured, 0)
	for _, client := range staticClients {
		authClient, err := c.createOrUpdateStackDependency(ctx, fmt.Sprintf("%s-%s", stack.GetName(), client.Id), stack.GetName(),
			stack, v1beta1.GroupVersion.WithKind("AuthClient"), map[string]any{
				"metadata": metadata,
				"spec": map[string]any{
					"id":     client.Id,
					"public": client.Public,
				},
			})
		if err != nil {
			logging.FromContext(ctx).Errorf("Unable to create AuthClient cluster side: %s", err)
			continue
		}
		expectedAuthClients = append(expectedAuthClients, authClient)
	}

	authClients, err := c.client.List(ctx, "AuthClients", stackLabels(stack.GetName()))
	if err != nil {
		logging.FromContext(ctx).Errorf("Unable to list AuthClient cluster side: %s", err)
		return
	}

	authClientsToDelete := collectionutils.Reduce(authClients, func(acc []string, item unstructured.Unstructured) []string {
		for _, expectedClient := range expectedAuthClients {
			if expectedClient.GetName() == item.GetName() {
				return acc
			}
		}
		return append(acc, item.GetName())
	}, []string{})

	for _, name := range authClientsToDelete {
		logging.FromContext(ctx).Infof("Deleting AuthClient %s", name)
		if err := c.client.EnsureNotExists(ctx, "AuthClients", name); err != nil {
			logging.FromContext(ctx).Errorf("Unable to delete AuthClient %s cluster side: %s", name, err)
		}
	}
}

func (c *membershipListener) deleteStack(ctx context.Context, stack *generated.DeletedStack) {
	if err := c.client.EnsureNotExists(ctx, "Stacks", stack.ClusterName); err != nil {
		logging.FromContext(ctx).Errorf("Deleting cluster side: %s", err)
		return
	}
	logging.FromContext(ctx).Infof("Stack %s deleted", stack.ClusterName)
}

func (c *membershipListener) disableStack(ctx context.Context, stack *generated.DisabledStack) {

	if err := c.client.Patch(ctx, "Stacks", stack.ClusterName, []byte(`{"spec": {"disabled": true}}`)); err != nil {
		logging.FromContext(ctx).Errorf("Disabling cluster side: %s", err)
		return
	}

	logging.FromContext(ctx).Infof("Stack %s disabled", stack.ClusterName)
}

func (c *membershipListener) enableStack(ctx context.Context, stack *generated.EnabledStack) {
	if err := c.client.Patch(ctx, "Stacks", stack.ClusterName, []byte(`{"spec": {"disabled": false}}`)); err != nil {
		logging.FromContext(ctx).Errorf("Disabling cluster side: %s", err)
		return
	}

	logging.FromContext(ctx).Infof("Stack %s enabled", stack.ClusterName)
}

func (c *membershipListener) createOrUpdate(ctx context.Context, gvk schema.GroupVersionKind, name string, stackName string, owner *metav1.OwnerReference, content map[string]any) (*unstructured.Unstructured, error) {

	logger := logging.FromContext(ctx).WithFields(map[string]any{
		"gvk": gvk,
	})
	logger.Infof("creating object '%s'", name)
	if content["metadata"] == nil {
		content["metadata"] = map[string]any{}
	}

	if content["metadata"].(map[string]any)["labels"] == nil {
		content["metadata"].(map[string]any)["labels"] = map[string]any{}
	}

	content["metadata"].(map[string]any)["labels"].(map[string]any)["formance.com/created-by-agent"] = "true"
	content["metadata"].(map[string]any)["labels"].(map[string]any)["formance.com/stack"] = stackName
	content["metadata"].(map[string]any)["name"] = name

	restMapping, err := c.restMapper.RESTMapping(gvk.GroupKind())
	if err != nil {
		return nil, errors.Wrap(err, "getting rest mapping")
	}

	u, err := c.client.Get(ctx, restMapping.Resource.Resource, name)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, errors.Wrap(err, "reading object")
		}

		logger.Infof("Object not found, create a new one")

		u := &unstructured.Unstructured{}
		u.SetUnstructuredContent(content)
		u.SetGroupVersionKind(gvk)
		u.SetName(name)
		if owner != nil {
			u.SetOwnerReferences([]metav1.OwnerReference{*owner})
		}

		if err := c.client.Create(ctx, restMapping.Resource.Resource, u); err != nil {
			return nil, errors.Wrap(err, "creating object")
		}

		return u, nil

	}

	if equality.Semantic.DeepDerivative(content, u.Object) {
		logger.Infof("Object found and has expected content, skip it")
		return u, nil
	}

	logger.Infof("Object exists and content differ, patch it")
	contentData, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	if err := c.client.Patch(ctx, restMapping.Resource.Resource, name, contentData); err != nil {
		return nil, errors.Wrap(err, "patching object")
	}

	return u, nil
}

func (c *membershipListener) createOrUpdateStackDependency(
	ctx context.Context,
	name string,
	stackName string,
	stack *unstructured.Unstructured,
	gvk schema.GroupVersionKind,
	content map[string]any,
) (*unstructured.Unstructured, error) {
	if _, ok := content["spec"]; !ok {
		content["spec"] = map[string]any{}
	}
	content["spec"].(map[string]any)["stack"] = stack.GetName()

	return c.createOrUpdate(ctx, gvk, name, stackName,
		&metav1.OwnerReference{
			APIVersion: "formance.com/v1beta1",
			Kind:       "Stack",
			Name:       stack.GetName(),
			UID:        stack.GetUID(),
		}, content)
}

func NewMembershipListener(client K8SClient, clientInfo ClientInfo, mapper meta.RESTMapper,
	orders MembershipClient) *membershipListener {
	return &membershipListener{
		client:     client,
		clientInfo: clientInfo,
		restMapper: mapper,
		orders:     orders,
		wp:         pond.New(5, 5),
	}
}

func must[T any](t *T, err error) T {
	if err != nil {
		panic(err)
	}
	return *t
}

func stackLabels(stackName string) labels.Selector {
	return labels.NewSelector().Add(
		must(labels.NewRequirement("formance.com/created-by-agent", selection.Equals, []string{"true"})),
		must(labels.NewRequirement("formance.com/stack", selection.Equals, []string{stackName})),
	)
}
