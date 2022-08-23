package cmd

import (
	"github.com/numary/webhooks/pkg/env"
	"github.com/spf13/viper"
)

func init() {
	env.LoadEnv(viper.GetViper())
}
