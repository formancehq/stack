package databases

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/settings"
)

func NewHostSetting(name, value string, stacks ...string) *v1beta1.Settings {
	return settings.New(name, "databases.host", value, stacks...)
}
