package internal

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	envPrefix = "formance_webhooks"
)

var envVarReplacer = strings.NewReplacer(".", "_", "-", "_")

func BindEnv(v *viper.Viper) {
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(envVarReplacer)
	v.AutomaticEnv()
}
