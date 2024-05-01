package interpreter

// Ast
type Destination interface {
	receive(int64, *evalDestinationCtx) int64
}

type AccountDest struct {
	Name string
}

type SeqDest struct {
	Destinations []Destination
}

type CappedDest struct {
	Cap         int64
	Destination Destination
}

type AllottedDest struct {
	Allotments []Allotment[Destination]
}

// eval

type evalDestinationCtx struct {
	Receivers []Receiver
}

func (d *AccountDest) receive(monetary int64, ctx *evalDestinationCtx) int64 {
	ctx.Receivers = append(ctx.Receivers, Receiver{
		Name:     d.Name,
		Monetary: monetary,
	})
	return monetary
}

func (d *SeqDest) receive(monetary int64, ctx *evalDestinationCtx) int64 {
	var receivedTotal int64
	for _, destination := range d.Destinations {
		receivedTotal += destination.receive(monetary-receivedTotal, ctx)
		if receivedTotal >= monetary {
			break
		}
	}

	return receivedTotal
}

func (d *CappedDest) receive(monetary int64, ctx *evalDestinationCtx) int64 {
	return d.Destination.receive(min(monetary, d.Cap), ctx)
}

func (d *AllottedDest) receive(monetary int64, ctx *evalDestinationCtx) int64 {
	var receivedTotal int64
	allot := makeAllotment(monetary, d.Allotments)
	for i, a := range d.Allotments {
		receivedTotal += a.Value.receive(allot[i], ctx)
	}
	return receivedTotal
}

// public api

type Receiver struct {
	Name     string
	Monetary int64
}

func EvalDestination(monetary int64, destination Destination) ([]Receiver, error) {
	ctx := evalDestinationCtx{}
	destination.receive(monetary, &ctx)
	return ctx.Receivers, nil
}
