package internal

import (
	"context"
	"path/filepath"
	osRuntime "runtime"
	"testing"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type testConfig struct {
	restConfig *rest.Config
	mapper     meta.RESTMapper
	client     *rest.RESTClient
}

func test(t *testing.T, fn func(context.Context, *testConfig)) {
	_, filename, _, _ := osRuntime.Caller(0)
	apiServer := envtest.APIServer{}
	apiServer.Configure().
		Set("service-cluster-ip-range", "10.0.0.0/20")

	require.NoError(t, v1beta1.AddToScheme(scheme.Scheme))
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join(filepath.Dir(filename), "..", "..", "..", "components", "operator",
				"config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: true,
		ControlPlane: envtest.ControlPlane{
			APIServer: &apiServer,
		},
		Scheme: scheme.Scheme,
	}

	restConfig, err := testEnv.Start()

	require.NoError(t, err)

	restConfig.GroupVersion = &v1beta1.GroupVersion
	restConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	restConfig.APIPath = "/apis"

	k8sClient, err := rest.RESTClientFor(restConfig)
	require.NoError(t, err)

	mapper, err := CreateRestMapper(restConfig)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			require.NoError(t, testEnv.Stop())
		},
	)
	fn(logging.TestingContext(), &testConfig{
		restConfig: restConfig,
		mapper:     mapper,
		client:     k8sClient,
	})
}
func TestDeleteModule(t *testing.T) {

	type testCase struct {
		name       string
		withLabels bool
	}

	testCases := []testCase{
		{
			name:       "with labels",
			withLabels: true,
		},
		{
			name:       "without labels",
			withLabels: false,
		},
	}
	test(t, func(ctx context.Context, testConfig *testConfig) {
		t.Parallel()

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				stackName := uuid.NewString()
				recon := v1beta1.Reconciliation{
					ObjectMeta: v1.ObjectMeta{
						Name: uuid.NewString(),
					},
				}
				if tc.withLabels {
					recon.Labels = map[string]string{
						"formance.com/created-by-agent": "true",
						"formance.com/stack":            stackName,
					}
				}

				require.NoError(t, testConfig.client.Post().Resource("Reconciliations").Body(&recon).Do(ctx).Error())
				orders := NewMembershipClientMock()

				membershipListener := NewMembershipListener(
					testConfig.client, ClientInfo{}, testConfig.mapper, orders, nil)

				if tc.withLabels {
					gvk := v1beta1.GroupVersion.WithKind("Reconciliation")
					resource, err := testConfig.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
					require.NoError(t, err)

					require.NoError(t, membershipListener.deleteModule(ctx, resource.Resource, stackName))
				}

				if !tc.withLabels {
					require.NoError(t, testConfig.client.Get().Resource("Reconciliations").Name(recon.Name).Do(ctx).Error())
				}
			})
		}
	})
}

func TestRetrieveModuleList(t *testing.T) {
	t.Parallel()
	test(t, func(ctx context.Context, testConfig *testConfig) {
		modules, eeModules, err := retrieveModuleList(ctx, testConfig.restConfig)
		require.NoError(t, err)
		require.NotEmpty(t, modules)
		require.NotEmpty(t, eeModules)
		for _, module := range eeModules {
			require.Contains(t, modules, module)
		}
	})
}
