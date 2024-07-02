package interpreter

// ---- Intermediate ASt
type Source interface {
	trySending(int64, *evalSourceCtx) int64
}

type AccountSrc struct {
	Name string
}

type CappedSrc struct {
	Cap    int64
	Source Source
}

type SeqSrc struct {
	Sources []Source
}

type AllottedSrc struct {
	Allotments []Allotment[Source]
}

// ---- Eval
type Sender struct {
	Name     string
	Monetary int64
}

type evalSourceCtx struct {
	Balances map[string]int64
	Senders  []Sender
}

func (s *AccountSrc) trySending(monetary int64, ctx *evalSourceCtx) int64 {
	if s.Name != "world" {
		monetary = min(ctx.Balances[s.Name], monetary)
	}

	ctx.Senders = append(ctx.Senders, Sender{
		Name:     s.Name,
		Monetary: monetary,
	})
	if ctx.Balances != nil {
		ctx.Balances[s.Name] -= monetary
	}

	return monetary
}

func (s *CappedSrc) trySending(monetary int64, ctx *evalSourceCtx) int64 {
	return s.Source.trySending(min(monetary, s.Cap), ctx)
}

func (s *SeqSrc) trySending(monetary int64, ctx *evalSourceCtx) int64 {
	var sentTotal int64
	for _, source := range s.Sources {
		sentTotal += source.trySending(monetary-sentTotal, ctx)
	}
	return sentTotal
}

// Precondition: sum of allotments is 1
func (s *AllottedSrc) trySending(monetary int64, ctx *evalSourceCtx) int64 {
	parts := makeAllotment(monetary, s.Allotments)

	var sentTotal int64
	for i, a := range s.Allotments {
		sentTotal += a.Value.trySending(parts[i], ctx)
	}

	return sentTotal
}

// ---- Public interface
type MissingFundsErr struct {
	error
	Missing int64
}

func EvalSource(monetary int64, balances map[string]int64, source Source) ([]Sender, error) {
	ctx := evalSourceCtx{
		Balances: balances,
	}

	sentTotal := source.trySending(monetary, &ctx)

	if sentTotal < monetary {
		return ctx.Senders, MissingFundsErr{Missing: monetary - sentTotal}
	}

	return ctx.Senders, nil
}
