package compiler

import (
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/numary/ledger/pkg/core"
	"github.com/numary/ledger/pkg/machine/script/parser"
	"github.com/numary/ledger/pkg/machine/vm/program"
	"github.com/pkg/errors"
)

type parseVisitor struct {
	errListener *ErrorListener
	vars        map[string]core.Type
}

func (p *parseVisitor) isWorld(expr parser.IExpressionContext) bool {
	if lit, ok := expr.(*parser.ExprLiteralContext); ok {
		_, value, _ := p.VisitLit(lit.GetLit())
		return core.ValueEquals(value, core.AccountAddress("world"))
	} else {
		return false
	}
}

func (p *parseVisitor) VisitExprTy(c parser.IExpressionContext, ty core.Type) (program.Expr, *CompileError) {
	exprTy, expr, err := p.VisitExpr(c)
	if exprTy != ty {
		return nil, LogicError(c, fmt.Errorf("wrong type: expected %v and found %v", ty, exprTy))
	}
	return expr, err
}

func (p *parseVisitor) VisitExpr(c parser.IExpressionContext) (core.Type, program.Expr, *CompileError) {
	switch c := c.(type) {
	case *parser.ExprAddSubContext:
		ty, lhs, err := p.VisitExpr(c.GetLhs())
		if err != nil {
			return 0, nil, err
		}
		if ty != core.TypeNumber {
			return 0, nil, LogicError(c, errors.New("tried to do arithmetic with wrong type"))
		}
		ty, rhs, err := p.VisitExpr(c.GetRhs())
		if err != nil {
			return 0, nil, err
		}
		if ty != core.TypeNumber {
			return 0, nil, LogicError(c, errors.New("tried to do arithmetic with wrong type"))
		}
		return core.TypeNumber, program.ExprInfix{
			Op:  0,
			Lhs: lhs,
			Rhs: rhs,
		}, nil
	case *parser.ExprLiteralContext:
		ty, value, err := p.VisitLit(c.GetLit())
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
		asset, compErr := p.VisitExprTy(c.Monetary().GetAsset(), core.TypeAsset)
		if compErr != nil {
			return 0, nil, compErr
		}
		amt, compErr := p.VisitExprTy(c.Monetary().GetAmt(), core.TypeNumber)
		if compErr != nil {
			return 0, nil, compErr
		}
		return core.TypeMonetary, program.ExprMonetaryNew{
			Asset:  asset,
			Amount: amt,
		}, nil
	default:
		return 0, nil, InternalError(c)
	}
}

func (p *parseVisitor) VisitLit(c parser.ILiteralContext) (core.Type, core.Value, *CompileError) {
	switch c := c.(type) {
	case *parser.LitAccountContext:
		account := core.AccountAddress(c.GetText()[1:])
		return core.TypeAccount, account, nil
	case *parser.LitAssetContext:
		asset := core.Asset(c.GetText())
		return core.TypeAsset, asset, nil
	case *parser.LitNumberContext:
		number, err := core.ParseNumber(c.GetText())
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		return core.TypeNumber, number, nil
	case *parser.LitStringContext:
		str := core.String(strings.Trim(c.GetText(), `"`))
		return core.TypeString, str, nil
	case *parser.LitPortionContext:
		portion, err := core.ParsePortionSpecific(c.GetText())
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		return core.TypePortion, *portion, nil
	default:
		return 0, nil, InternalError(c)
	}
}

func (p *parseVisitor) VisitSend(c *parser.SendContext) (*program.StatementAllocate, *CompileError) {
	mon, err := p.VisitExprTy(c.GetMon(), core.TypeMonetary)
	if err != nil {
		return nil, err
	}
	value_aware_source, err := p.VisitValueAwareSource(c.GetSrc())
	if err != nil {
		return nil, err
	}
	destination, err := p.VisitDestination(c.GetDest())
	if err != nil {
		return nil, err
	}
	return &program.StatementAllocate{
		Funding: program.ExprTake{
			Amount: mon,
			Source: value_aware_source,
		},
		Destination: destination,
	}, nil
}

func (p *parseVisitor) VisitSendAll(c *parser.SendAllContext) (*program.StatementAllocate, *CompileError) {
	source, has_fallback, err := p.VisitSource(c.GetSrc())
	if err != nil {
		return nil, err
	}
	asset, err := p.VisitExprTy(c.GetMonAll().GetAsset(), core.TypeAsset)
	if err != nil {
		return nil, err
	}
	if has_fallback {
		return nil, LogicError(c, errors.New("cannot take all balance of an unlimited source"))
	}
	destination, err := p.VisitDestination(c.GetDest())
	if err != nil {
		return nil, err
	}
	return &program.StatementAllocate{
		Funding: program.ExprTakeAll{
			Asset:  asset,
			Source: source,
		},
		Destination: destination,
	}, nil
}

