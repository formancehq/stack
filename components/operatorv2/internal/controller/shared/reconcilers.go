package shared

import ctrl "sigs.k8s.io/controller-runtime"

func SetupReconcilers(mgr ctrl.Manager, reconcilers ...interface {
	SetupWithManager(mgr ctrl.Manager) error
}) error {
	for _, reconciler := range reconcilers {
		if err := reconciler.SetupWithManager(mgr); err != nil {
			return err
		}
	}
	return nil
}
