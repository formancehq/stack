package compiler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/formancehq/ledger/pkg/machine/internal"
	"github.com/formancehq/ledger/pkg/machine/script/parser"
	"github.com/formancehq/ledger/pkg/machine/vm/program"
	"github.com/pkg/errors"
)

type parseVisitor struct {
	errListener  *ErrorListener
	instructions []byte
	// resources must not exceed 65536 elements
	resources []program.Resource
	// sources store all source accounts
	// a source can be also a destination of another posting
	sources map[internal.Address]struct{}
	// varIdx maps name to resource index
	varIdx map[string]internal.Address
	// needBalances store for each account, the set of assets needed
	neededBalances map[internal.Address]map[internal.Address]struct{}
}

// Allocates constants if it hasn't already been,
// and returns its resource address.
func (p *parseVisitor) findConstant(constant program.Constant) (*internal.Address, bool) {
	for i := 0; i < len(p.resources); i++ {
		if c, ok := p.resources[i].(program.Constant); ok {
			if internal.ValueEquals(c.Inner, constant.Inner) {
				addr := internal.Address(i)
				return &addr, true
			}
		}
	}
	return nil, false
}

func (p *parseVisitor) AllocateResource(res program.Resource) (*internal.Address, error) {
	if c, ok := res.(program.Constant); ok {
		idx, ok := p.findConstant(c)
		if ok {
			return idx, nil
		}
	}
	if len(p.resources) >= 65536 {
		return nil, errors.New("number of unique constants exceeded 65536")
	}
	p.resources = append(p.resources, res)
	addr := internal.NewAddress(uint16(len(p.resources) - 1))
	return &addr, nil
}

func (p *parseVisitor) isWorld(addr internal.Address) bool {
	idx := int(addr)
	if idx < len(p.resources) {
		if c, ok := p.resources[idx].(program.Constant); ok {
			if acc, ok := c.Inner.(internal.AccountAddress); ok {
				if string(acc) == "world" {
					return true
				}
			}
		}
	}
	return false
}

func (p *parseVisitor) VisitVariable(c parser.IVariableContext, push bool) (internal.Type, *internal.Address, *CompileError) {
	name := c.GetText()[1:] // strip '$' prefix
	if idx, ok := p.varIdx[name]; ok {
		res := p.resources[idx]
		if push {
			p.PushAddress(idx)
		}
		return res.GetType(), &idx, nil
	} else {
		return 0, nil, LogicError(c, errors.New("variable not declared"))
	}
}

func (p *parseVisitor) VisitExpr(c parser.IExpressionContext, push bool) (internal.Type, *internal.Address, *CompileError) {
	switch c := c.(type) {
	case *parser.ExprAddSubContext:
		ty, _, err := p.VisitExpr(c.GetLhs(), push)
		if err != nil {
			return 0, nil, err
		}
		if ty != internal.TypeNumber {
			return 0, nil, LogicError(c, errors.New("tried to do arithmetic with wrong type"))
		}
		ty, _, err = p.VisitExpr(c.GetRhs(), push)
		if err != nil {
			return 0, nil, err
		}
		if ty != internal.TypeNumber {
			return 0, nil, LogicError(c, errors.New("tried to do arithmetic with wrong type"))
		}
		if push {
			switch c.GetOp().GetTokenType() {
			case parser.NumScriptLexerOP_ADD:
				p.AppendInstruction(program.OP_IADD)
			case parser.NumScriptLexerOP_SUB:
				p.AppendInstruction(program.OP_ISUB)
			}
		}
		return internal.TypeNumber, nil, nil
	case *parser.ExprLiteralContext:
		ty, addr, err := p.VisitLit(c.GetLit(), push)
		if err != nil {
			return 0, nil, err
		}
		return ty, addr, nil
	case *parser.ExprVariableContext:
		ty, addr, err := p.VisitVariable(c.GetVar_(), push)
		return ty, addr, err
	default:
		return 0, nil, InternalError(c)
	}
}