func (p *parseVisitor) VisitSetTxMeta(ctx *parser.SetTxMetaContext) (*program.StatementSetTxMeta, *CompileError) {
	_, value, err := p.VisitExpr(ctx.GetValue())
	if err != nil {
		return nil, err
	}
	return &program.StatementSetTxMeta{
		Key:   strings.Trim(ctx.GetKey().GetText(), `"`),
		Value: value,
	}, nil

}

func (p *parseVisitor) VisitSetAccountMeta(ctx *parser.SetAccountMetaContext) (*program.StatementSetAccountMeta, *CompileError) {
	account, err := p.VisitExprTy(ctx.GetAcc(), core.TypeAccount)
	if err != nil {
		return nil, err
	}

	_, value, err := p.VisitExpr(ctx.GetValue())
	if err != nil {
		return nil, err
	}

	return &program.StatementSetAccountMeta{
		Account: account,
		Key:     strings.Trim(ctx.GetKey().GetText(), `"`),
		Value:   value,
	}, nil
}

func (p *parseVisitor) VisitPrint(ctx *parser.PrintContext) (program.Statement, *CompileError) {
	_, expr, err := p.VisitExpr(ctx.GetExpr())
	if err != nil {
		return nil, err
	}
	return program.StatementPrint{Expr: expr}, nil
}

func (p *parseVisitor) VisitVars(c *parser.VarListDeclContext) ([]program.VarDecl, *CompileError) {
	vars_decl := make([]program.VarDecl, 0)

	for _, v := range c.GetV() {
		name := v.GetName().GetText()[1:]
		if _, ok := p.vars[name]; ok {
			return nil, LogicError(c, fmt.Errorf("duplicate variable $%s", name))
		}
		var ty core.Type
		switch v.GetTy().GetText() {
		case "account":
			ty = core.TypeAccount
		case "asset":
			ty = core.TypeAsset
		case "number":
			ty = core.TypeNumber
		case "string":
			ty = core.TypeString
		case "monetary":
			ty = core.TypeMonetary
		case "portion":
			ty = core.TypePortion
		default:
			return nil, InternalError(c)
		}

		p.vars[name] = ty

		var_decl := program.VarDecl{
			Ty:     ty,
			Name:   name,
			Origin: nil,
		}

		switch c := v.GetOrig().(type) {
		case *parser.OriginAccountMetaContext:
			account, compErr := p.VisitExprTy(c.GetAccount(), core.TypeAccount)
			if compErr != nil {
				return nil, compErr
			}
			key := strings.Trim(c.GetKey().GetText(), `"`)
			var_decl.Origin = program.VarOriginMeta{
				Account: account,
				Key:     key,
			}
		case *parser.OriginAccountBalanceContext:
			if ty != core.TypeMonetary {
				return nil, LogicError(c, fmt.Errorf(
					"variable $%s: type should be 'monetary' to pull account balance", name))
			}
			account, compErr := p.VisitExprTy(c.GetAccount(), core.TypeAccount)
			if compErr != nil {
				return nil, compErr
			}
			asset, compErr := p.VisitExprTy(c.GetAsset(), core.TypeAsset)
			if compErr != nil {
				return nil, compErr
			}
			var_decl.Origin = program.VarOriginBalance{
				Account: account,
				Asset:   asset,
			}
		}
		vars_decl = append(vars_decl, var_decl)
	}

	return vars_decl, nil
}

func (p *parseVisitor) VisitScript(c parser.IScriptContext) (*program.Program, *CompileError) {
	var vars_decl []program.VarDecl
	var statements []program.Statement
	var err *CompileError
	switch c := c.(type) {
	case *parser.ScriptContext:
		vars := c.GetVars()
		if vars != nil {
			switch c := vars.(type) {
			case *parser.VarListDeclContext:
				vars_decl, err = p.VisitVars(c)
				if err != nil {
					return nil, err
				}
			default:
				return nil, InternalError(c)
			}
		}

		for _, statement := range c.GetStmts() {
			var stmt program.Statement
			var err *CompileError
			switch c := statement.(type) {
			case *parser.PrintContext:
				stmt, err = p.VisitPrint(c)
			case *parser.FailContext:
				stmt = program.StatementFail{}
			case *parser.SendContext:
				stmt, err = p.VisitSend(c)
			case *parser.SetTxMetaContext:
				stmt, err = p.VisitSetTxMeta(c)
			case *parser.SetAccountMetaContext:
				stmt, err = p.VisitSetAccountMeta(c)
			default:
				return nil, InternalError(c)
			}
			if err != nil {
				return nil, err
			}
			statements = append(statements, stmt)
		}
	default:
		return nil, InternalError(c)
	}

	return &program.Program{
		VarsDecl:   vars_decl,
		Statements: statements,
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
		vars:        make(map[string]core.Type),
	}

	program, err := visitor.VisitScript(tree)
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
