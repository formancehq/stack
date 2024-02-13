package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
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
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	ServiceName = "agent"
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

func init() {
	if err := v1beta1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
}

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return bindFlagsToViper(cmd)
	},
	RunE: runAgent,
}

func exitWithCode(code int, v ...any) {
	fmt.Fprintln(os.Stdout, v...)
	os.Exit(code)
}

func runAgent(cmd *cobra.Command, args []string) error {
	serverAddress := viper.GetString(serverAddressFlag)
	if serverAddress == "" {
		return errors.New("missing server address")
	}

	agentID := viper.GetString(idFlag)
	if agentID == "" {
		return errors.New("missing id")
	}

	credentials, err := createGRPCTransportCredentials(cmd.Context())
	if err != nil {
		return err
	}

	dialOptions := make([]grpc.DialOption, 0)
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials))

	baseUrlString := viper.GetString(baseUrlFlag)
	if baseUrlString == "" {
		return errors.New("missing base url")
	}

	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return err
	}

	authenticator, err := createAuthenticator(agentID)
	if err != nil {
		return err
	}

	options := []fx.Option{
		fx.Provide(newK8SConfig),
		fx.NopLogger,
		internal.NewModule(serverAddress, authenticator, internal.ClientInfo{
			ID:         agentID,
			BaseUrl:    baseUrl,
			Production: viper.GetBool(productionFlag),
			Version:    Version,
		}, dialOptions...),
		sharedotlptraces.CLITracesModule(),
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
}

func newK8SConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		sharedlogging.Info("Does not seems to be in cluster, trying to load k8s client from kube config file")
		config, err = clientcmd.BuildConfigFromFlags("", viper.GetString(kubeConfigFlag))
		if err != nil {
			return nil, err
		}
	}

	config.GroupVersion = &v1beta1.GroupVersion
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.APIPath = "/apis"

	return config, nil
}

func createAuthenticator(agentID string) (internal.Authenticator, error) {
	var authenticator internal.Authenticator
	switch viper.GetString(authenticationModeFlag) {
	case "token":
		token := viper.GetString(authenticationTokenFlag)
		if token == "" {
			return nil, errors.New("missing authentication token")
		}
		authenticator = internal.TokenAuthenticator(token)
	case "bearer":
		clientSecret := viper.GetString(authenticationClientSecretFlag)
		if clientSecret == "" {
			return nil, errors.New("missing client secret")
		}
		issuer := viper.GetString(authenticationIssuerFlag)
		if issuer == "" {
			return nil, errors.New("missing issuer")
		}

		authenticator = internal.BearerAuthenticator(issuer, agentID, clientSecret)
	default:
		return nil, errors.New("authentication mode not specified")
	}
	return authenticator, nil
}

func createGRPCTransportCredentials(ctx context.Context) (credentials.TransportCredentials, error) {
	var credential credentials.TransportCredentials
	if !viper.GetBool(tlsEnabledFlag) {
		sharedlogging.FromContext(ctx).Infof("TLS not enabled")
		credential = insecure.NewCredentials()
	} else {
		sharedlogging.FromContext(ctx).Infof("TLS enabled")
		certPool := x509.NewCertPool()
		if ca := viper.GetString(tlsCACertificateFlag); ca != "" {
			sharedlogging.FromContext(ctx).Infof("Load server certificate from config")
			if !certPool.AppendCertsFromPEM([]byte(ca)) {
				return nil, fmt.Errorf("failed to add server CA's certificate")
			}
		}

		if viper.GetBool(tlsInsecureSkipVerifyFlag) {
			sharedlogging.FromContext(ctx).Infof("Disable certificate checks")
		}
		credential = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: viper.GetBool(tlsInsecureSkipVerifyFlag),
			RootCAs:            certPool,
		})
	}
	return credential, nil
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
