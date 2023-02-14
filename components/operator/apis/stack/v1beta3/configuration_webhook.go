package v1beta3

import (
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var _ webhook.Defaulter = &Configuration{}

// log is for logging in this package.
//
//nolint:unused
var configurationLog = logf.Log.WithName("configuration-resource")

//+kubebuilder:webhook:path=/mutate-stack-formance-com-v1beta3-configuration,mutating=true,failurePolicy=fail,sideEffects=None,groups=stack.formance.com,resources=configurations,verbs=create;update,versions=v1beta3,name=mconfiguration.kb.io,admissionReviewVersions=v1

func (r *Configuration) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
