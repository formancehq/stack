package flag

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	Debug                     = "debug"
	LogLevel                  = "log-level"
	Listen                    = "listen"
	Worker                    = "worker"
	RetriesSchedule           = "retries-schedule"
	RetriesCron               = "retries-cron"
	StoragePostgresConnString = "storage-postgres-conn-string"
	KafkaTopics               = "kafka-topics"
)

const (
	DefaultBindAddressServer  = ":8080"
	DefaultPostgresConnString = "postgresql://webhooks:webhooks@127.0.0.1/webhooks?sslmode=disable"
	DefaultKafkaTopic         = "default"
)

var (
	DefaultRetriesSchedule = []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 5 * time.Hour, 24 * time.Hour}
	DefaultRetriesCron     = time.Minute
)

var ErrScheduleInvalid = errors.New("the retry schedule should only contain durations of at least 1 second")

func Init(flagSet *pflag.FlagSet) (retriesSchedule []time.Duration, err error) {
	flagSet.Bool(Debug, false, "Debug mode")
	flagSet.String(LogLevel, logrus.InfoLevel.String(), "Log level")

	flagSet.String(Listen, DefaultBindAddressServer, "server HTTP bind address")
	flagSet.DurationSlice(RetriesSchedule, DefaultRetriesSchedule, "worker retries schedule")
	flagSet.Duration(RetriesCron, DefaultRetriesCron, "worker retries cron")
	flagSet.String(StoragePostgresConnString, DefaultPostgresConnString, "Postgres connection string")
	flagSet.Bool(Worker, false, "Enable worker on server")

	flagSet.StringSlice(KafkaTopics, []string{DefaultKafkaTopic}, "Kafka topics")

	retriesSchedule, err = flagSet.GetDurationSlice(RetriesSchedule)
	if err != nil {
		return nil, errors.Wrap(err, "flagSet.GetDurationSlice")
	}

	// Check that the schedule is valid
	for _, s := range retriesSchedule {
		if s < time.Second {
			return nil, ErrScheduleInvalid
		}
	}

	return retriesSchedule, nil
}

func LoadEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}

func init() {
	LoadEnv(viper.GetViper())
}
