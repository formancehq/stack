/*
Copyright 2023.

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
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/controller"
	"github.com/formancehq/operator/v2/internal/controller/shared"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	formancev1beta1 "github.com/formancehq/operator/v2/api/v1beta1"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(formancev1beta1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
		probeAddr            string
		region               string
		env                  string
	)
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&region, "region", "eu-west-1", "The cloud region in use for the operator")
	flag.StringVar(&env, "env", "staging", "The current environment in use for the operator")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "6e1085e1.com",
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

	platform := shared.Platform{
		Region:      region,
		Environment: env,
	}

	if err := reconcilers.Setup(mgr,
		reconcilers.New[*v1beta1.Stack](mgr.GetClient(), mgr.GetScheme(), controller.ForStack(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Topic](mgr.GetClient(), mgr.GetScheme(), controller.ForTopic(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.TopicQuery](mgr.GetClient(), mgr.GetScheme(), controller.ForTopicQuery(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Ledger](mgr.GetClient(), mgr.GetScheme(), controller.ForLedger(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.HTTPAPI](mgr.GetClient(), mgr.GetScheme(), controller.ForHTTPAPI(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Gateway](mgr.GetClient(), mgr.GetScheme(), controller.ForGateway(mgr.GetClient(), mgr.GetScheme(), platform)),
		reconcilers.New[*v1beta1.Auth](mgr.GetClient(), mgr.GetScheme(), controller.ForAuth(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Database](mgr.GetClient(), mgr.GetScheme(), controller.ForDatabase(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.AuthClient](mgr.GetClient(), mgr.GetScheme(), controller.ForAuthClient(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Wallets](mgr.GetClient(), mgr.GetScheme(), controller.ForWallets(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Orchestration](mgr.GetClient(), mgr.GetScheme(), controller.ForOrchestration(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Payments](mgr.GetClient(), mgr.GetScheme(), controller.ForPayments(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Reconciliation](mgr.GetClient(), mgr.GetScheme(), controller.ForReconciliation(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Webhooks](mgr.GetClient(), mgr.GetScheme(), controller.ForWebhooks(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Search](mgr.GetClient(), mgr.GetScheme(), controller.ForSearch(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.StreamProcessor](mgr.GetClient(), mgr.GetScheme(), controller.ForStreamProcessor(mgr.GetClient(), mgr.GetScheme())),
		reconcilers.New[*v1beta1.Stream](mgr.GetClient(), mgr.GetScheme(), controller.ForStream(mgr.GetClient(), mgr.GetScheme())),
	); err != nil {
		setupLog.Error(err, "unable to create controllers")
		os.Exit(1)
	}

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
