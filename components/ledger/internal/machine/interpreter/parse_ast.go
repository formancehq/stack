package interpreter

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/formancehq/ledger/internal/machine"
	"github.com/formancehq/ledger/internal/machine/script/parser"
)

type SendStatement struct {
	Amount      int64
	Source      Source
	Destination Destination
}

type Program struct {
	Statements []SendStatement
}

type parseVisitor struct {
	errListener *ErrorListener
}

type CompileError struct {
	StartL, StartC int
	EndL, EndC     int
	Msg            string
}

type CompileArtifacts struct {
	Errors  []CompileError
	Program *Program
}

type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []CompileError
}

func CompileFull(input string) CompileArtifacts {
	artifacts := CompileArtifacts{}

	errListener := &ErrorListener{}

	is := antlr.NewInputStream(input)
	lexer := parser.NewNumScriptLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.LexerDefaultTokenChannel)
	p := parser.NewNumScriptParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)

	p.BuildParseTrees = true

	tree := p.Script()

	artifacts.Errors = append(artifacts.Errors, errListener.Errors...)
	if len(errListener.Errors) != 0 {
		return artifacts
	}

	visitor := parseVisitor{
		errListener: errListener,
	}

	program, err := visitor.visitScript(tree)
	if err != nil {
		fmt.Printf("#### err: %#v\n", err)
		panic("TODO handle")
		// return artifacts
	}

	artifacts.Program = program
	return artifacts
}

func (p *parseVisitor) visitScript(c parser.IScriptContext) (*Program, *CompileError) {
	program := &Program{}

	switch c := c.(type) {
	case *parser.ScriptContext:
		// TODO fetch vars

		for _, stmt := range c.GetStmts() {
			switch c := stmt.(type) {
			case *parser.SendContext:
				statement, err := p.visitSend(c)
				if err != nil {
					return nil, err
				}
				program.Statements = append(program.Statements, *statement)

			default:
				panic("TODO unimplemented statement type")
			}
		}

	default:
		panic("TODO unimplemented script type")
	}

	return program, nil
}

func (p *parseVisitor) visitSend(c *parser.SendContext) (*SendStatement, *CompileError) {
	statement := SendStatement{}

	monExpr := c.GetMon().(*parser.ExprLiteralContext)
	mon, err := p.visitMonetaryLit(monExpr.GetLit())
	if err != nil {
		return nil, err
	}
	statement.Amount = int64(mon.Uint64())

	switch c := c.GetSrc().(type) {
	case *parser.SrcContext:
		src, err := p.visitSource(c.Source())
		if err != nil {
			return nil, err
		}
		statement.Source = src

	case *parser.SrcAllotmentContext:
		panic("TODO handle src allotment")
	}

	dest, err := p.visitDestination(c.GetDest())
	if err != nil {
		return nil, err
	}
	statement.Destination = dest

	return &statement, nil
}

func (p *parseVisitor) visitMonetaryLit(c parser.ILiteralContext) (*machine.MonetaryInt, *CompileError) {
	monCtx, ok := c.(*parser.LitMonetaryContext)
	if !ok {
		return nil, InternalError(c)
	}

	amt, err := machine.ParseMonetaryInt(monCtx.Monetary().GetAmt().GetText())
	if err != nil {
		return nil, LogicError(c, err)
	}
	return amt, nil

}

func (p *parseVisitor) visitAccountLit(c parser.ILiteralContext) (*machine.AccountAddress, *CompileError) {
	switch c := c.(type) {
	case *parser.LitAccountContext:
		account := machine.AccountAddress(c.GetText()[1:])
		return &account, nil

	default:
		panic("INVALID TOKEN (expeting accoutn)")
	}
}

