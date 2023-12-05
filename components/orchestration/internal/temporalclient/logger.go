package temporalclient

import (
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.temporal.io/sdk/log"
)

func keyvalsToMap(keyvals ...interface{}) map[string]any {
	ret := make(map[string]any)
	for i := 0; i < len(keyvals); i += 2 {
		ret[keyvals[i].(string)] = keyvals[i+1]
	}
	return ret
}

type logger struct {
	logger logging.Logger
}

func (l logger) Debug(msg string, keyvals ...interface{}) {
	l.logger.WithFields(keyvalsToMap(keyvals...)).Debugf(msg)
}

func (l logger) Info(msg string, keyvals ...interface{}) {
	l.logger.WithFields(keyvalsToMap(keyvals...)).Infof(msg)
}

func (l logger) Warn(msg string, keyvals ...interface{}) {
	l.logger.WithFields(keyvalsToMap(keyvals...)).Errorf(msg)
}

func (l logger) Error(msg string, keyvals ...interface{}) {
	l.logger.WithFields(keyvalsToMap(keyvals...)).Errorf(msg)
}

var _ log.Logger = (*logger)(nil)

func newLogger(l logging.Logger) *logger {
	return &logger{
		logger: l,
	}
}
