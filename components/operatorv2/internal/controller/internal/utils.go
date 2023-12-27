package internal

import (
	"context"
	"github.com/google/uuid"
	"reflect"

	"github.com/formancehq/operator/v2/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateDatabase(ctx context.Context, client client.Client, stack *v1beta1.Stack, name string) (*v1beta1.Database, error) {
	database, _, err := CreateOrUpdate[*v1beta1.Database](ctx, client, types.NamespacedName{
		Name: GetObjectName(stack.Name, name),
	}, func(t *v1beta1.Database) {
		t.Spec.Stack = stack.Name
		t.Spec.Service = name
	})
	if !database.Status.Ready {
		return nil, ErrPending
	}
	return database, err
}

func CreateHTTPAPI(ctx context.Context, _client client.Client, scheme *runtime.Scheme,
	stack *v1beta1.Stack, owner client.Object, objectName string, options ...func(spec *v1beta1.HTTPAPISpec)) error {
	_, _, err := CreateOrUpdate[*v1beta1.HTTPAPI](ctx, _client, types.NamespacedName{
		Name: GetObjectName(stack.Name, objectName),
	},
		func(t *v1beta1.HTTPAPI) {
			t.Spec = v1beta1.HTTPAPISpec{
				Stack:              stack.Name,
				PortName:           "http",
				HasVersionEndpoint: true,
				Name:               objectName,
			}
			for _, option := range options {
				option(&t.Spec)
			}
		},
		WithController[*v1beta1.HTTPAPI](scheme, owner),
	)
	return err
}

func CreateAuthClient(ctx context.Context, _client client.Client, scheme *runtime.Scheme,
	stack *v1beta1.Stack, owner client.Object, objectName string, options ...func(spec *v1beta1.AuthClientSpec)) (*v1beta1.AuthClient, error) {
	authClient, _, err := CreateOrUpdate[*v1beta1.AuthClient](ctx, _client, types.NamespacedName{
		Name: GetObjectName(stack.Name, objectName),
	},
		func(t *v1beta1.AuthClient) {
			t.Spec.Stack = stack.Name
			if t.Spec.ID == "" {
				t.Spec.ID = uuid.NewString()
			}
			if t.Spec.Secret == "" {
				t.Spec.Secret = uuid.NewString()
			}
			for _, option := range options {
				option(&t.Spec)
			}
		},
		WithController[*v1beta1.AuthClient](scheme, owner),
	)
	if err != nil {
		return nil, err
	}
	return authClient, err
}

func TopicExists(ctx context.Context, _client client.Client, stack *v1beta1.Stack, name string) (bool, error) {
	topicList := &v1beta1.TopicList{}
	if err := _client.List(ctx, topicList, client.MatchingFields{
		".spec.service": name,
		".spec.stack":   stack.Name,
	}); err != nil {
		return false, err
	}
	return len(topicList.Items) > 0, nil
}

func DeleteIfExists[T client.Object](ctx context.Context, _client client.Client, name types.NamespacedName) error {
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	if err := _client.Get(ctx, name, t); err != nil {
		if client.IgnoreNotFound(err) == nil {
			return nil
		}
		return err
	}
	return _client.Delete(ctx, t)
}
