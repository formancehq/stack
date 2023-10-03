package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	innerGrpc "github.com/formancehq/stack/components/agent/internal/grpc"
	"github.com/formancehq/stack/components/agent/internal/k8s"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	sharedotlptraces "github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	ServiceName = "membership"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
	debugFlag                      = "debug"
	kubeConfigFlag                 = "kube-config"
	serverAddressFlag              = "server-address"
	tlsEnabledFlag                 = "tls-enabled"
	tlsInsecureSkipVerifyFlag      = "tls-insecure-skip-verify"
	tlsCACertificateFlag           = "tls-ca-cert"
	idFlag                         = "id"
	authenticationModeFlag         = "authentication-mode"
	authenticationTokenFlag        = "authentication-token"
	authenticationIssuerFlag       = "authentication-issuer"
	authenticationClientSecretFlag = "authentication-client-secret"
	baseUrlFlag                    = "base-url"
	productionFlag                 = "production"
)

func newK8SConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		sharedlogging.Info("Does not seems to be in cluster, trying to load k8s client from kube config file")
		config, err = clientcmd.BuildConfigFromFlags("", viper.GetString(kubeConfigFlag))
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		serverAddress := viper.GetString(serverAddressFlag)
		if serverAddress == "" {
			return errors.New("missing server address")
		}

		id := viper.GetString(idFlag)
		if id == "" {
			return errors.New("missing id")
		}

		dialOptions := make([]grpc.DialOption, 0)
		var credential credentials.TransportCredentials
		if !viper.GetBool(tlsEnabledFlag) {
			sharedlogging.FromContext(cmd.Context()).Infof("TLS not enabled")
			credential = insecure.NewCredentials()
		} else {
			sharedlogging.FromContext(cmd.Context()).Infof("TLS enabled")
			certPool := x509.NewCertPool()
			if ca := viper.GetString(tlsCACertificateFlag); ca != "" {
				sharedlogging.FromContext(cmd.Context()).Infof("Load server certificate from config")
				if !certPool.AppendCertsFromPEM([]byte(ca)) {
					return fmt.Errorf("failed to add server CA's certificate")
				}
			}

			if viper.GetBool(tlsInsecureSkipVerifyFlag) {
				sharedlogging.FromContext(cmd.Context()).Infof("Disable certificate checks")
			}
			credential = credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: viper.GetBool(tlsInsecureSkipVerifyFlag),
				RootCAs:            certPool,
			})
		}
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credential))

		baseUrlString := viper.GetString(baseUrlFlag)
		if baseUrlString == "" {
			return errors.New("missing base url")
		}
		baseUrl, err := url.Parse(baseUrlString)
		if err != nil {
			return err
		}

		var authenticator innerGrpc.Authenticator
		switch viper.GetString(authenticationModeFlag) {
		case "token":
			token := viper.GetString(authenticationTokenFlag)
			if token == "" {
				return errors.New("missing authentication token")
			}
			authenticator = innerGrpc.TokenAuthenticator(token)
		case "bearer":
			clientSecret := viper.GetString(authenticationClientSecretFlag)
			if clientSecret == "" {
				return errors.New("missing client secret")
			}
			issuer := viper.GetString(authenticationIssuerFlag)
			if issuer == "" {
				return errors.New("missing issuer")
			}

			authenticator = innerGrpc.BearerAuthenticator(issuer, id, clientSecret)
		default:
			return errors.New("authentication mode not specified")
		}

		options := []fx.Option{
			fx.Provide(newK8SConfig),
			fx.NopLogger,
			k8s.NewModule(),
			innerGrpc.NewModule(serverAddress, authenticator, innerGrpc.ClientInfo{
				ID:         id,
				BaseUrl:    baseUrl,
				Production: viper.GetBool(productionFlag),
				Version:    Version,
			}, dialOptions...),
			sharedotlptraces.CLITracesModule(viper.GetViper()),
		}

		return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
	},
}

func exitWithCode(code int, v ...any) {
	fmt.Fprintln(os.Stdout, v...)
	os.Exit(code)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithCode(1, err)
	}
}

func init() {
	var kubeConfigFilePath string
	if home := homedir.HomeDir(); home != "" {
		kubeConfigFilePath = filepath.Join(home, ".kube", "config")
	}
	rootCmd.Flags().String(kubeConfigFlag, kubeConfigFilePath, "")
	rootCmd.Flags().String(serverAddressFlag, "localhost:8081", "")
	rootCmd.Flags().Bool(tlsEnabledFlag, false, "")
	rootCmd.Flags().Bool(tlsInsecureSkipVerifyFlag, false, "")
	rootCmd.Flags().String(tlsCACertificateFlag, "", "")
	rootCmd.Flags().String(idFlag, "", "")
	rootCmd.Flags().String(authenticationModeFlag, "", "")
	rootCmd.Flags().String(authenticationTokenFlag, "", "")
	rootCmd.Flags().String(authenticationClientSecretFlag, "", "")
	rootCmd.Flags().String(authenticationIssuerFlag, "", "")
	rootCmd.Flags().String(baseUrlFlag, "", "")
	rootCmd.Flags().Bool(productionFlag, false, "Is a production agent")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP(debugFlag, "d", false, "Debug mode")
}
