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

	"github.com/alitto/pond"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MembershipClient interface {
	Orders() chan *generated.Order
	Send(message *generated.Message) error
}

type MembershipClientMock struct {
	orders   chan *generated.Order
	messages []*generated.Message
}

func (m MembershipClientMock) Orders() chan *generated.Order {
	return m.orders
}

func (m *MembershipClientMock) Send(message *generated.Message) error {
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
	restClient *rest.RESTClient

	stacksModules InMemoryStacksModules
	restMapper    meta.RESTMapper
	orders        MembershipClient
	wp            *pond.WorkerPool
}

func getExpectedModules() []string {
	return []string{
		"Stargate", "Wallets", "Ledger", "Payments", "Webhooks", "Auth", "Orchestration", "Search", "Gateway",
	}
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
				sharedlogging.FromContext(ctx).Infof("Got message from membership: %T", msg.GetMessage())
				switch msg := msg.Message.(type) {
				case *generated.Order_ExistingStack:
					if msg.ExistingStack.Modules == nil || len(msg.ExistingStack.Modules) == 0 {
						msg.ExistingStack.Modules = make([]*generated.Module, 0)
						for _, module := range getExpectedModules() {
							msg.ExistingStack.Modules = append(msg.ExistingStack.Modules, &generated.Module{
								Name: module,
							})
						}
					}
					c.stacksModules[msg.ExistingStack.ClusterName] = collectionutils.Map(msg.ExistingStack.Modules, func(module *generated.Module) string {
						return strings.ToUpper(string(module.Name[0])) + module.Name[1:]
					})
					c.syncExistingStack(ctx, msg.ExistingStack)
				case *generated.Order_DeletedStack:
					c.deleteStack(ctx, msg.DeletedStack)
				case *generated.Order_DisabledStack:
					c.disableStack(ctx, msg.DisabledStack)
				case *generated.Order_EnabledStack:
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
		sharedlogging.FromContext(ctx).Errorf("Unable to create stack cluster side: %s", err)
		return
	}

	c.syncModules(ctx, metadata, stack, membershipStack)
	c.syncStargate(ctx, metadata, stack, membershipStack)
	c.syncAuthClients(ctx, metadata, stack, membershipStack.StaticClients)

	sharedlogging.FromContext(ctx).Infof("Stack %s updated cluster side", stack.GetName())
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

	for gvk, rtype := range scheme.Scheme.AllKnownTypes() {
		object := reflect.New(rtype).Interface()
		if _, ok := object.(v1beta1.Module); !ok {
			continue
		}

		if !slices.Contains(c.stacksModules[stack.GetName()], gvk.Kind) {
			if err := c.deleteModule(ctx, gvk, stack.GetName()); err != nil {
				sharedlogging.FromContext(ctx).Errorf("Unable to get and delete module %s cluster side: %s", gvk.Kind, err)
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
				sharedlogging.FromContext(ctx).Errorf("Unable to create module Auth cluster side: %s", err)
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
				sharedlogging.FromContext(ctx).Errorf("Unable to create module Stargate cluster side: %s", err)
			}
		default:
			if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, gvk, map[string]any{
				"metadata": metadata,
			}); err != nil {
				sharedlogging.FromContext(ctx).Errorf("Unable to create module %s cluster side: %s", gvk.Kind, err)
			}
		}

	}
}

func (c *membershipListener) deleteModule(ctx context.Context, gvk schema.GroupVersionKind, stackName string) error {
	logging.FromContext(ctx).Debugf("Deleting module %s", gvk.Kind)
	if err := c.restClient.Delete().
		Resource(gvk.Kind).
		VersionedParams(
			&metav1.ListOptions{
				LabelSelector: "formance.com/created-by-agent=true,formance.com/stack=" + stackName,
			}, scheme.ParameterCodec).
		Do(ctx).
		Error(); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return errors.Wrap(err, "Unable to delete object")
	}

	return nil
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
			sharedlogging.FromContext(ctx).Errorf("Unable to create module Stargate cluster side: %s", err)
		}
	} else {
		if err := c.restClient.Delete().
			Name(stargateName).
			Do(ctx).
			Error(); client.IgnoreNotFound(err) != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to delete module Stargate cluster side: %s", err)
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
			sharedlogging.FromContext(ctx).Errorf("Unable to create AuthClient cluster side: %s", err)
			continue
		}
		expectedAuthClients = append(expectedAuthClients, authClient)
	}

	authClientList := &unstructured.UnstructuredList{}
	if err := c.restClient.Get().
		Resource("AuthClients").
		VersionedParams(&metav1.ListOptions{
			LabelSelector: "formance.com/created-by-agent=true,formance.com/stack=" + stack.GetName(),
		}, scheme.ParameterCodec).
		Do(ctx).
		Into(authClientList); err != nil {
		sharedlogging.FromContext(ctx).Errorf("Unable to list AuthClient cluster side: %s", err)
	}

	authClientsToDelete := collectionutils.Reduce(authClientList.Items, func(acc []string, item unstructured.Unstructured) []string {
		for _, expectedClient := range expectedAuthClients {
			if expectedClient.GetName() == item.GetName() {
				return acc
			}
		}
		return append(acc, item.GetName())
	}, []string{})

	for _, name := range authClientsToDelete {
		sharedlogging.FromContext(ctx).Infof("Deleting AuthClient %s", name)
		if err := c.restClient.Delete().
			Resource("AuthClients").
			Name(name).
			Do(ctx).
			Error(); err != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to delete AuthClient %s cluster side: %s", name, err)
		}
	}
}

