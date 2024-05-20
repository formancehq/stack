package internal

import (
	"reflect"

	"github.com/formancehq/stack/components/agent/internal/generated"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

func fromUnstructuredToModuleStatusChanged(unstructuredModule *unstructured.Unstructured, status *structpb.Struct) *generated.Message {
	return &generated.Message{
		Message: &generated.Message_ModuleStatusChanged{
			ModuleStatusChanged: &generated.ModuleStatusChanged{
				ClusterName: unstructuredModule.GetName(),
				Vk: &generated.VersionKind{
					Version: unstructuredModule.GetObjectKind().GroupVersionKind().Version,
					Kind:    unstructuredModule.GetObjectKind().GroupVersionKind().Kind,
				},
				Status: status,
			},
		},
	}
}

type ModuleEventHandler struct {
	logger logging.Logger
	client MembershipClient
}

func (h *ModuleEventHandler) AddFunc(obj interface{}) {
	unstructuredModule := obj.(*unstructured.Unstructured)
	status, err := getStatus(unstructuredModule)
	if err != nil {
		h.logger.Debug("unstructuredModule", unstructuredModule)
		h.logger.Errorf("Unable to generate message module added: %s", err)
		return
	}

	if status == nil {
		return
	}

	message := fromUnstructuredToModuleStatusChanged(unstructuredModule, status)
	if err := h.client.Send(message); err != nil {
		h.logger.Errorf("Unable to send message module added: %s", err)
		return
	}

	h.logger.Infof("Detect module '%s' added", unstructuredModule.GetName())
}

func (h *ModuleEventHandler) UpdateFunc(oldObj, newObj any) {
	oldVersions := oldObj.(*unstructured.Unstructured)
	newVersions := newObj.(*unstructured.Unstructured)

	oldStatus, err := getStatus(oldVersions)
	if err != nil {
		h.logger.Errorf("Unable to get status from old versions: %s", err)
	}
	newStatus, err := getStatus(newVersions)
	if err != nil {
		h.logger.Errorf("Unable to get status from new versions: %s", err)
		return
	}

	if newStatus == nil || reflect.DeepEqual(oldStatus, newStatus) {
		return
	}

	message := fromUnstructuredToModuleStatusChanged(newVersions, newStatus)
	if err := h.client.Send(message); err != nil {
		h.logger.Errorf("Unable to send message module update: %s", err)
	}
	h.logger.Infof("Detect module '%s' updated", newVersions.GetName())
}

func (h *ModuleEventHandler) DeleteFunc(obj interface{}) {
	unstructuredModule := obj.(*unstructured.Unstructured)
	if err := h.client.Send(&generated.Message{
		Message: &generated.Message_ModuleDeleted{
			ModuleDeleted: &generated.ModuleDeleted{
				ClusterName: unstructuredModule.GetName(),
				Vk: &generated.VersionKind{
					Version: unstructuredModule.GetObjectKind().GroupVersionKind().Version,
					Kind:    unstructuredModule.GetObjectKind().GroupVersionKind().Kind,
				},
			},
		},
	}); err != nil {
		h.logger.Errorf("Unable to send message module deleted: %s", err)
		return
	}
	h.logger.Infof("Detect module '%s' deleted", unstructuredModule.GetName())
}

func NewModuleEventHandler(logger logging.Logger, membershipClient MembershipClient) cache.ResourceEventHandlerFuncs {
	moduleEventHandler := &ModuleEventHandler{
		logger: logger,
		client: membershipClient,
	}

	return cache.ResourceEventHandlerFuncs{
		AddFunc:    moduleEventHandler.AddFunc,
		UpdateFunc: moduleEventHandler.UpdateFunc,
		DeleteFunc: moduleEventHandler.DeleteFunc,
	}
}
