package common

import (
	"github.com/formancehq/operator/internal/collectionutils"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MatchingLabels(label string) func(deployment *v1.Deployment) {
	return func(t *v1.Deployment) {
		matchLabels := collectionutils.CreateMap("app.kubernetes.io/name", label)
		t.Spec.Selector = &v12.LabelSelector{
			MatchLabels: matchLabels,
		}
		t.Spec.Template.ObjectMeta = v12.ObjectMeta{
			Labels: matchLabels,
		}
	}
}