func (c *membershipListener) deleteStack(ctx context.Context, stack *generated.DeletedStack) {
	if err := c.restClient.Delete().
		Resource("stacks").
		Name(stack.ClusterName).
		Do(ctx).
		Error(); err != nil {
		if apierrors.IsNotFound(err) {
			sharedlogging.FromContext(ctx).Infof("Cannot delete not existing stack: %s", stack.ClusterName)
			return
		}
		sharedlogging.FromContext(ctx).Errorf("Deleting cluster side: %s", err)
		return
	}
	sharedlogging.FromContext(ctx).Infof("Stack %s deleted", stack.ClusterName)
}

func (c *membershipListener) disableStack(ctx context.Context, stack *generated.DisabledStack) {
	if err := c.restClient.Patch(types.MergePatchType).
		Name(stack.ClusterName).
		Body([]byte(`{"spec": {"disabled": true}}`)).
		Resource("Stacks").
		Do(ctx).
		Error(); err != nil {
		if apierrors.IsNotFound(err) {
			sharedlogging.FromContext(ctx).Infof("Cannot disable not existing stack: %s", stack.ClusterName)
			return
		}
		sharedlogging.FromContext(ctx).Errorf("Disabling cluster side: %s", err)
		return
	}

	sharedlogging.FromContext(ctx).Infof("Stack %s disabled", stack.ClusterName)
}

func (c *membershipListener) enableStack(ctx context.Context, stack *generated.EnabledStack) {
	if err := c.restClient.Patch(types.MergePatchType).
		Name(stack.ClusterName).
		Body([]byte(`{"spec": {"disabled": false}}`)).
		Resource("Stacks").
		Do(ctx).
		Error(); err != nil {
		if apierrors.IsNotFound(err) {
			sharedlogging.FromContext(ctx).Infof("Cannot enable not existing stack: %s", stack.ClusterName)
			return
		}
		sharedlogging.FromContext(ctx).Errorf("Enabling cluster side: %s", err)
		return
	}

	sharedlogging.FromContext(ctx).Infof("Stack %s enabled", stack.ClusterName)
}

func (c *membershipListener) createOrUpdate(ctx context.Context, gvk schema.GroupVersionKind, name string, stackName string, owner *metav1.OwnerReference, content map[string]any) (*unstructured.Unstructured, error) {

	logger := sharedlogging.FromContext(ctx).WithFields(map[string]any{
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

	u := &unstructured.Unstructured{}
	if err := c.restClient.Get().
		Resource(restMapping.Resource.Resource).
		Name(name).
		Do(ctx).
		Into(u); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, errors.Wrap(err, "reading object")
		}

		logger.Infof("Object not found, create a new one")

		u.SetUnstructuredContent(content)
		u.SetGroupVersionKind(gvk)
		u.SetName(name)
		if owner != nil {
			u.SetOwnerReferences([]metav1.OwnerReference{*owner})
		}

		if err := c.restClient.
			Post().
			Resource(restMapping.Resource.Resource).
			Body(u).
			Do(ctx).
			Into(u); err != nil {
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

	if err := c.restClient.
		Patch(types.MergePatchType).
		Resource(restMapping.Resource.Resource).
		Name(name).
		Body(contentData).
		Do(ctx).
		Into(u); err != nil {
		return nil, errors.Wrap(err, "patching object")
	}

	return u, nil
}

func (c *membershipListener) createOrUpdateStackDependency(ctx context.Context, name string, stackName string, stack *unstructured.Unstructured, gvk schema.GroupVersionKind, content map[string]any) (*unstructured.Unstructured, error) {
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

func NewMembershipListener(restClient *rest.RESTClient, clientInfo ClientInfo, mapper meta.RESTMapper,
	orders MembershipClient, stacksModules InMemoryStacksModules) *membershipListener {
	return &membershipListener{
		restClient:    restClient,
		clientInfo:    clientInfo,
		stacksModules: stacksModules,
		restMapper:    mapper,
		orders:        orders,
		wp:            pond.New(5, 5),
	}
}