func (p *parseVisitor) VisitLit(c parser.ILiteralContext, push bool) (internal.Type, *internal.Address, *CompileError) {
	switch c := c.(type) {
	case *parser.LitAccountContext:
		account := internal.AccountAddress(c.GetText()[1:])
		addr, err := p.AllocateResource(program.Constant{Inner: account})
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		if push {
			p.PushAddress(*addr)
		}
		return internal.TypeAccount, addr, nil
	case *parser.LitAssetContext:
		asset := internal.Asset(c.GetText())
		addr, err := p.AllocateResource(program.Constant{Inner: asset})
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		if push {
			p.PushAddress(*addr)
		}
		return internal.TypeAsset, addr, nil
	case *parser.LitNumberContext:
		number, err := internal.ParseNumber(c.GetText())
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		if push {
			err := p.PushInteger(number)
			if err != nil {
				return 0, nil, LogicError(c, err)
			}
		}
		return internal.TypeNumber, nil, nil
	case *parser.LitStringContext:
		addr, err := p.AllocateResource(program.Constant{
			Inner: internal.String(strings.Trim(c.GetText(), `"`)),
		})
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		if push {
			p.PushAddress(*addr)
		}
		return internal.TypeString, addr, nil
	case *parser.LitPortionContext:
		portion, err := internal.ParsePortionSpecific(c.GetText())
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		addr, err := p.AllocateResource(program.Constant{Inner: *portion})
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		if push {
			p.PushAddress(*addr)
		}
		return internal.TypePortion, addr, nil
	case *parser.LitMonetaryContext:
		typ, assetAddr, compErr := p.VisitExpr(c.Monetary().GetAsset(), false)
		if compErr != nil {
			return 0, nil, compErr
		}
		if typ != internal.TypeAsset {
			return 0, nil, LogicError(c, fmt.Errorf(
				"the expression in monetary literal should be of type '%s' instead of '%s'",
				internal.TypeAsset, typ))
		}
		amt, err := internal.ParseMonetaryInt(c.Monetary().GetAmt().GetText())
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		monAddr, err := p.AllocateResource(program.Monetary{
			Asset:  *assetAddr,
			Amount: amt,
		})
		if err != nil {
			return 0, nil, LogicError(c, err)
		}
		if push {
			p.PushAddress(*monAddr)
		}
		return internal.TypeMonetary, monAddr, nil
	default:
		return 0, nil, InternalError(c)
	}
}

func (p *parseVisitor) VisitSend(c *parser.SendContext) *CompileError {
	var (
		accounts map[internal.Address]struct{}
		addr     *internal.Address
		compErr  *CompileError
		typ      internal.Type
	)

	if monAll := c.GetMonAll(); monAll != nil {
		typ, addr, compErr = p.VisitExpr(monAll.GetAsset(), false)
		if compErr != nil {
			return compErr
		}
		if typ != internal.TypeAsset {
			return LogicError(c, fmt.Errorf(
				"send monetary all: the expression should be of type 'asset' instead of '%s'", typ))
		}

		accounts, compErr = p.VisitValueAwareSource(c.GetSrc(), func() {
			p.PushAddress(*addr)
		}, nil)
		if compErr != nil {
			return compErr
		}
	} else if mon := c.GetMon(); mon != nil {
		typ, addr, compErr = p.VisitExpr(mon, false)
		if compErr != nil {
			return compErr
		}
		if typ != internal.TypeMonetary {
			return LogicError(c, fmt.Errorf(
				"send monetary: the expression should be of type 'monetary' instead of '%s'", typ))
		}

		accounts, compErr = p.VisitValueAwareSource(c.GetSrc(), func() {
			p.PushAddress(*addr)
			p.AppendInstruction(program.OP_ASSET)
		}, addr)
		if compErr != nil {
			return compErr
		}
	}

	for acc := range accounts {
		if b, ok := p.neededBalances[acc]; ok {
			b[*addr] = struct{}{}
		} else {
			p.neededBalances[acc] = map[internal.Address]struct{}{
				*addr: {},
			}
		}
	}

	if err := p.VisitDestination(c.GetDest()); err != nil {
		return err
	}

	return nil
}

