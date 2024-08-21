package registries_test

import (
	"testing"

	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTranslateImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	type testCase struct {
		url string
	}
	// ghcr.io/<organization>/<repository>:<version>
	// public.ecr.aws/<organization>/jeffail/benthos:<version>
	// docker.io/<organization|user>/<image>:<version>
	// version: "v2.0.0-rc.35-scratch@sha256:4a29620448a90f3ae50d2e375c993b86ef141ead4b6ac1edd1674e9ff6b933f8"
	// docker.io/<organization|user>/<image>:v2.0.0-rc.35-scratch@sha256:4a29620448a90f3ae50d2e375c993b86ef141ead4b6ac1edd1674e9ff6b933f8
	testCases := []testCase{
		{
			url: "ghcr.io/formancehq/stack:latest",
		},
		{
			url: "public.ecr.aws/formance-internal/jeffail/benthos:latest",
		},
		{
			url: "docker.io/natsio/nats-box:latest",
		},
		{
			url: "docker.io/caddy/caddy:2.7.6-alpine",
		},
		{
			url: "ghcr.io/formancehq/operator-utils:v2.0.0-rc.35-scratch@sha256:4a29620448a90f3ae50d2e375c993b86ef141ead4b6ac1edd1674e9ff6b933f8",
		},
	}

	for _, tc := range testCases {
		mockImageSettingsOverrider := registries.NewMockImageSettingsOverrider(ctrl)
		mockImageSettingsOverrider.EXPECT().OverrideWithSetting(gomock.Any(), gomock.Any()).Return(nil)

		fullImageString, err := registries.TranslateImage("stackName", mockImageSettingsOverrider, tc.url)
		require.NoError(t, err)
		require.Equal(t, tc.url, fullImageString)
	}
}
