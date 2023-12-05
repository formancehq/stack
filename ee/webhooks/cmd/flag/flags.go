package flag

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	Debug    = "debug"
	LogLevel = "log-level"
	Listen   = "listen"
	Worker   = "worker"

	RetriesSchedule = "retries-schedule"
	RetriesCron     = "retries-cron"

	StoragePostgresConnString = "storage-postgres-conn-string"

	KafkaTopics = "kafka-topics"
)

const (
	DefaultBindAddressServer = ":8080"

	DefaultPostgresConnString = "postgresql://webhooks:webhooks@127.0.0.1/webhooks?sslmode=disable"

	DefaultKafkaTopic = "default"
)

var (
	DefaultRetriesSchedule = []time.Duration{time.Minute, 5 * time.Minute, 30 * time.Minute, 5 * time.Hour, 24 * time.Hour}
	DefaultRetriesCron     = time.Minute
)

func Init(flagSet *pflag.FlagSet) {
	flagSet.Bool(Debug, false, "Debug mode")
	flagSet.String(LogLevel, logrus.InfoLevel.String(), "Log level")

	flagSet.String(Listen, DefaultBindAddressServer, "server HTTP bind address")
	flagSet.DurationSlice(RetriesSchedule, DefaultRetriesSchedule, "worker retries schedule")
	flagSet.Duration(RetriesCron, DefaultRetriesCron, "worker retries cron")
	flagSet.String(StoragePostgresConnString, DefaultPostgresConnString, "Postgres connection string")
	flagSet.Bool(Worker, false, "Enable worker on server")

	flagSet.StringSlice(KafkaTopics, []string{DefaultKafkaTopic}, "Kafka topics")
}

func LoadEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}

func init() {
	LoadEnv(viper.GetViper())
}
