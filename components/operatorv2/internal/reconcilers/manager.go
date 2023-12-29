package reconcilers

import (
	"github.com/formancehq/operator/v2/internal/controller/shared"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Manager interface {
	ctrl.Manager
	GetPlatform() shared.Platform
}

type defaultManager struct {
	ctrl.Manager
	platform shared.Platform
}

func (d defaultManager) GetPlatform() shared.Platform {
	return d.platform
}

var _ Manager = (*defaultManager)(nil)

func newDefaultManager(m ctrl.Manager, platform shared.Platform) *defaultManager {
	return &defaultManager{
		Manager:  m,
		platform: platform,
	}
}