func (p *parseVisitor) visitSource(c parser.ISourceContext) (Source, *CompileError) {
	switch c := c.(type) {
	case *parser.SrcAccountContext:
		accountLit := c.SourceAccount().GetAccount().(*parser.ExprLiteralContext)
		account, err := p.visitAccountLit(accountLit.GetLit())
		if err != nil {
			return nil, err
		}
		return &AccountSrc{Name: string(*account)}, nil

	case *parser.SrcMaxedContext:
		panic("TODO handle maxed")
		// accounts, _, subsourceFallback, compErr := p.visitSource(c.SourceMaxed().GetSrc(), pushAsset, false)
		// if compErr != nil {
		// 	return nil, nil, nil, compErr
		// }
		// ty, _, compErr := p.visitExpr(c.SourceMaxed().GetMax(), true)
		// if compErr != nil {
		// 	return nil, nil, nil, compErr
		// }
		// if ty != machine.TypeMonetary {
		// 	return nil, nil, nil, LogicError(c, errors.New("wrong type: expected monetary as max"))
		// }
		// for k, v := range accounts {
		// 	neededAccounts[k] = v
		// }
		// err := p.Bump(1)
		// if err != nil {
		// 	return nil, nil, nil, LogicError(c, err)
		// }
		// if subsourceFallback != nil {
		// 	p.PushAddress(machine.Address(*subsourceFallback))
		// 	err := p.Bump(2)
		// 	if err != nil {
		// 		return nil, nil, nil, LogicError(c, err)
		// 	}
		// 	err = p.PushInteger(machine.NewNumber(2))
		// 	if err != nil {
		// 		return nil, nil, nil, LogicError(c, err)
		// 	}
		// } else {
		// 	err := p.Bump(1)
		// 	if err != nil {
		// 		return nil, nil, nil, LogicError(c, err)
		// 	}
		// }

	case *parser.SrcInOrderContext:
		seq := &SeqSrc{}

		for _, source := range c.SourceInOrder().GetSources() {
			src, err := p.visitSource(source)
			if err != nil {
				return nil, err
			}

			seq.Sources = append(seq.Sources, src)
		}

		return seq, nil
	}

	return nil, nil
}

func (p *parseVisitor) visitDestination(c parser.IDestinationContext) (Destination, *CompileError) {
	switch c := c.(type) {
	case *parser.DestAccountContext:
		e, ok := c.Expression().(*parser.ExprLiteralContext)
		if !ok {
			return nil, InternalError(c)
		}
		account, err := p.visitAccountLit(e.GetLit())
		if err != nil {
			return nil, err
		}
		return &AccountDest{Name: string(*account)}, nil

	case *parser.DestInOrderContext:
		// dests := c.DestinationInOrder().GetDests()
		// amounts := c.DestinationInOrder().GetAmounts()
		// n := len(dests)

		// for i := 0; i < n; i++ {
		// 	ty, _, compErr := p.VisitExpr(amounts[i], true)
		// 	if compErr != nil {
		// 		return compErr
		// 	}
		// 	if ty != machine.TypeMonetary {
		// 		return LogicError(c, errors.New("wrong type: expected monetary as max"))
		// 	}

		// 	if compErr != nil {
		// 		return compErr
		// 	}

		// 	if err != nil {
		// 		return LogicError(c, err)
		// 	}

		// }

		// cerr := p.VisitKeptOrDestination(c.DestinationInOrder().GetRemainingDest())

		return nil, nil
	case *parser.DestAllotmentContext:
		// err := p.VisitDestinationAllotment(c.DestinationAllotment())
		// return err
		return nil, nil
	default:
		return nil, InternalError(c)
	}
}

func LogicError(c antlr.ParserRuleContext, err error) *CompileError {
	endC := c.GetStop().GetColumn() + len(c.GetStop().GetText())
	return &CompileError{
		StartL: c.GetStart().GetLine(),
		StartC: c.GetStart().GetColumn(),
		EndL:   c.GetStop().GetLine(),
		EndC:   endC,
		Msg:    err.Error(),
	}
}

const InternalErrorMsg = "internal compiler error, please report to the issue tracker"

func InternalError(c antlr.ParserRuleContext) *CompileError {
	return &CompileError{
		StartL: c.GetStart().GetLine(),
		StartC: c.GetStart().GetColumn(),
		EndL:   c.GetStop().GetLine(),
		EndC:   c.GetStop().GetColumn(),
		Msg:    InternalErrorMsg,
	}
}
