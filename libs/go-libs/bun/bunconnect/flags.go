package bunconnect

import (
	"context"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/lib/pq"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"time"
)

func InitFlags(flags *pflag.FlagSet) {
	flags.String("postgres-uri", "", "Postgres URI")
	flags.Bool("postgres-aws-enable-iam", false, "Enable AWS IAM authentication")
	flags.Int("postgres-max-idle-conns", 0, "Max Idle connections")
	flags.Duration("postgres-conn-max-idle-time", time.Minute, "Max Idle time for connections")
	flags.Int("postgres-max-open-conns", 20, "Max opened connections")
}

func ConnectionOptionsFromFlags(v *viper.Viper, output io.Writer, debug bool) (*ConnectionOptions, error) {
	var connector func(string) (driver.Connector, error)
	if v.GetBool("postgres-aws-enable-iam") {
		cfg, err := config.LoadDefaultConfig(context.Background(), iam.LoadOptionFromViper(v))
		if err != nil {
			return nil, err
		}

		connector = func(s string) (driver.Connector, error) {
			return &iamConnector{
				dsn: s,
				driver: &iamDriver{
					awsConfig: cfg,
				},
			}, nil
		}
	} else {
		connector = func(dsn string) (driver.Connector, error) {
			return pq.NewConnector(dsn)
		}
	}
	return &ConnectionOptions{
		DatabaseSourceName: v.GetString("postgres-uri"),
		Debug:              debug,
		Writer:             output,
		MaxIdleConns:       v.GetInt("postgres-max-idle-conns"),
		ConnMaxIdleTime:    v.GetDuration("postgres-conn-max-idle-time"),
		MaxOpenConns:       v.GetInt("postgres-max-open-conns"),
		Connector:          connector,
	}, nil
}
