package interpreter

import (
	"math/big"
	"slices"
	"testing"
)

// Test.todo balances are updated
type SourceTestCase struct {
	Monetary    int64
	Balances    map[string]int64
	Source      Source
	Expected    []Sender
	ExpectedErr error
}

func runSourceTestCase(t *testing.T, tc *SourceTestCase) {
	got, err := EvalSource(tc.Monetary, tc.Balances, tc.Source)

	if !slices.Equal(got, tc.Expected) {
		t.Fatalf(`Expected = %#v; got = %#v`, tc.Expected, got)
	}

	if err != tc.ExpectedErr {
		t.Fatalf(`Expected err = %#v; got = %#v`, tc.ExpectedErr, err)
	}
}

func TestWorldSource(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 100,
		Source:   &AccountSrc{Name: "world"},
		Expected: []Sender{
			{Name: "world", Monetary: 100},
		},
	})
}

func TestSourceWithMissingBalance(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 100,
		Source:   &AccountSrc{Name: "src"},
		Balances: nil,
		Expected: []Sender{
			{Name: "src", Monetary: 0},
		},
		ExpectedErr: MissingFundsErr{Missing: 100},
	})
}

func TestSourceWithRightBalance(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 80,
		Source:   &AccountSrc{Name: "src"},
		Balances: map[string]int64{
			"src": 100,
		},
		Expected: []Sender{
			{Name: "src", Monetary: 80},
		},
	})
}

func TestCappedSourceNoOpWhenCapIsHigher(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 1,
		Source: &CappedSrc{
			Cap:    100,
			Source: &AccountSrc{Name: "world"},
		},
		Expected: []Sender{
			{Name: "world", Monetary: 1},
		},
	})
}

func TestCappedSource(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 100,
		Source: &CappedSrc{
			Cap:    20,
			Source: &AccountSrc{Name: "world"},
		},
		Expected: []Sender{
			{Name: "world", Monetary: 20},
		},
		ExpectedErr: MissingFundsErr{Missing: 80},
	})
}

func TestSeq(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 100,
		Balances: map[string]int64{
			"s1": 30,
			"s2": 50,
			"s3": 99999,
		},
		Source: &SeqSrc{
			Sources: []Source{
				&AccountSrc{Name: "s1"},
				&AccountSrc{Name: "s2"},
				&AccountSrc{Name: "s3"},
			},
		},
		Expected: []Sender{
			{Name: "s1", Monetary: 30},
			{Name: "s2", Monetary: 50},
			{Name: "s3", Monetary: 20},
		},
	})
}

func TestSeqFail(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 100,
		Balances: map[string]int64{
			"s1": 30,
		},
		Source: &SeqSrc{
			Sources: []Source{
				&AccountSrc{Name: "s1"},
				&AccountSrc{Name: "s2"},
			},
		},
		Expected: []Sender{
			{Name: "s1", Monetary: 30},
			{Name: "s2", Monetary: 0},
		},
		ExpectedErr: MissingFundsErr{Missing: 70},
	})
}

func TestCapInSeq(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 100,
		Balances: map[string]int64{
			"s1": 99999,
		},
		Source: &SeqSrc{
			Sources: []Source{
				&CappedSrc{
					Cap:    20,
					Source: &AccountSrc{Name: "world"},
				},
				&AccountSrc{Name: "s1"},
			},
		},
		Expected: []Sender{
			{Name: "world", Monetary: 20},
			{Name: "s1", Monetary: 80},
		},
	})
}

func TestSimpleAllotment(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 90,
		Balances: map[string]int64{
			"s1": 99999,
		},
		Source: &AllottedSrc{
			Allotments: []Allotment[Source]{
				{*big.NewRat(1, 3), &AccountSrc{Name: "s1"}},
				{*big.NewRat(2, 3), &AccountSrc{Name: "world"}},
			},
		},
		Expected: []Sender{
			{Name: "s1", Monetary: 30},
			{Name: "world", Monetary: 60},
		},
	})
}

func TestSimpleAllotmentOfUneven(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 99,
		Balances: map[string]int64{
			"s1": 99999,
		},
		Source: &AllottedSrc{
			Allotments: []Allotment[Source]{
				{*big.NewRat(1, 2), &AccountSrc{Name: "s1"}},
				{*big.NewRat(1, 2), &AccountSrc{Name: "world"}},
			},
		},
		Expected: []Sender{
			{Name: "s1", Monetary: 50},
			{Name: "world", Monetary: 49},
		},
	})
}

func TestComplexAllotmentOfUneven(t *testing.T) {
	runSourceTestCase(t, &SourceTestCase{
		Monetary: 99,
		Balances: map[string]int64{
			"s1": 1000,
			"s2": 1000,
			"s3": 1000,
			"s4": 1000,
			"s5": 1000,
		},
		Source: &AllottedSrc{
			Allotments: []Allotment[Source]{
				{*big.NewRat(1, 5), &AccountSrc{Name: "s1"}},
				{*big.NewRat(1, 5), &AccountSrc{Name: "s2"}},
				{*big.NewRat(1, 5), &AccountSrc{Name: "s3"}},
				{*big.NewRat(1, 5), &AccountSrc{Name: "s4"}},
				{*big.NewRat(1, 5), &AccountSrc{Name: "s5"}},
			},
		},
		Expected: []Sender{
			{Name: "s1", Monetary: 20},
			{Name: "s2", Monetary: 20},
			{Name: "s3", Monetary: 20},
			{Name: "s4", Monetary: 20},
			{Name: "s5", Monetary: 19},
		},
	})
}
