/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	authcomponentsv1beta1 "github.com/formancehq/operator/apis/auth.components/v1beta1"
	componentsv1beta1 "github.com/formancehq/operator/apis/components/v1beta1"
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
	"github.com/formancehq/operator/internal/controllers/stack"
	traefik "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	authcomponentsv1beta2 "github.com/formancehq/operator/apis/auth.components/v1beta2"
	benthoscomponentsv1beta2 "github.com/formancehq/operator/apis/benthos.components/v1beta2"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	stackv1beta1 "github.com/formancehq/operator/apis/stack/v1beta1"
	stackv1beta2 "github.com/formancehq/operator/apis/stack/v1beta2"
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(authcomponentsv1beta1.AddToScheme(scheme))
	utilruntime.Must(authcomponentsv1beta2.AddToScheme(scheme))
	utilruntime.Must(traefik.AddToScheme(scheme))
	utilruntime.Must(benthoscomponentsv1beta2.AddToScheme(scheme))
	utilruntime.Must(stackv1beta1.AddToScheme(scheme))
	utilruntime.Must(stackv1beta2.AddToScheme(scheme))
	utilruntime.Must(stackv1beta3.AddToScheme(scheme))
	utilruntime.Must(componentsv1beta1.AddToScheme(scheme))
	utilruntime.Must(componentsv1beta2.AddToScheme(scheme))

	//+kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
		probeAddr            string
		issuerRefName        string
		issuerRefKind        string
		disableWebhooks      bool
	)
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&issuerRefName, "issuer-ref-name", "", "")
	flag.StringVar(&issuerRefKind, "issuer-ref-kind", "ClusterIssuer", "")
	flag.BoolVar(&disableWebhooks, "disable-webhooks", false, "Disable wehooks")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "68fe8eef.formance.com",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if !disableWebhooks {
		if err := (&stackv1beta3.Stack{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to install stack webhook", "controller", "Stack")
			os.Exit(1)
		}
		if err := (&stackv1beta3.Configuration{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to install stack configuration", "controller", "Configuration")
			os.Exit(1)
		}
	}

	stackReconciler := stack.NewReconciler(mgr.GetClient(), mgr.GetScheme())
	if err = stackReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Stack")
		os.Exit(1)
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
