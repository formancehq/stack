package iam

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitFlags(flags *pflag.FlagSet) {
	flags.String("aws-region", "", "Specify AWS region")
	flags.String("aws-access-key-id", "", "AWS access key id")
	flags.String("aws-secret-access-key", "", "AWS secret access key")
	flags.String("aws-session-token", "", "AWS session token")
	flags.String("aws-profile", "", "AWS profile")
}

func LoadOptionFromViper(v *viper.Viper) func(opts *config.LoadOptions) error {
	return func(opts *config.LoadOptions) error {
		if v.GetString("aws-access-key-id") != "" {
			opts.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     v.GetString("aws-access-key-id"),
					SecretAccessKey: v.GetString("aws-secret-access-key"),
					SessionToken:    v.GetString("aws-session-token"),
					Source:          "flags",
					CanExpire:       false,
				}, nil
			})
		}
		opts.Region = v.GetString("aws-region")
		opts.SharedConfigProfile = v.GetString("aws-profile")

		return nil
	}
}