func (p *parseVisitor) VisitSetTxMeta(ctx *parser.SetTxMetaContext) *CompileError {
	_, _, compErr := p.VisitExpr(ctx.GetValue(), true)
	if compErr != nil {
		return compErr
	}

	keyAddr, err := p.AllocateResource(program.Constant{
		Inner: internal.String(strings.Trim(ctx.GetKey().GetText(), `"`)),
	})
	if err != nil {
		return LogicError(ctx, err)
	}
	p.PushAddress(*keyAddr)

	p.AppendInstruction(program.OP_TX_META)

	return nil
}

func (p *parseVisitor) VisitSetAccountMeta(ctx *parser.SetAccountMetaContext) *CompileError {
	_, _, compErr := p.VisitExpr(ctx.GetValue(), true)
	if compErr != nil {
		return compErr
	}

	keyAddr, err := p.AllocateResource(program.Constant{
		Inner: internal.String(strings.Trim(ctx.GetKey().GetText(), `"`)),
	})
	if err != nil {
		return LogicError(ctx, err)
	}
	p.PushAddress(*keyAddr)

	ty, accAddr, compErr := p.VisitExpr(ctx.GetAcc(), false)
	if compErr != nil {
		return compErr
	}
	if ty != internal.TypeAccount {
		return LogicError(ctx, fmt.Errorf(
			"set_account_meta: expression is of type %s, and should be of type account", ty))
	}
	p.PushAddress(*accAddr)

	p.AppendInstruction(program.OP_ACCOUNT_META)

	return nil
}

func (p *parseVisitor) VisitSaveFromAccount(c *parser.SaveFromAccountContext) *CompileError {
	var (
		typ     internal.Type
		addr    *internal.Address
		compErr *CompileError
	)
	if monAll := c.GetMonAll(); monAll != nil {
		typ, addr, compErr = p.VisitExpr(monAll.GetAsset(), false)
		if compErr != nil {
			return compErr
		}
		if typ != internal.TypeAsset {
			return LogicError(c, fmt.Errorf(
				"save monetary all from account: the first expression should be of type 'asset' instead of '%s'", typ))
		}
	} else if mon := c.GetMon(); mon != nil {
		typ, addr, compErr = p.VisitExpr(mon, false)
		if compErr != nil {
			return compErr
		}
		if typ != internal.TypeMonetary {
			return LogicError(c, fmt.Errorf(
				"save monetary from account: the first expression should be of type 'monetary' instead of '%s'", typ))
		}
	}
	p.PushAddress(*addr)

	typ, addr, compErr = p.VisitExpr(c.GetAcc(), false)
	if compErr != nil {
		return compErr
	}
	if typ != internal.TypeAccount {
		return LogicError(c, fmt.Errorf(
			"save monetary from account: the second expression should be of type 'account' instead of '%s'", typ))
	}
	p.PushAddress(*addr)

	p.AppendInstruction(program.OP_SAVE)

	return nil
}

func (p *parseVisitor) VisitPrint(ctx *parser.PrintContext) *CompileError {
	_, _, err := p.VisitExpr(ctx.GetExpr(), true)
	if err != nil {
		return err
	}

	p.AppendInstruction(program.OP_PRINT)

	return nil
}

