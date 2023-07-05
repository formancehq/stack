package compiler

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/formancehq/ledger/pkg/machine/internal"
	"github.com/formancehq/ledger/pkg/machine/script/parser"
	"github.com/formancehq/ledger/pkg/machine/vm/program"
	"github.com/pkg/errors"
)

type parseVisitor struct {
	errListener *ErrorListener
	vars        map[string]internal.Type
}

func (p *parseVisitor) isWorld(expr parser.IExpressionContext) bool {
	if lit, ok := expr.(*parser.ExprLiteralContext); ok {
		_, value, _ := p.CompileLit(lit.GetLit())
		return internal.ValueEquals(value, internal.AccountAddress("world"))
	} else {
		return false
	}
}

func (p *parseVisitor) CompileExprTy(c parser.IExpressionContext, ty internal.Type) (program.Expr, *CompileError) {
	exprTy, expr, err := p.CompileExpr(c)
	if err != nil {
		return nil, err
	}
	if exprTy != ty {
		return nil, LogicError(c, fmt.Errorf("wrong type: expected %v and found %v", ty, exprTy))
	}
	return expr, err
}

func (p *parseVisitor) CompileExpr(c parser.IExpressionContext) (internal.Type, program.Expr, *CompileError) {
	switch c := c.(type) {
	case *parser.ExprAddSubContext:
		lhsType, lhs, err := p.CompileExpr(c.GetLhs())
		if err != nil {
			return 0, nil, err
		}
		switch lhsType {
		case internal.TypeNumber:
			rhs, err := p.CompileExprTy(c.GetRhs(), internal.TypeNumber)
			if err != nil {
				return 0, nil, err
			}
			var expr program.Expr
			switch c.GetOp().GetTokenType() {
			case parser.NumScriptLexerOP_ADD:
				expr = program.ExprNumberAdd{
					Lhs: lhs,
					Rhs: rhs,
				}
			case parser.NumScriptLexerOP_SUB:
				expr = program.ExprNumberSub{
					Lhs: lhs,
					Rhs: rhs,
				}
			}
			return internal.TypeNumber, expr, nil
		case internal.TypeMonetary:
			rhs, err := p.CompileExprTy(c.GetRhs(), internal.TypeMonetary)
			if err != nil {
				return 0, nil, err
			}
			var expr program.Expr
			switch c.GetOp().GetTokenType() {
			case parser.NumScriptLexerOP_ADD:
				expr = program.ExprMonetaryAdd{
					Lhs: lhs,
					Rhs: rhs,
				}
			case parser.NumScriptLexerOP_SUB:
				expr = program.ExprMonetarySub{
					Lhs: lhs,
					Rhs: rhs,
				}
			}
			return internal.TypeMonetary, expr, nil
		default:
			return 0, nil, LogicError(c, errors.New("tried to do arithmetic with wrong type"))
		}
	case *parser.ExprLiteralContext:
		ty, value, err := p.CompileLit(c.GetLit())
		if err != nil {
			return 0, nil, err
		}
		return ty, program.ExprLiteral{Value: value}, nil
	case *parser.ExprVariableContext:
		name := c.GetVar_().GetText()[1:] // strip '$' prefix
		if ty, ok := p.vars[name]; ok {
			return ty, program.ExprVariable(name), nil
		} else {
			return 0, nil, LogicError(c, errors.New("variable not declared"))
		}
	case *parser.ExprMonetaryNewContext:
		asset, compErr := p.CompileExprTy(c.Monetary().GetAsset(), internal.TypeAsset)
		if compErr != nil {
			return 0, nil, compErr
		}
		amt, compErr := p.CompileExprTy(c.Monetary().GetAmt(), internal.TypeNumber)
		if compErr != nil {
			return 0, nil, compErr
		}
		return internal.TypeMonetary, program.ExprMonetaryNew{
			Asset:  asset,
			Amount: amt,
		}, nil
	default:
		return 0, nil, InternalError(c)
	}
}

func (p *parseVisitor) CompileLit(c parser.ILiteralContext) (internal.Type, internal.Value, *CompileError) {
	switch c := c.(type) {
	case *parser.LitAccountContext:
		account := internal.AccountAddress(c.GetText()[1:])
		return internal.TypeAccount, account, nil
	case *parser.LitAssetContext:
		asset := internal.Asset(c.GetText())
		return internal.TypeAsset, asset, nil
	case *parser.LitNumberContext:
		number, err := internal.ParseNumber(c.GetText())
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		return internal.TypeNumber, number, nil
	case *parser.LitStringContext:
		str := internal.String(strings.Trim(c.GetText(), `"`))
		return internal.TypeString, str, nil
	case *parser.LitPortionContext:
		portion, err := internal.ParsePortionSpecific(c.GetText())
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		return internal.TypePortion, *portion, nil
	default:
		return 0, nil, InternalError(c)
	}
}

