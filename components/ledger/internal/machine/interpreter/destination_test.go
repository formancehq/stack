package interpreter

import (
	"math/big"
	"slices"
	"testing"
)

type DestinationTestCase struct {
	Monetary    int64
	Destination Destination
	Expected    []Receiver
	ExpectedErr error
}

func runDestinationTestCase(t *testing.T, tc *DestinationTestCase) {
	got, err := EvalDestination(tc.Monetary, tc.Destination)

	if !slices.Equal(got, tc.Expected) {
		t.Fatalf(`Expected = %#v; got = %#v`, tc.Expected, got)
	}

	if err != tc.ExpectedErr {
		t.Fatalf(`Expected err = %#v; got = %#v`, tc.ExpectedErr, err)
	}
}

func TestAddressDest(t *testing.T) {
	runDestinationTestCase(t, &DestinationTestCase{
		Monetary:    100,
		Destination: &AccountDest{Name: "x"},
		Expected: []Receiver{
			{"x", 100},
		},
	})
}

func TestFirstOfSeqReceivesAllAmount(t *testing.T) {
	runDestinationTestCase(t, &DestinationTestCase{
		Monetary: 100,
		Destination: &SeqDest{
			Destinations: []Destination{
				&AccountDest{Name: "d1"},
				&AccountDest{Name: "d2"},
				&AccountDest{Name: "d3"},
			},
		},
		Expected: []Receiver{
			{"d1", 100},
		},
	})
}

func TestEmptySeqHasNoReceivers(t *testing.T) {
	runDestinationTestCase(t, &DestinationTestCase{
		Monetary: 100,
		Destination: &SeqDest{
			Destinations: []Destination{},
		},
		Expected: []Receiver{},
	})
}

func TestCappedDestIsNoopWhenCapIsHigherThanAmt(t *testing.T) {
	runDestinationTestCase(t, &DestinationTestCase{
		Monetary: 10,
		Destination: &CappedDest{
			Cap:         99999,
			Destination: &AccountDest{Name: "x"},
		},
		Expected: []Receiver{
			{"x", 10},
		},
	})
}

func TestCappedDest(t *testing.T) {
	runDestinationTestCase(t, &DestinationTestCase{
		Monetary: 100,
		Destination: &CappedDest{
			Cap:         10,
			Destination: &AccountDest{Name: "x"},
		},
		Expected: []Receiver{
			{"x", 10},
		},
	})
}

func TestCappedDestInSeq(t *testing.T) {
	/*
		send 100 (
			destination = {
				max 10 to @d1
				@d2
			}
		)
		=> [(@d1, 10), (@d2, 90)]
	*/

	runDestinationTestCase(t, &DestinationTestCase{
		Monetary: 100,
		Destination: &SeqDest{
			Destinations: []Destination{
				&CappedDest{
					Cap:         10,
					Destination: &AccountDest{Name: "d1"},
				},
				&AccountDest{Name: "d2"},
			},
		},
		Expected: []Receiver{
			{"d1", 10},
			{"d2", 90},
		},
	})
}

func TestAllottedDest(t *testing.T) {
	runDestinationTestCase(t, &DestinationTestCase{
		Monetary: 90,
		Destination: &AllottedDest{
			Allotments: []Allotment[Destination]{
				{*big.NewRat(1, 3), &AccountDest{"x"}},
				{*big.NewRat(2, 3), &AccountDest{"y"}},
			},
		},
		Expected: []Receiver{
			{"x", 30},
			{"y", 60},
		},
	})
}
