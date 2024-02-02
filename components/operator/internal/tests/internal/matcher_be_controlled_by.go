package internal

import (
	gomegaTypes "github.com/onsi/gomega/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func BeControlledBy(owner client.Object) gomegaTypes.GomegaMatcher {
	return BeOwnedBy(owner, func(matcher *beOwnedByMatcher) {
		matcher.controller = true
		matcher.blockOwnerDeletion = true
	})
}