func (p *parseVisitor) VisitVars(c *parser.VarListDeclContext) *CompileError {
	if len(c.GetV()) > 32768 {
		return LogicError(c, fmt.Errorf("number of variables exceeded %v", 32768))
	}

	for _, v := range c.GetV() {
		name := v.GetName().GetText()[1:]
		if _, ok := p.varIdx[name]; ok {
			return LogicError(c, fmt.Errorf("duplicate variable $%s", name))
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
			return InternalError(c)
		}

		var addr *internal.Address
		var err error
		if v.GetOrig() == nil {
			addr, err = p.AllocateResource(program.Variable{Typ: ty, Name: name})
			if err != nil {
				return &CompileError{
					Msg: errors.Wrap(err,
						"allocating variable resource").Error(),
				}
			}
			p.varIdx[name] = *addr
			continue
		}

		switch c := v.GetOrig().(type) {
		case *parser.OriginAccountMetaContext:
			srcTy, src, compErr := p.VisitExpr(c.GetAccount(), false)
			if compErr != nil {
				return compErr
			}
			if srcTy != internal.TypeAccount {
				return LogicError(c, fmt.Errorf(
					"variable $%s: type should be 'account' to pull account metadata", name))
			}
			key := strings.Trim(c.GetKey().GetText(), `"`)
			addr, err = p.AllocateResource(program.VariableAccountMetadata{
				Typ:     ty,
				Name:    name,
				Account: *src,
				Key:     key,
			})
		case *parser.OriginAccountBalanceContext:
			if ty != internal.TypeMonetary {
				return LogicError(c, fmt.Errorf(
					"variable $%s: type should be 'monetary' to pull account balance", name))
			}
			accTy, accAddr, compErr := p.VisitExpr(c.GetAccount(), false)
			if compErr != nil {
				return compErr
			}
			if accTy != internal.TypeAccount {
				return LogicError(c, fmt.Errorf(
					"variable $%s: the first argument to pull account balance should be of type 'account'", name))
			}

			assTy, assAddr, compErr := p.VisitExpr(c.GetAsset(), false)
			if compErr != nil {
				return compErr
			}
			if assTy != internal.TypeAsset {
				return LogicError(c, fmt.Errorf(
					"variable $%s: the second argument to pull account balance should be of type 'asset'", name))
			}
			addr, err = p.AllocateResource(program.VariableAccountBalance{
				Name:    name,
				Account: *accAddr,
				Asset:   *assAddr,
			})
			if err != nil {
				return LogicError(c, err)
			}
		}
		if err != nil {
			return LogicError(c, err)
		}

		p.varIdx[name] = *addr
	}

	return nil
}

func (p *parseVisitor) VisitScript(c parser.IScriptContext) *CompileError {
	switch c := c.(type) {
	case *parser.ScriptContext:
		vars := c.GetVars()
		if vars != nil {
			switch c := vars.(type) {
			case *parser.VarListDeclContext:
				if err := p.VisitVars(c); err != nil {
					return err
				}
			default:
				return InternalError(c)
			}
		}

		for _, stmt := range c.GetStmts() {
			var err *CompileError
			switch c := stmt.(type) {
			case *parser.PrintContext:
				err = p.VisitPrint(c)
			case *parser.FailContext:
				p.AppendInstruction(program.OP_FAIL)
			case *parser.SendContext:
				err = p.VisitSend(c)
			case *parser.SetTxMetaContext:
				err = p.VisitSetTxMeta(c)
			case *parser.SetAccountMetaContext:
				err = p.VisitSetAccountMeta(c)
			case *parser.SaveFromAccountContext:
				err = p.VisitSaveFromAccount(c)
			default:
				return InternalError(c)
			}
			if err != nil {
				return err
			}
		}
	default:
		return InternalError(c)
	}

	return nil
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
		errListener:    errListener,
		instructions:   make([]byte, 0),
		resources:      make([]program.Resource, 0),
		varIdx:         make(map[string]internal.Address),
		neededBalances: make(map[internal.Address]map[internal.Address]struct{}),
		sources:        map[internal.Address]struct{}{},
	}

	err := visitor.VisitScript(tree)
	if err != nil {
		artifacts.Errors = append(artifacts.Errors, *err)
		return artifacts
	}

	sources := make(internal.Addresses, 0)
	for address := range visitor.sources {
		sources = append(sources, address)
	}
	sort.Stable(sources)

	artifacts.Program = &program.Program{
		Instructions:   visitor.instructions,
		Resources:      visitor.resources,
		NeededBalances: visitor.neededBalances,
		Sources:        sources,
	}

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
