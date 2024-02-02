package authclients

import (
	"fmt"

	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Create(ctx core.Context, stack *v1beta1.Stack, owner client.Object, objectName string, options ...func(spec *v1beta1.AuthClientSpec)) (*v1beta1.AuthClient, error) {
	authClient, _, err := core.CreateOrUpdate[*v1beta1.AuthClient](ctx, types.NamespacedName{
		Name: core.GetObjectName(stack.Name, objectName),
	},
		func(t *v1beta1.AuthClient) error {
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

			return nil
		},
		core.WithController[*v1beta1.AuthClient](ctx.GetScheme(), owner),
	)
	if err != nil {
		return nil, err
	}
	return authClient, err
}

func WithScopes(scopes ...string) func(client *v1beta1.AuthClientSpec) {
	return func(clientSpec *v1beta1.AuthClientSpec) {
		clientSpec.Scopes = scopes
	}
}

func GetEnvVars(authClient *v1beta1.AuthClient) []v1.EnvVar {
	return []v1.EnvVar{
		core.EnvFromSecret("STACK_CLIENT_ID", fmt.Sprintf("auth-client-%s", authClient.Name), "id"),
		core.EnvFromSecret("STACK_CLIENT_SECRET", fmt.Sprintf("auth-client-%s", authClient.Name), "secret"),
	}
}
