package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
)

func GetModuleVersion(stack *v1beta1.Stack, defaultVersion string) string {
	if defaultVersion != "" {
		return defaultVersion
	}
	if stack.GetVersion() != "" {
		return stack.GetVersion()
	}

	return "latest"
}
