package internal

import (
	"reflect"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/stack/components/agent/internal/generated"
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

	logger := h.logger.WithField("func", "Add").WithField("module", unstructuredModule.GetName())

	status, err := getStatus(unstructuredModule)
	if err != nil {
		logger.Errorf("Unable to generate message module added: %s", err)
		return
	}

	if status == nil {
		return
	}

	message := fromUnstructuredToModuleStatusChanged(unstructuredModule, status)
	if err := h.client.Send(message); err != nil {
		logger.Errorf("Unable to send message module added: %s", err)
		return
	}

	logger.Infof("Detect module '%s' added", unstructuredModule.GetName())
}

func (h *ModuleEventHandler) UpdateFunc(oldObj, newObj any) {

	oldVersions := oldObj.(*unstructured.Unstructured)
	newVersions := newObj.(*unstructured.Unstructured)

	logger := h.logger.WithField("func", "Update").WithField("module", newVersions.GetName())

	oldStatus, err := getStatus(oldVersions)
	if err != nil {
		logger.Errorf("Unable to get status from old versions: %s", err)
	}
	newStatus, err := getStatus(newVersions)
	if err != nil {
		logger.Errorf("Unable to get status from new versions: %s", err)
		return
	}

	if newStatus == nil || reflect.DeepEqual(oldStatus, newStatus) {
		return
	}

	message := fromUnstructuredToModuleStatusChanged(newVersions, newStatus)
	if err := h.client.Send(message); err != nil {
		logger.Errorf("Unable to send message module update: %s", err)
		return
	}
	logger.Infof("Detect module '%s' updated", newVersions.GetName())
}

func (h *ModuleEventHandler) DeleteFunc(obj interface{}) {

	unstructuredModule := obj.(*unstructured.Unstructured)
	logger := h.logger.WithField("func", "Delete").WithField("module", unstructuredModule.GetName())

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
		logger.Errorf("Unable to send message module deleted: %s", err)
		return
	}
	logger.Infof("Detect module '%s' deleted", unstructuredModule.GetName())
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
