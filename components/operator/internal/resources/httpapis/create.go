package httpapis

import (
	v1beta1 "github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

type option func(spec *v1beta1.HTTPAPI)

var defaultOptions = []option{
	WithRules(RuleSecured()),
}

func Create(ctx core.Context, owner core.Module, options ...option) error {
	objectName := strings.ToLower(owner.GetObjectKind().GroupVersionKind().Kind)
	_, _, err := core.CreateOrUpdate[*v1beta1.HTTPAPI](ctx, types.NamespacedName{
		Name: core.GetObjectName(owner.GetStack(), core.GetModuleName(owner)),
	},
		func(t *v1beta1.HTTPAPI) {
			t.Spec = v1beta1.HTTPAPISpec{
				StackDependency: v1beta1.StackDependency{
					Stack: owner.GetStack(),
				},
				Name: objectName,
			}
			for _, option := range append(defaultOptions, options...) {
				option(t)
			}
		},
		core.WithController[*v1beta1.HTTPAPI](ctx.GetScheme(), owner),
	)
	return err
}

func WithRules(rules ...v1beta1.HTTPAPIRule) func(httpapi *v1beta1.HTTPAPI) {
	return func(httpapi *v1beta1.HTTPAPI) {
		httpapi.Spec.Rules = rules
	}
}

func WithServiceConfiguration(serviceConfiguration *v1beta1.ServiceConfiguration) func(httpapi *v1beta1.HTTPAPI) {
	return func(httpapi *v1beta1.HTTPAPI) {
		httpapi.Spec.Service = serviceConfiguration
	}
}

func RuleSecured() v1beta1.HTTPAPIRule {
	return v1beta1.HTTPAPIRule{}
}

func RuleUnsecured() v1beta1.HTTPAPIRule {
	return v1beta1.HTTPAPIRule{
		Secured: true,
	}
}
