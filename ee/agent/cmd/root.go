package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/stack/components/agent/internal"
	"github.com/formancehq/stack/libs/go-libs/licence"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/homedir"
)

var (
	ServiceName = "agent"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

const (
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
	resyncPeriodFlag               = "resync-period"
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	RunE:         runAgent,
}

func init() {
	if err := v1beta1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}

	var kubeConfigFilePath string
	if home := homedir.HomeDir(); home != "" {
		kubeConfigFilePath = filepath.Join(home, ".kube", "config")
	}

	service.AddFlags(rootCmd.PersistentFlags())
	otlptraces.AddFlags(rootCmd.PersistentFlags())
	licence.AddFlags(rootCmd.PersistentFlags())

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
	rootCmd.Flags().Duration(resyncPeriodFlag, 5*time.Minute, "Resync period of K8S resources")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() {
	service.Execute(rootCmd)
}

func runAgent(cmd *cobra.Command, _ []string) error {
	serverAddress, _ := cmd.Flags().GetString(serverAddressFlag)
	if serverAddress == "" {
		return errors.New("missing server address")
	}

	agentID, _ := cmd.Flags().GetString(idFlag)
	if agentID == "" {
		return errors.New("missing id")
	}

	credentials, err := createGRPCTransportCredentials(cmd)
	if err != nil {
		return err
	}

	dialOptions := make([]grpc.DialOption, 0)
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials))

	baseUrlString, _ := cmd.Flags().GetString(baseUrlFlag)
	if baseUrlString == "" {
		return errors.New("missing base url")
	}

	baseUrl, err := url.Parse(baseUrlString)
	if err != nil {
		return err
	}

	authenticator, err := createAuthenticator(cmd)
	if err != nil {
		return err
	}

	kubeConfig, _ := cmd.Flags().GetString(kubeConfigFlag)

	restConfig, err := internal.NewK8SConfig(kubeConfig)
	if err != nil {
		return err
	}

	isProduction, _ := cmd.Flags().GetBool(productionFlag)
	resyncPeriod, _ := cmd.Flags().GetDuration(resyncPeriodFlag)

	options := []fx.Option{
		fx.Supply(restConfig),
		fx.NopLogger,
		internal.NewModule(serverAddress, authenticator, internal.ClientInfo{
			ID:         agentID,
			BaseUrl:    baseUrl,
			Production: isProduction,
			Version:    Version,
		}, resyncPeriod, dialOptions...),
		otlptraces.FXModuleFromFlags(cmd),
		licence.FXModuleFromFlags(cmd, ServiceName),
	}

	return service.New(cmd.OutOrStdout(), options...).Run(cmd)
}

func createAuthenticator(cmd *cobra.Command) (internal.Authenticator, error) {
	var authenticator internal.Authenticator
	authenticationMode, _ := cmd.Flags().GetString(authenticationModeFlag)
	agentID, _ := cmd.Flags().GetString(idFlag)

	switch authenticationMode {
	case "token":

		token, _ := cmd.Flags().GetString(authenticationTokenFlag)
		if token == "" {
			return nil, errors.New("missing authentication token")
		}
		authenticator = internal.TokenAuthenticator(token)
	case "bearer":
		clientSecret, _ := cmd.Flags().GetString(authenticationClientSecretFlag)
		if clientSecret == "" {
			return nil, errors.New("missing client secret")
		}
		issuer, _ := cmd.Flags().GetString(authenticationIssuerFlag)
		if issuer == "" {
			return nil, errors.New("missing issuer")
		}

		authenticator = internal.BearerAuthenticator(issuer, agentID, clientSecret)
	default:
		return nil, errors.New("authentication mode not specified")
	}
	return authenticator, nil
}

func createGRPCTransportCredentials(cmd *cobra.Command) (credentials.TransportCredentials, error) {
	var credential credentials.TransportCredentials
	tlsEnabled, _ := cmd.Flags().GetBool(tlsEnabledFlag)
	if !tlsEnabled {
		sharedlogging.FromContext(cmd.Context()).Infof("TLS not enabled")
		credential = insecure.NewCredentials()
	} else {
		sharedlogging.FromContext(cmd.Context()).Infof("TLS enabled")
		certPool := x509.NewCertPool()
		ca, _ := cmd.Flags().GetString(tlsCACertificateFlag)
		if ca != "" {
			sharedlogging.FromContext(cmd.Context()).Infof("Load server certificate from config")
			if !certPool.AppendCertsFromPEM([]byte(ca)) {
				return nil, fmt.Errorf("failed to add server CA's certificate")
			}
		}

		tlsInsecure, _ := cmd.Flags().GetBool(tlsInsecureSkipVerifyFlag)
		if tlsInsecure {
			sharedlogging.FromContext(cmd.Context()).Infof("Disable certificate checks")
		}
		credential = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: tlsInsecure,
			RootCAs:            certPool,
		})
	}
	return credential, nil
}
