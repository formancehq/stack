package common

import (
	"github.com/formancehq/go-libs/platform/postgres"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"github.com/formancehq/ledger/internal/machine"
	"github.com/formancehq/ledger/internal/machine/script/compiler"
)

type ErrInsufficientFunds = ledgercontroller.ErrInsufficientFunds
type ErrAlreadyReverted = ledgercontroller.ErrAlreadyReverted
type ErrInvalidQuery = ledgercontroller.ErrInvalidQuery
type ErrMissingFeature = ledgercontroller.ErrMissingFeature
type ErrConstraintsFailed = postgres.ErrConstraintsFailed

var ErrNoPostings = ledgercontroller.ErrNoPostings
var ErrNotFound = ledgercontroller.ErrNotFound

// todo: should not depend on machine
// notes(gfyrag): I will wait for the new numscript interpreter to refactor error handling

type ErrInvalidVars = machine.ErrInvalidVars
type ErrCompileErrorList = compiler.CompileErrorList
type ErrMetadataOverride = machine.ErrMetadataOverride
