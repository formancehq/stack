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
	if exprTy != ty {
		return nil, LogicError(c, fmt.Errorf("wrong type: expected %v and found %v", ty, exprTy))
	}
	return expr, err
}

func (p *parseVisitor) CompileExpr(c parser.IExpressionContext) (internal.Type, program.Expr, *CompileError) {
	switch c := c.(type) {
	case *parser.ExprAddSubContext:
		ty, lhs, err := p.CompileExpr(c.GetLhs())
		if err != nil {
			return 0, nil, err
		}
		if ty != internal.TypeNumber {
			return 0, nil, LogicError(c, errors.New("tried to do arithmetic with wrong type"))
		}
		ty, rhs, err := p.CompileExpr(c.GetRhs())
		if err != nil {
			return 0, nil, err
		}
		if ty != internal.TypeNumber {
			return 0, nil, LogicError(c, errors.New("tried to do arithmetic with wrong type"))
		}
		var op byte
		switch c.GetOp().GetTokenType() {
		case parser.NumScriptLexerOP_ADD:
			op = program.OP_ADD
		case parser.NumScriptLexerOP_SUB:
			op = program.OP_SUB
		}
		return internal.TypeNumber, program.ExprInfix{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		}, nil
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

func (p *parseVisitor) CompileSend(c *parser.SendContext) (program.Statement, *CompileError) {
	mon, err := p.CompileExprTy(c.GetMon(), internal.TypeMonetary)
	if err != nil {
		return nil, err
	}
	value_aware_source, err := p.CompileValueAwareSource(c.GetSrc())
	if err != nil {
		return nil, err
	}
	destination, err := p.CompileDestination(c.GetDest())
	if err != nil {
		return nil, err
	}
	return program.StatementAllocate{
		Funding: program.ExprTake{
			Amount: mon,
			Source: value_aware_source,
		},
		Destination: destination,
	}, nil
}

func (p *parseVisitor) CompileSendAll(c *parser.SendAllContext) (program.Statement, *CompileError) {
	source, has_fallback, err := p.CompileSource(c.GetSrc())
	if err != nil {
		return nil, err
	}
	asset, err := p.CompileExprTy(c.GetMonAll().GetAsset(), internal.TypeAsset)
	if err != nil {
		return nil, err
	}
	if has_fallback {
		return nil, LogicError(c, errors.New("cannot take all balance of an unlimited source"))
	}
	destination, err := p.CompileDestination(c.GetDest())
	if err != nil {
		return nil, err
	}
	return program.StatementAllocate{
		Funding: program.ExprTakeAll{
			Asset:  asset,
			Source: source,
		},
		Destination: destination,
	}, nil
}

func (p *parseVisitor) CompileSetTxMeta(ctx *parser.SetTxMetaContext) (program.Statement, *CompileError) {
	_, value, err := p.CompileExpr(ctx.GetValue())
	if err != nil {
		return nil, err
	}
	return program.StatementSetTxMeta{
		Key:   strings.Trim(ctx.GetKey().GetText(), `"`),
		Value: value,
	}, nil

}

func (p *parseVisitor) CompileSetAccountMeta(ctx *parser.SetAccountMetaContext) (program.Statement, *CompileError) {
	account, err := p.CompileExprTy(ctx.GetAcc(), internal.TypeAccount)
	if err != nil {
		return nil, err
	}

	_, value, err := p.CompileExpr(ctx.GetValue())
	if err != nil {
		return nil, err
	}

	return program.StatementSetAccountMeta{
		Account: account,
		Key:     strings.Trim(ctx.GetKey().GetText(), `"`),
		Value:   value,
	}, nil
}

func (p *parseVisitor) CompileSave(ctx *parser.SaveFromAccountContext) (program.Statement, *CompileError) {
	typ, amt_expr, err := p.CompileExpr(ctx.GetMon())
	if err != nil {
		return nil, err
	}
	acc_expr, err := p.CompileExprTy(ctx.GetAcc(), internal.TypeAccount)
	if err != nil {
		return nil, err
	}
	switch typ {
	case internal.TypeAsset:
		return program.StatementSaveAll{Asset: amt_expr, Account: acc_expr}, nil
	case internal.TypeMonetary:
		return program.StatementSave{Amount: amt_expr, Account: acc_expr}, nil
	default:
		return nil, InternalError(ctx)
	}
}

func (p *parseVisitor) CompilePrint(ctx *parser.PrintContext) (program.Statement, *CompileError) {
	_, expr, err := p.CompileExpr(ctx.GetExpr())
	if err != nil {
		return nil, err
	}
	return program.StatementPrint{Expr: expr}, nil
}

func (p *parseVisitor) CompileVars(c *parser.VarListDeclContext) ([]program.VarDecl, *CompileError) {
	vars_decl := make([]program.VarDecl, 0)

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

		var_decl := program.VarDecl{
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
			var_decl.Origin = program.VarOriginMeta{
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
			var_decl.Origin = program.VarOriginBalance{
				Account: account,
				Asset:   asset,
			}
		}
		vars_decl = append(vars_decl, var_decl)
	}

	return vars_decl, nil
}

func (p *parseVisitor) CompileScript(c parser.IScriptContext) (*program.Program, *CompileError) {
	var vars_decl []program.VarDecl
	var statements []program.Statement
	var err *CompileError
	switch c := c.(type) {
	case *parser.ScriptContext:
		vars := c.GetVars()
		if vars != nil {
			switch c := vars.(type) {
			case *parser.VarListDeclContext:
				vars_decl, err = p.CompileVars(c)
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
				stmt, err = p.CompilePrint(c)
			case *parser.FailContext:
				stmt = program.StatementFail{}
			case *parser.SaveFromAccountContext:
				stmt, err = p.CompileSave(c)
			case *parser.SendContext:
				stmt, err = p.CompileSend(c)
			case *parser.SendAllContext:
				stmt, err = p.CompileSendAll(c)
			case *parser.SetTxMetaContext:
				stmt, err = p.CompileSetTxMeta(c)
			case *parser.SetAccountMetaContext:
				stmt, err = p.CompileSetAccountMeta(c)
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
