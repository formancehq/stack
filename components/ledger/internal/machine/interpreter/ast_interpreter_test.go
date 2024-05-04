package interpreter

import (
	"errors"
	"math/big"
	"testing"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/machine"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	tc := NewTestCase()
	tc.compile(t, `send [EUR/2 100] (
		source=@alice
		destination=@bob
	)`)
	tc.setBalance("alice", "EUR/2", 100)
	tc.expected = CaseResult{
		Printed: []machine.Value{},
		Postings: []Posting{
			{
				Amount:      100,
				Source:      "alice",
				Destination: "bob",
			},
		},
		Error: nil,
	}
	test(t, tc)
}

// ---- Test utilities
type CaseResult struct {
	Printed       []machine.Value
	Postings      []Posting
	Metadata      map[string]machine.Value
	Error         error
	ErrorContains string
}

type TestCase struct {
	program  *Program
	vars     map[string]string
	meta     map[string]metadata.Metadata
	balances map[string]int64
	expected CaseResult
}

func NewTestCase() TestCase {
	return TestCase{
		vars:     make(map[string]string),
		meta:     make(map[string]metadata.Metadata),
		balances: make(map[string]int64),
		expected: CaseResult{
			Printed:  []machine.Value{},
			Postings: []Posting{},
			Metadata: make(map[string]machine.Value),
			Error:    nil,
		},
	}
}

func (c *TestCase) compile(_ *testing.T, code string) {

	artifacts := CompileFull(code)
	c.program = artifacts.Program
}

func (c *TestCase) setBalance(account, _ string, amount int64) {
	c.balances[account] = amount
}

type StaticStore map[string]*AccountWithBalances
type AccountWithBalances struct {
	ledger.Account
	Balances map[string]*big.Int
}

func test(t *testing.T, testCase TestCase) {
	statement := testCase.program.Statements[0]
	expected := testCase.expected

	postings, err := EvalSend(statement.Amount, testCase.balances, statement.Source, statement.Destination)

	if expected.Error != nil {
		require.True(t, errors.Is(err, expected.Error), "got wrong error, want: %v, got: %v", expected.Error, err)
		if expected.ErrorContains != "" {
			require.ErrorContains(t, err, expected.ErrorContains)
		}
	} else {
		require.NoError(t, err)
	}
	if err != nil {
		return
	}

	if expected.Postings == nil {
		expected.Postings = make([]Posting, 0)
	}
	if expected.Metadata == nil {
		expected.Metadata = make(map[string]machine.Value)
	}

	assert.Equalf(t, expected.Postings, postings, "unexpected postings output: %v", postings)
}
