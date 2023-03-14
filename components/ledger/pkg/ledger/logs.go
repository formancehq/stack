package ledger

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
	"github.com/pkg/errors"
)

type appendLog func(context.Context, ...core.Log) <-chan error

type LogHandler struct {
	appendLog appendLog
	log       core.Log

	errChan <-chan error
}

func writeLog(ctx context.Context, appendLog appendLog, log core.Log) (*LogHandler, error) {
	lh := &LogHandler{
		appendLog: appendLog,
		log:       log,
	}
	if err := lh.write(ctx); err != nil {
		return nil, errors.Wrap(err, "writing logs")
	}
	return lh, nil
}

func (ls *LogHandler) write(ctx context.Context) error {
	ls.errChan = ls.appendLog(ctx, ls.log)

	return nil
}

func (ls *LogHandler) Wait(ctx context.Context) error {
	if ls.errChan == nil {
		// Nothing to wait on
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ls.errChan:
		if err != nil {
			return errors.Wrap(err, "appending logs")
		}
	}

	return nil
}
