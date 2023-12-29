package reconcilers

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

type Manager interface {
	ctrl.Manager
	GetPlatform() Platform
}

type defaultManager struct {
	ctrl.Manager
	platform Platform
}

func (d defaultManager) GetPlatform() Platform {
	return d.platform
}

var _ Manager = (*defaultManager)(nil)

func newDefaultManager(m ctrl.Manager, platform Platform) *defaultManager {
	return &defaultManager{
		Manager:  m,
		platform: platform,
	}
}
