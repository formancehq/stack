package flag

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	LogLevel = "log-level"
	Listen   = "listen"
	Worker   = "worker"

	RetryPeriod     = "retry-period"
	AbortAfter      = "abort-after"
	MinBackoffDelay = "min-backoff-delay"
	MaxBackoffDelay = "max-backoff-delay"

	KafkaTopics = "kafka-topics"
	AutoMigrate = "auto-migrate"
)

const (
	DefaultBindAddressServer = ":8080"

	DefaultPostgresConnString = "postgresql://webhooks:webhooks@127.0.0.1/webhooks?sslmode=disable"

	DefaultKafkaTopic = "default"
)

var (
	DefaultRetryPeriod = time.Minute
)

func Init(flagSet *pflag.FlagSet) {
	flagSet.String(LogLevel, logrus.InfoLevel.String(), "Log level")

	flagSet.String(Listen, DefaultBindAddressServer, "server HTTP bind address")
	flagSet.Duration(RetryPeriod, DefaultRetryPeriod, "worker retry period")
	flagSet.Bool(Worker, false, "Enable worker on server")

	flagSet.StringSlice(KafkaTopics, []string{DefaultKafkaTopic}, "Kafka topics")

	flagSet.Duration(AbortAfter, 30*24*time.Hour, "consider a webhook as failed after retrying it for this duration.")
	flagSet.Duration(MinBackoffDelay, time.Minute, "minimum backoff delay")
	flagSet.Duration(MaxBackoffDelay, time.Hour, "maximum backoff delay")
	flagSet.Bool(AutoMigrate, false, "auto migrate database")
}
