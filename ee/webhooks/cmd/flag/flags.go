package flag

import (
	"strings"
	"time"

	cache "github.com/formancehq/webhooks/internal/app/cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	LogLevel = "log-level"
	Listen   = "listen"

	MaxCall   = "max-call"
	MaxRetry  = "max-retry"
	TimeOut   = "time-out"
	DelayPull = "delay-pull"

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

	flagSet.Int(TimeOut, 2000, "Set time out for hook request (ms)")
	flagSet.Int(MaxRetry, 60, "Set max  number of retries for failed attempt")
	flagSet.Int(MaxCall, 20, "Set max number of http request at the same time")
	flagSet.Int(DelayPull, 1, "Period of time for pulling the database and synchronise cached data")

	flagSet.StringSlice(KafkaTopics, []string{DefaultKafkaTopic}, "Kafka topics")
	flagSet.Bool(AutoMigrate, false, "auto migrate database")
}

func LoadEnv(v *viper.Viper) {
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()
}

func init() {
	LoadEnv(viper.GetViper())
}

func LoadRunnerParams() cache.CacheParams {
	stateParams := cache.DefaultCacheParams()
	stateParams.MaxRetry = viper.GetInt(MaxRetry)
	stateParams.MaxCall = viper.GetInt(MaxCall)
	stateParams.TimeOut = viper.GetInt(TimeOut)
	stateParams.DelayPull = viper.GetInt(DelayPull)

	return stateParams
}
