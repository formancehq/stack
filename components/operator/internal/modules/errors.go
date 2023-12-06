package modules

import (
	"fmt"
	"sync"
)

type serviceErrors struct {
	errors map[string]error
	mu     sync.Mutex
}

func (m *serviceErrors) Error() string {
	ret := ""
	for service, err := range m.errors {
		ret = fmt.Sprintf("%s: %s\r\n", service, err)
	}
	return ret
}

func (m *serviceErrors) setError(service string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.errors == nil {
		m.errors = map[string]error{}
	}
	m.errors[service] = err
}
