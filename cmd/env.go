package cmd

import (
	"github.com/numary/webhooks-cloud/internal/env"
	"github.com/spf13/viper"
)

func init() {
	env.LoadEnv(viper.GetViper())
}
