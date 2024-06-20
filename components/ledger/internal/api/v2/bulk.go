package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/engine"
	"github.com/formancehq/ledger/internal/machine"
	"github.com/formancehq/ledger/internal/opentelemetry/tracer"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"

	"github.com/formancehq/ledger/internal/api/backend"
	"github.com/formancehq/ledger/internal/engine/command"
)

const (
	ActionCreateTransaction = "CREATE_TRANSACTION"
	ActionAddMetadata       = "ADD_METADATA"
	ActionRevertTransaction = "REVERT_TRANSACTION"
	ActionDeleteMetadata    = "DELETE_METADATA"
)

type Bulk []Element

type Element struct {
	Action         string          `json:"action"`
	IdempotencyKey string          `json:"ik"`
	Data           json.RawMessage `json:"data"`
}

type Result struct {
	ErrorCode        string `json:"errorCode,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
	ErrorDetails     string `json:"errorDetails,omitempty"`
	Data             any    `json:"data,omitempty"`
	ResponseType     string `json:"responseType"` // Added for sdk generation (discriminator in oneOf)
}

func processBulkElement(
	ctx context.Context,
	l backend.Ledger,
	parameters command.Parameters,
	element Element,
	i int,
) (error, string, any) {
	switch element.Action {
	case ActionCreateTransaction:
		req := &ledger.TransactionRequest{}
		if err := json.Unmarshal(element.Data, req); err != nil {
			return fmt.Errorf("error parsing element %d: %s", i, err), "", nil
		}
		rs := req.ToRunScript()

		tx, err := l.CreateTransaction(ctx, parameters, *rs)
		if err != nil {
			var code string
			switch {
			case machine.IsInsufficientFundError(err):
				code = ErrInsufficientFund
			case engine.IsCommandError(err):
				code = ErrValidation
			default:
				code = sharedapi.ErrorInternal
			}
			return err, code, nil
		} else {
			return nil, "", tx
		}
	case ActionAddMetadata:
		type addMetadataRequest struct {
			TargetType string            `json:"targetType"`
			TargetID   json.RawMessage   `json:"targetId"`
			Metadata   metadata.Metadata `json:"metadata"`
		}
		req := &addMetadataRequest{}
		if err := json.Unmarshal(element.Data, req); err != nil {
			return fmt.Errorf("error parsing element %d: %s", i, err), "", nil
		}

		var targetID any
		switch req.TargetType {
		case ledger.MetaTargetTypeAccount:
			targetID = ""
		case ledger.MetaTargetTypeTransaction:
			targetID = big.NewInt(0)
		}
		if err := json.Unmarshal(req.TargetID, &targetID); err != nil {
			return err, "", nil
		}

		if err := l.SaveMeta(ctx, parameters, req.TargetType, targetID, req.Metadata); err != nil {
			var code string
			switch {
			case command.IsSaveMetaError(err, command.ErrSaveMetaCodeTransactionNotFound):
				code = sharedapi.ErrorCodeNotFound
			default:
				code = sharedapi.ErrorInternal
			}
			return err, code, nil
		} else {
			return nil, "", nil
		}
	case ActionRevertTransaction:
		type revertTransactionRequest struct {
			ID              *big.Int `json:"id"`
			Force           bool     `json:"force"`
			AtEffectiveDate bool     `json:"atEffectiveDate"`
		}
		req := &revertTransactionRequest{}
		if err := json.Unmarshal(element.Data, req); err != nil {
			return fmt.Errorf("error parsing element %d: %s", i, err), "", nil
		}

		tx, err := l.RevertTransaction(ctx, parameters, req.ID, req.Force, req.AtEffectiveDate)
		if err != nil {
			var code string
			switch {
			case engine.IsCommandError(err):
				code = ErrValidation
			default:
				code = sharedapi.ErrorInternal
			}
			return err, code, nil
		} else {
			return nil, "", tx
		}
	case ActionDeleteMetadata:
		type deleteMetadataRequest struct {
			TargetType string          `json:"targetType"`
			TargetID   json.RawMessage `json:"targetId"`
			Key        string          `json:"key"`
		}
		req := &deleteMetadataRequest{}
		if err := json.Unmarshal(element.Data, req); err != nil {
			return err, "", nil
		}

		var targetID any
		switch req.TargetType {
		case ledger.MetaTargetTypeAccount:
			targetID = ""
		case ledger.MetaTargetTypeTransaction:
			targetID = big.NewInt(0)
		}
		if err := json.Unmarshal(req.TargetID, &targetID); err != nil {
			return err, "", nil
		}

		err := l.DeleteMetadata(ctx, parameters, req.TargetType, targetID, req.Key)
		if err != nil {
			var code string
			switch {
			case command.IsDeleteMetaError(err, command.ErrSaveMetaCodeTransactionNotFound):
				code = sharedapi.ErrorCodeNotFound
			default:
				code = sharedapi.ErrorInternal
			}
			return err, code, nil
		} else {
			return nil, "", nil
		}
	}

	return fmt.Errorf("unknown action %s", element.Action), "", nil
}

func ProcessBulk(
	ctx context.Context,
	ledger backend.Ledger,
	bulk Bulk,
	continueOnFailure bool,
	parallel bool,
) ([]Result, bool, error) {

	ctx, span := tracer.Start(ctx, "Bulk")
	defer span.End()

	ret := make([]Result, len(bulk))

	errorsInBulk := false

	var bulkError = func(index int, action, code string, err error) {
		ret[index] = Result{
			ErrorCode:        code,
			ErrorDescription: err.Error(),
			ResponseType:     "ERROR",
		}
		errorsInBulk = true
	}

	var bulkSuccess = func(index int, action string, data any) {
		ret[index] = Result{
			Data:         data,
			ResponseType: action,
		}
	}

	wg := sync.WaitGroup{}

	for i, element := range bulk {
		parameters := command.Parameters{
			DryRun:         false,
			IdempotencyKey: element.IdempotencyKey,
		}

		wg.Add(1)

		go func(element Element, index int) {
			err, code, data := processBulkElement(ctx, ledger, parameters, element, index)
			if err != nil {
				bulkError(index, element.Action, code, err)
			} else {
				bulkSuccess(index, element.Action, data)
			}
			wg.Done()
		}(element, i)

		if !parallel {
			wg.Wait()
		}
	}

	if parallel {
		wg.Wait()
	}

	return ret, errorsInBulk, nil
}
