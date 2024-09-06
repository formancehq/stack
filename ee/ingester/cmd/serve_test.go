package cmd

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/stretchr/testify/require"
)

func TestEvalServerConfiguration(t *testing.T) {
	t.Parallel()

	type test struct {
		name string
		// init return args to pass to the cobra command, and a string, corresponding to an error that can be returned by the command
		init func(t *testing.T) ([]string, string)
	}

	for _, test := range []test{
		{
			name: "no args",
			init: func(t *testing.T) ([]string, string) {
				return []string{"serve"}, "no stack specified"
			},
		},
		{
			name: "with stack url specified",
			init: func(t *testing.T) ([]string, string) {
				return []string{
					"serve",
					"--" + StackFlag, "test",
					"--" + ModuleURLTplFlag, "http://localhost",
				}, "parsing postgres connection flags: missing postgres uri"
			},
		},
		{
			name: "with postgres uri specified",
			init: func(t *testing.T) ([]string, string) {
				return []string{
					"serve",
					"--" + StackFlag, "test",
					"--" + ModuleURLTplFlag, "http://localhost",
					"--" + bunconnect.PostgresURIFlag, "postgres://localhost:5432",
				}, ""
			},
		},
		{
			name: "with client id specified without secret",
			init: func(t *testing.T) ([]string, string) {
				return []string{
					"serve",
					"--" + StackFlag, "test",
					"--" + ModuleURLTplFlag, "http://localhost",
					"--" + bunconnect.PostgresURIFlag, "postgres://localhost:5432",
					"--" + stackClientIDFlag, "test",
				}, "no stack client secret specified"
			},
		},
		{
			name: "with client id and secret",
			init: func(t *testing.T) ([]string, string) {
				return []string{
					"serve",
					"--" + StackFlag, "test",
					"--" + bunconnect.PostgresURIFlag, "postgres://localhost:5432",
					"--" + stackClientIDFlag, "test",
					"--" + stackClientSecretFlag, "test",
					"--" + ModuleURLTplFlag, "http://localhost",
				}, "no stack issuer specified"
			},
		},
		{
			name: "with missing module url template",
			init: func(t *testing.T) ([]string, string) {
				return []string{
					"serve",
					"--" + StackFlag, "test",
					"--" + bunconnect.PostgresURIFlag, "postgres://localhost:5432",
					"--" + stackClientIDFlag, "test",
					"--" + stackClientSecretFlag, "test",
				}, "no module url template specified"
			},
		},
		{
			name: "with invalid module url template",
			init: func(t *testing.T) ([]string, string) {
				return []string{
					"serve",
					"--" + StackFlag, "test",
					"--" + bunconnect.PostgresURIFlag, "postgres://localhost:5432",
					"--" + stackClientIDFlag, "test",
					"--" + stackClientSecretFlag, "test",
					"--" + ModuleURLTplFlag, "http://{{",
				}, "failed to parse module url template: template: :1: unclosed action"
			},
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			args, expectError := test.init(t)
			args = append(args, "--"+BindFlag, ":0")
			if testing.Verbose() {
				args = append(args, "--"+service.DebugFlag)
			}
			serveCmd := NewServeCommand()
			require.NoError(t, serveCmd.ParseFlags(args))

			_, err := evalServeConfiguration(serveCmd)
			if expectError != "" {
				require.NotNil(t, err)
				require.Equal(t, expectError, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