func (p *parseVisitor) CompileSend(c *parser.SendContext) (program.Instruction, *CompileError) {
	mon, err := p.CompileExprTy(c.GetMon(), internal.TypeMonetary)
	if err != nil {
		return nil, err
	}
	valueAwareSource, err := p.CompileValueAwareSource(c.GetSrc())
	if err != nil {
		return nil, err
	}
	destination, err := p.CompileDestination(c.GetDest())
	if err != nil {
		return nil, err
	}
	return program.InstructionAllocate{
		Funding: program.ExprTake{
			Amount: mon,
			Source: valueAwareSource,
		},
		Destination: destination,
	}, nil
}

func (p *parseVisitor) CompileSendAll(c *parser.SendAllContext) (program.Instruction, *CompileError) {
	source, hasFallback, err := p.CompileSource(c.GetSrc())
	if err != nil {
		return nil, err
	}
	asset, err := p.CompileExprTy(c.GetMonAll().GetAsset(), internal.TypeAsset)
	if err != nil {
		return nil, err
	}
	if hasFallback {
		return nil, LogicError(c, errors.New("cannot take all balance of an unlimited source"))
	}
	destination, err := p.CompileDestination(c.GetDest())
	if err != nil {
		return nil, err
	}
	return program.InstructionAllocate{
		Funding: program.ExprTakeAll{
			Asset:  asset,
			Source: source,
		},
		Destination: destination,
	}, nil
}

func (p *parseVisitor) CompileSetTxMeta(ctx *parser.SetTxMetaContext) (program.Instruction, *CompileError) {
	_, value, err := p.CompileExpr(ctx.GetValue())
	if err != nil {
		return nil, err
	}
	return program.InstructionSetTxMeta{
		Key:   strings.Trim(ctx.GetKey().GetText(), `"`),
		Value: value,
	}, nil

}

func (p *parseVisitor) CompileSetAccountMeta(ctx *parser.SetAccountMetaContext) (program.Instruction, *CompileError) {
	account, err := p.CompileExprTy(ctx.GetAcc(), internal.TypeAccount)
	if err != nil {
		return nil, err
	}

	_, value, err := p.CompileExpr(ctx.GetValue())
	if err != nil {
		return nil, err
	}

	return program.InstructionSetAccountMeta{
		Account: account,
		Key:     strings.Trim(ctx.GetKey().GetText(), `"`),
		Value:   value,
	}, nil
}

func (p *parseVisitor) CompileSave(ctx *parser.SaveFromAccountContext) (program.Instruction, *CompileError) {
	if monAll := ctx.GetMonAll(); monAll != nil {
		asset, err := p.CompileExprTy(ctx.MonetaryAll().GetAsset(), internal.TypeAsset)
		if err != nil {
			return nil, err
		}
		account, err := p.CompileExprTy(ctx.GetAcc(), internal.TypeAccount)
		if err != nil {
			return nil, err
		}
		return program.InstructionSaveAll{
			Asset:   asset,
			Account: account,
		}, nil
	} else if mon := ctx.GetMon(); mon != nil {
		mon, err := p.CompileExprTy(ctx.GetMon(), internal.TypeMonetary)
		if err != nil {
			return nil, err
		}
		account, err := p.CompileExprTy(ctx.GetAcc(), internal.TypeAccount)
		if err != nil {
			return nil, err
		}
		return program.InstructionSave{
			Amount:  mon,
			Account: account,
		}, nil
	} else {
		return nil, InternalError(ctx)
	}
}

func (p *parseVisitor) CompilePrint(ctx *parser.PrintContext) (program.Instruction, *CompileError) {
	_, expr, err := p.CompileExpr(ctx.GetExpr())
	if err != nil {
		return nil, err
	}
	return program.InstructionPrint{Expr: expr}, nil
}

