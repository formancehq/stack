package bunconnect

import (
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
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

func ConnectionOptionsFromFlags(v *viper.Viper, output io.Writer, debug bool) ConnectionOptions {
	return ConnectionOptions{
		DatabaseSourceName: v.GetString("postgres-uri"),
		Debug:              debug,
		Writer:             output,
		MaxIdleConns:       v.GetInt("postgres-max-idle-conns"),
		ConnMaxIdleTime:    v.GetDuration("postgres-conn-max-idle-time"),
		MaxOpenConns:       v.GetInt("postgres-max-open-conns"),
		Opener: func() Opener {
			if v.GetBool("postgres-aws-enable-iam") {
				return IAMOpener(iam.LoadOptionFromViper(v))
			}
			return DefaultOpener
		}(),
	}
}
