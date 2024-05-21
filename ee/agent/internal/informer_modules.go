package internal

import (
	"encoding/json"
	"reflect"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"
)

func Restrict[T any](obj map[string]interface{}) (map[string]interface{}, error) {
	if len(obj) == 0 {
		return nil, errors.New("obj is empty")
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	res := new(T)
	err = json.Unmarshal(jsonBytes, &res)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal json in"+reflect.TypeOf(res).String())
	}

	filtered, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	var tmp map[string]interface{}
	if err := json.Unmarshal(filtered, &tmp); err != nil {
		return nil, err
	}

	return tmp, nil
}

func getStatus(unstructuredModule *unstructured.Unstructured) (*structpb.Struct, error) {
	status, found, err := unstructured.NestedMap(unstructuredModule.Object, "status")
	if err != nil {
		return nil, errors.Wrap(err, "unable to get status from unstructured")
	}

	if !found {
		return nil, nil
	}

	status, err = Restrict[v1beta1.Status](status)
	if err != nil {
		return nil, errors.Wrap(err, "unable to restrict status according to v1beta1.StatusWithConditions")
	}

	protoStatus, err := structpb.NewStruct(status)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert status to proto struct")
	}
	return protoStatus, nil

}

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

func NewModuleEventHandler(logger logging.Logger, membershipClient MembershipClient) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			unstructuredModule := obj.(*unstructured.Unstructured)
			status, err := getStatus(unstructuredModule)
			if err != nil {
				logger.Debug("unstructuredModule", unstructuredModule)
				logger.Errorf("Unable to generate message module added: %s", err)
				return
			}

			if status == nil {
				return
			}

			message := fromUnstructuredToModuleStatusChanged(unstructuredModule, status)
			if err := membershipClient.Send(message); err != nil {
				logger.Errorf("Unable to send message module added: %s", err)
				return
			}

			logger.Infof("Detect module '%s' added", unstructuredModule.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldVersions := oldObj.(*unstructured.Unstructured)
			newVersions := newObj.(*unstructured.Unstructured)

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
			if err := membershipClient.Send(message); err != nil {
				logger.Errorf("Unable to send message module update: %s", err)
			}
			logger.Infof("Detect module '%s' updated", newVersions.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			module := obj.(*unstructured.Unstructured)

			if err := membershipClient.Send(&generated.Message{
				Message: &generated.Message_ModuleDeleted{
					ModuleDeleted: &generated.ModuleDeleted{
						ClusterName: module.GetName(),
						Vk: &generated.VersionKind{
							Version: module.GetObjectKind().GroupVersionKind().Version,
							Kind:    module.GetObjectKind().GroupVersionKind().Kind,
						},
					},
				},
			}); err != nil {
				logger.Errorf("Unable to send message module delete: %s", err)
			}

			logger.Infof("Detect module '%s' deleted", module.GetName())
		},
	}
}