func (p *parseVisitor) CompileVars(c *parser.VarListDeclContext) ([]program.VarDecl, *CompileError) {
	varsDecl := make([]program.VarDecl, 0)

	for _, v := range c.GetV() {
		name := v.GetName().GetText()[1:]
		if _, ok := p.vars[name]; ok {
			return nil, LogicError(c, fmt.Errorf("duplicate variable $%s", name))
		}
		var ty internal.Type
		switch v.GetTy().GetText() {
		case "account":
			ty = internal.TypeAccount
		case "asset":
			ty = internal.TypeAsset
		case "number":
			ty = internal.TypeNumber
		case "string":
			ty = internal.TypeString
		case "monetary":
			ty = internal.TypeMonetary
		case "portion":
			ty = internal.TypePortion
		default:
			return nil, InternalError(c)
		}

		p.vars[name] = ty

		varDecl := program.VarDecl{
			Typ:    ty,
			Name:   name,
			Origin: nil,
		}

		switch c := v.GetOrig().(type) {
		case *parser.OriginAccountMetaContext:
			account, compErr := p.CompileExprTy(c.GetAccount(), internal.TypeAccount)
			if compErr != nil {
				return nil, compErr
			}
			key := strings.Trim(c.GetKey().GetText(), `"`)
			varDecl.Origin = program.VarOriginMeta{
				Account: account,
				Key:     key,
			}
		case *parser.OriginAccountBalanceContext:
			if ty != internal.TypeMonetary {
				return nil, LogicError(c, fmt.Errorf(
					"variable $%s: type should be 'monetary' to pull account balance", name))
			}
			account, compErr := p.CompileExprTy(c.GetAccount(), internal.TypeAccount)
			if compErr != nil {
				return nil, compErr
			}
			asset, compErr := p.CompileExprTy(c.GetAsset(), internal.TypeAsset)
			if compErr != nil {
				return nil, compErr
			}
			varDecl.Origin = program.VarOriginBalance{
				Account: account,
				Asset:   asset,
			}
		}
		varsDecl = append(varsDecl, varDecl)
	}

	return varsDecl, nil
}

func (p *parseVisitor) CompileScript(c parser.IScriptContext) (*program.Program, *CompileError) {
	var varsDecl []program.VarDecl
	var instructions []program.Instruction
	var err *CompileError
	switch c := c.(type) {
	case *parser.ScriptContext:
		vars := c.GetVars()
		if vars != nil {
			switch c := vars.(type) {
			case *parser.VarListDeclContext:
				varsDecl, err = p.CompileVars(c)
				if err != nil {
					return nil, err
				}
			default:
				return nil, InternalError(c)
			}
		}

		for _, statement := range c.GetStmts() {
			var instr program.Instruction
			var err *CompileError
			switch c := statement.(type) {
			case *parser.PrintContext:
				instr, err = p.CompilePrint(c)
			case *parser.FailContext:
				instr = program.InstructionFail{}
			case *parser.SaveFromAccountContext:
				instr, err = p.CompileSave(c)
			case *parser.SendContext:
				instr, err = p.CompileSend(c)
			case *parser.SendAllContext:
				instr, err = p.CompileSendAll(c)
			case *parser.SetTxMetaContext:
				instr, err = p.CompileSetTxMeta(c)
			case *parser.SetAccountMetaContext:
				instr, err = p.CompileSetAccountMeta(c)
			default:
				return nil, InternalError(c)
			}
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, instr)
		}
	default:
		return nil, InternalError(c)
	}

	return &program.Program{
		VarsDecl:    varsDecl,
		Instruction: instructions,
	}, nil
}

type CompileArtifacts struct {
	Source  string
	Tokens  []antlr.Token
	Errors  []CompileError
	Program *program.Program
}

func CompileFull(input string) CompileArtifacts {
	artifacts := CompileArtifacts{
		Source: input,
	}

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

	artifacts.Tokens = stream.GetAllTokens()
	artifacts.Errors = append(artifacts.Errors, errListener.Errors...)

	if len(errListener.Errors) != 0 {
		return artifacts
	}

	visitor := parseVisitor{
		errListener: errListener,
		vars:        make(map[string]internal.Type),
	}

	program, err := visitor.CompileScript(tree)
	if err != nil {
		artifacts.Errors = append(artifacts.Errors, *err)
		return artifacts
	}

	artifacts.Program = program

	return artifacts
}

func Compile(input string) (*program.Program, error) {
	artifacts := CompileFull(input)
	if len(artifacts.Errors) > 0 {
		err := CompileErrorList{
			Errors: artifacts.Errors,
			Source: artifacts.Source,
		}
		return nil, &err
	}

	return artifacts.Program, nil
}
