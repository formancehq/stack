package internal

import (
	"reflect"

	sharedlogging "github.com/formancehq/go-libs/logging"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal/generated"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func convertUnstructured[T client.Object](v any) T {
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(
		v.(*unstructured.Unstructured).Object, t); err != nil {
		panic(err)
	}
	return t
}

func VersionsEventHandler(logger sharedlogging.Logger, membershipClient MembershipClient) cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

			version := convertUnstructured[*v1beta1.Versions](obj)

			logger.Infof("Detect versions '%s' added", version.Name)
			if err := membershipClient.Send(&generated.Message{
				Message: &generated.Message_AddedVersion{
					AddedVersion: &generated.AddedVersion{
						Name:     version.Name,
						Versions: version.Spec,
					},
				},
			}); err != nil {
				logger.Errorf("Unable to send version update: %s", err)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {

			oldVersions := convertUnstructured[*v1beta1.Versions](oldObj)
			newVersions := convertUnstructured[*v1beta1.Versions](newObj)

			if reflect.DeepEqual(oldVersions.Spec, newVersions.Spec) {
				return
			}

			logger.Infof("Detect versions '%s' modified", newVersions.Name)
			if err := membershipClient.Send(&generated.Message{
				Message: &generated.Message_UpdatedVersion{
					UpdatedVersion: &generated.UpdatedVersion{
						Name:     newVersions.Name,
						Versions: newVersions.Spec,
					},
				},
			}); err != nil {
				logger.Errorf("Unable to send version update: %s", err)
			}
		},
		DeleteFunc: func(obj interface{}) {
			version := convertUnstructured[*v1beta1.Versions](obj)

			logger.Infof("Detect versions '%s' as deleted", version.Name)
			if err := membershipClient.Send(&generated.Message{
				Message: &generated.Message_DeletedVersion{
					DeletedVersion: &generated.DeletedVersion{
						Name: version.Name,
					},
				},
			}); err != nil {
				logger.Errorf("Unable to send version update: %s", err)
			}
		},
	}
}
