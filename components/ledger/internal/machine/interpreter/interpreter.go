package interpreter

import (
	"math/big"
	"slices"
)

type Allotment[T interface{}] struct {
	Ratio big.Rat
	Value T
}

func makeAllotment[T interface{}](monetary int64, allotments []Allotment[T]) []int64 {
	parts := make([]int64, len(allotments))

	var totalAllocated int64

	for i, allot := range allotments {
		var product big.Rat
		product.Mul(&allot.Ratio, big.NewRat(monetary, 1))

		floored := product.Num().Int64() / product.Denom().Int64()

		parts[i] = floored
		totalAllocated += floored
	}

	for i := range parts {
		if totalAllocated >= monetary {
			break
		}

		parts[i]++
		totalAllocated++
	}

	return parts
}

type Posting struct {
	Source      string
	Destination string
	Amount      int64
}

type ReconcileError struct {
	error
}

func Reconcile(senders []Sender, receivers []Receiver) ([]Posting, error) {
	var postings []Posting

	for {
		receiver, empty := popStack(&receivers)
		if empty {
			break
		}

		sender, empty := popStack(&senders)
		if empty {
			return nil, ReconcileError{}
		}

		var postingAmount int64
		if sender.Monetary == receiver.Monetary {
			postingAmount = sender.Monetary
		} else if sender.Monetary > receiver.Monetary {
			senders = append(senders, Sender{
				Name:     sender.Name,
				Monetary: sender.Monetary - receiver.Monetary,
			})
			postingAmount = receiver.Monetary
		} else /* if sender.Monetary < receiver.Monetary */ {
			receivers = append(receivers, Receiver{
				Name:     receiver.Name,
				Monetary: receiver.Monetary - sender.Monetary,
			})
			postingAmount = sender.Monetary
		}

		postings = append(postings, Posting{
			Source:      sender.Name,
			Destination: receiver.Name,
			Amount:      postingAmount,
		})
	}

	slices.Reverse(postings)
	return postings, nil
}

func popStack[T any](stack *[]T) (T, bool) {
	l := len(*stack)
	if l == 0 {
		var t T
		return t, true
	}

	popped := (*stack)[l-1]
	*stack = (*stack)[:l-1]
	return popped, false
}

func EvalSend(
	monetary int64,
	balances map[string]int64,
	source Source,
	destination Destination,
) ([]Posting, error) {
	senders, err := EvalSource(monetary, balances, source)
	if err != nil {
		return nil, err
	}

	receivers, err := EvalDestination(monetary, destination)
	if err != nil {
		return nil, err
	}

	postings, err := Reconcile(senders, receivers)
	if err != nil {
		return nil, err
	}

	return postings, nil
}
