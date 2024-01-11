package core

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
)

const (
	OperatorVersionKey = "formance.com/operator-version"
	OperatorVersion    = "2"
)

func GetVersion(stack *v1beta1.Stack, defaultVersion string) string {
	if defaultVersion == "" {
		return stack.GetVersion()
	}
	return defaultVersion
}
