package auth

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	AuthEnabled                  = "auth-enabled"
	AuthIssuerFlag               = "auth-issuer"
	AuthIssuersFlag              = "auth-issuers"
	AuthReadKeySetMaxRetriesFlag = "auth-read-key-set-max-retries"
	AuthCheckScopesFlag          = "auth-check-scopes"
	AuthServiceFlag              = "auth-service"
)

func InitAuthFlags(flags *flag.FlagSet) {
	flags.Bool(AuthEnabled, false, "Enable auth")
	flags.String(AuthIssuerFlag, "", "Issuer (single issuer, for backward compatibility)")
	flags.StringSlice(AuthIssuersFlag, nil, "Trusted issuers (comma-separated, e.g. --auth-issuers=https://issuer1,https://issuer2)")
	flags.Int(AuthReadKeySetMaxRetriesFlag, 10, "ReadKeySetMaxRetries")
	flags.Bool(AuthCheckScopesFlag, false, "CheckScopes")
	flags.String(AuthServiceFlag, "", "Service")
}

func CLIAuthModule() fx.Option {
	authIssuer := viper.GetString(AuthIssuerFlag)
	authIssuers := viper.GetStringSlice(AuthIssuersFlag)

	// Merge --auth-issuer into --auth-issuers for backward compatibility
	if authIssuer != "" {
		found := false
		for _, iss := range authIssuers {
			if iss == authIssuer {
				found = true
				break
			}
		}
		if !found {
			authIssuers = append(authIssuers, authIssuer)
		}
	}

	return Module(ModuleConfig{
		Enabled:              viper.GetBool(AuthEnabled),
		Issuers:              authIssuers,
		ReadKeySetMaxRetries: viper.GetInt(AuthReadKeySetMaxRetriesFlag),
		CheckScopes:          viper.GetBool(AuthCheckScopesFlag),
		Service:              viper.GetString(AuthServiceFlag),
	})
}
