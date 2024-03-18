//nolint:nosnakecase
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/alitto/pond"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
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
				sharedlogging.FromContext(ctx).Infof("Got message from membership: %T", msg.GetMessage())
				switch msg := msg.Message.(type) {
				case *generated.Order_ExistingStack:
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

// convert name to pascal case
func pascalNameModule(modules ...*generated.Service) []string {
	var names []string
	for _, module := range modules {
		moduleName := strings.ToUpper(string(module.Name[0])) + module.Name[1:]
		names = append(names, moduleName)
	}
	return names
}

func ExpectedModules(stack *generated.Stack) []string {
	modules := []string{}
	if stack.StargateConfig != nil && stack.StargateConfig.Enabled {
		modules = append(modules, "Stargate")
	}

	modules = append(modules, pascalNameModule(stack.Services...)...)

	if len(modules) > 0 {
		modules = append(modules, "Gateway")
	}

	return modules
}

func NotExpectedModules(stack *unstructured.Unstructured) []string {
	status, ok := stack.Object["status"].(map[string]any)
	if !ok {
		panic("status not found")
	}

	var installedModules []string
	installedModulesI, ok := status["modules"].([]any)
	if !ok {
		installedModules = []string{}
	} else {
		installedModules = collectionutils.Map(installedModulesI, func(module any) string {
			return module.(string)
		})
	}

	spec, ok := stack.Object["spec"].(map[string]any)
	if !ok {
		panic("spec not found")
	}

	expectedModulesI, ok := spec["expectedModules"].([]any)
	if !ok {
		panic("expectedModules not found")
	}

	expectedModules := collectionutils.Map(expectedModulesI, func(module any) string {
		return module.(string)
	})

	return Reduce(installedModules, func(acc []string, module string) []string {
		if collectionutils.Contains(expectedModules, module) {
			return acc
		}
		return append(acc, module)
	}, []string{})
}

func ExpectedModulesAsGVK(modules []string) []schema.GroupVersionKind {
	return collectionutils.Map(modules, func(module string) schema.GroupVersionKind {
		return v1beta1.GroupVersion.WithKind(module)
	})
}

func (c *membershipListener) syncExistingStack(ctx context.Context, membershipStack *generated.Stack) {

	versions := membershipStack.Versions
	if versions == "" {
		versions = "default"
	}

	additionalLabels := collectionutils.ConvertMap(membershipStack.AdditionalLabels, func(value string) any {
		return value
	})

	expectedStringModules := ExpectedModules(membershipStack)
	expectedModules := ExpectedModulesAsGVK(expectedStringModules)
	stack, err := c.createOrUpdate(ctx, v1beta1.GroupVersion.WithKind("Stack"), membershipStack.ClusterName, membershipStack.ClusterName, nil, map[string]any{
		"spec": map[string]any{
			"versionsFromFile": versions,
			"disabled":         membershipStack.Disabled,
			"enableAudit":      membershipStack.EnableAudit,
			"expectedModules":  expectedStringModules,
		},
	}, additionalLabels)
	if err != nil {
		sharedlogging.FromContext(ctx).Errorf("Unable to create stack cluster side: %s", err)
		return
	}

	for _, gvkModule := range expectedModules {
		if gvkModule.Kind == "Stargate" || gvkModule.Kind == "Auth" || gvkModule.Kind == "Gateway" {
			continue
		}

		if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, gvkModule, map[string]any{}, additionalLabels); err != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to create module %s cluster side: %s", gvkModule.Kind, err)
		}
	}

	if collectionutils.Contains(expectedStringModules, "Auth") {
		if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, v1beta1.GroupVersion.WithKind("Auth"), map[string]any{
			"spec": map[string]any{
				"delegatedOIDCServer": map[string]any{
					"issuer":       membershipStack.AuthConfig.Issuer,
					"clientID":     membershipStack.AuthConfig.ClientId,
					"clientSecret": membershipStack.AuthConfig.ClientSecret,
				},
			},
		}, additionalLabels); err != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to create module Auth cluster side: %s", err)
		}
	} else {
		// If the module is not in the expectedModules list, then it should delete the module
		if err := c.restClient.Delete().
			Resource("Auths").
			Name(stack.GetName()).
			Do(ctx).
			Error(); client.IgnoreNotFound(err) != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to delete module Auth cluster side: %s", err)
		}
	}

	if collectionutils.Contains(expectedStringModules, "Stargate") {
		parts := strings.Split(stack.GetName(), "-")
		if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, v1beta1.GroupVersion.WithKind("Stargate"), map[string]any{
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
		}, additionalLabels); err != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to create module Stargate cluster side: %s", err)
		}
	}

	if collectionutils.Contains(expectedStringModules, "Gateway") {
		if _, err := c.createOrUpdateStackDependency(ctx, stack.GetName(), stack.GetName(), stack, v1beta1.GroupVersion.WithKind("Gateway"), map[string]any{
			"spec": map[string]any{
				"ingress": map[string]any{
					"host":   fmt.Sprintf("%s.%s", stack.GetName(), c.clientInfo.BaseUrl.Host),
					"scheme": c.clientInfo.BaseUrl.Scheme,
				},
			},
		}, additionalLabels); client.IgnoreNotFound(err) != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to create module Stargate cluster side: %s", err)
		}
	}

	syncAuthClients := make([]*unstructured.Unstructured, 0)
	if collectionutils.Contains(expectedStringModules, "Auth") {
		for _, client := range membershipStack.StaticClients {
			authClient, err := c.createOrUpdateStackDependency(ctx, fmt.Sprintf("%s-%s", stack.GetName(), client.Id), stack.GetName(),
				stack, v1beta1.GroupVersion.WithKind("AuthClient"), map[string]any{
					"spec": map[string]any{
						"id":     client.Id,
						"public": client.Public,
					},
				}, additionalLabels)
			if err != nil {
				sharedlogging.FromContext(ctx).Errorf("Unable to create AuthClient cluster side: %s", err)
			}
			syncAuthClients = append(syncAuthClients, authClient)
		}
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

	// Here i need to find all AuthClients that differ from syncAuthClients and delete them
	// If the modules is not in the expectedModules list, then it should delete all AuthClients
	authClientsToDelete := Reduce(authClientList.Items, func(acc []string, item unstructured.Unstructured) []string {
		for _, syncAuthClient := range syncAuthClients {
			if syncAuthClient.GetName() == item.GetName() {
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

	moduleToDelete := NotExpectedModules(stack)
	for _, module := range moduleToDelete {
		moduleAPI := ModuleToPlural(module)
		sharedlogging.FromContext(ctx).Infof("Deleting module %s", module)
		if err := c.restClient.Delete().
			Resource(moduleAPI).
			Name(stack.GetName()).
			Do(ctx).
			Error(); err != nil {
			sharedlogging.FromContext(ctx).Errorf("Unable to delete module %s cluster side: %s", module, err)
		}
	}

	sharedlogging.FromContext(ctx).Infof("Stack %s updated cluster side", stack.GetName())
}

func ModuleToPlural(module string) string {
	switch {
	case module == "Search":
		return "Searches"
	case string(module[len(module)-1:]) == "s":
		return module
	default:
		return module + "s"
	}
}

func Reduce[TYPE any, ACC any](input []TYPE, reducer func(ACC, TYPE) ACC, initial ACC) ACC {
	ret := initial
	for _, i := range input {
		ret = reducer(ret, i)
	}
	return ret
}

func (c *membershipListener) deleteStack(ctx context.Context, stack *generated.DeletedStack) {
	if err := c.restClient.Delete().
		Resource("Stacks").
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

func (c *membershipListener) createOrUpdate(ctx context.Context, gvk schema.GroupVersionKind, name string, stackName string, owner *metav1.OwnerReference, content map[string]any, additionalLabel map[string]any) (*unstructured.Unstructured, error) {

	logger := sharedlogging.WithFields(map[string]any{
		"gvk": gvk,
	})
	logger.Infof("creating object '%s'", name)
	if content["metadata"] == nil {
		content["metadata"] = map[string]any{}
	}

	content["metadata"].(map[string]any)["labels"] = map[string]any{}
	for k, v := range additionalLabel {
		content["metadata"].(map[string]any)["labels"].(map[string]any)["formance.com/"+k] = v
	}
	content["metadata"].(map[string]any)["labels"].(map[string]any)["formance.com/created-by-agent"] = "true"
	content["metadata"].(map[string]any)["labels"].(map[string]any)["formance.com/stack"] = stackName
	content["metadata"].(map[string]any)["name"] = name

	restMapping, err := c.restMapper.RESTMapping(gvk.GroupKind())
	if err != nil {
		panic(err)
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

func (c *membershipListener) createOrUpdateStackDependency(ctx context.Context, name string, stackName string, stack *unstructured.Unstructured, gvk schema.GroupVersionKind, content map[string]any, additionalLabel map[string]any) (*unstructured.Unstructured, error) {
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
		}, content, additionalLabel)
}

func NewMembershipListener(restClient *rest.RESTClient, clientInfo ClientInfo, mapper meta.RESTMapper,
	orders MembershipClient) *membershipListener {
	return &membershipListener{
		restClient: restClient,
		clientInfo: clientInfo,
		restMapper: mapper,
		orders:     orders,
		wp:         pond.New(5, 5),
	}
}
